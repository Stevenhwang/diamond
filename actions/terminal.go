package actions

import (
	"context"
	"diamond/misc"
	"diamond/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Stevenhwang/gommon/slice"
	"github.com/Stevenhwang/gommon/ssh"
	"github.com/Stevenhwang/gommon/times"
	"github.com/Stevenhwang/gommon/tools"
	"gorm.io/gorm"

	gossh "golang.org/x/crypto/ssh"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 超时控制，分钟
var maxTimeout = 30 * time.Minute
var idleTimeout = 5 * time.Minute

func terminal(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		misc.Logger.Error().Err(err).Str("from", "terminal").Msg("")
	}
	defer ws.Close()
	// 检查server id
	id := c.QueryParam("id")
	if len(id) == 0 {
		ws.WriteMessage(websocket.TextMessage, []byte("请指定服务器id"))
		return nil
	}
	// 检查token
	token := c.QueryParam("token")
	if len(token) == 0 {
		ws.WriteMessage(websocket.TextMessage, []byte("请携带token请求"))
		return nil
	}
	// parse token
	t, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(misc.Config.GetString("jwt.secret")), nil
	})
	// 解析token
	if t != nil {
		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			ws.WriteMessage(websocket.TextMessage, []byte("failed to cast claims as jwt.MapClaims"))
			return nil
		}
		uid := uint(claims["uid"].(float64))
		username := claims["username"].(string)
		// 还要判断用户是否有权限连接此服务器
		if username != "admin" {
			user := models.User{}
			models.DB.Preload("Servers", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at desc")
			}).First(&user, uid)
			// 如果用户被禁用
			if !user.IsActive {
				ws.WriteMessage(websocket.TextMessage, []byte("账号禁用"))
				return nil
			}
			var sids []int
			for _, s := range user.Servers {
				sids = append(sids, int(s.ID))
			}
			intid, err := strconv.Atoi(id)
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return nil
			}
			if !slice.FindValInIntSlice(sids, intid) {
				ws.WriteMessage(websocket.TextMessage, []byte("你没有权限连接此服务器"))
				return nil
			}
		}
		// 显示
		ws.WriteMessage(websocket.TextMessage, []byte("================================\r\n"))
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("连接最长%s, 空闲超时%s\r\n", maxTimeout, idleTimeout)))
		ws.WriteMessage(websocket.TextMessage, []byte("================================\r\n"))
		// 开始连接服务器
		server := models.Server{}
		if res := models.DB.First(&server, id); res.Error != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("服务器不存在"))
			return nil
		}
		// 开始连接流程
		credential := models.Credential{}
		if res := models.DB.First(&credential, server.CredentialID); res.Error != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("认证信息不存在"))
			return nil
		}
		var client *gossh.Client
		spass := tools.AesDecrypt(credential.AuthContent, "0123456789012345")
		if credential.AuthType == 1 {
			cli, err := ssh.GetSSHClientByPassword(server.IP, spass, ssh.SSHOptions{Port: server.Port, User: credential.AuthUser})
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return nil
			}
			client = cli
		}
		if credential.AuthType == 2 {
			cli, err := ssh.GetSSHClientByKey(server.IP, []byte(spass), ssh.SSHOptions{Port: server.Port, User: credential.AuthUser})
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return nil
			}
			client = cli
		}
		defer client.Close()
		session, err := client.NewSession()
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return nil
		}
		defer session.Close()
		// Set up terminal modes
		modes := gossh.TerminalModes{
			gossh.ECHO:          1,     // disable echoing
			gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		// Request pseudo terminal
		cols, _ := strconv.Atoi(c.QueryParam("cols"))
		rows, _ := strconv.Atoi(c.QueryParam("rows"))
		if err := session.RequestPty("xterm", rows, cols, modes); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return nil
		}
		// stdin, stdout
		stdin, err := session.StdinPipe()
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return nil
		}
		stdout, err := session.StdoutPipe()
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return nil
		}
		// Start remote shell
		if err := session.Shell(); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return nil
		}
		// 连接超时控制，默认30分钟
		ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
		defer cancel()
		ch := make(chan []byte)
		exitCh := make(chan bool)
		idletimer := time.NewTimer(idleTimeout) // 如果用户空闲5分钟没输入，也退出
		// 读取用户输入，并传递给远程主机
		go func() {
			for {
				// misc.Logger.Info().Str("from", "terminal").Msg("开始读取用户输入")
				_, msg, err := ws.ReadMessage()
				ch <- msg
				// 收到用户输入，重置timer
				if !idletimer.Stop() {
					<-idletimer.C
				}
				idletimer.Reset(idleTimeout)
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						misc.Logger.Error().Err(err).Str("from", "terminal").Msg("")
					}
					exitCh <- true
					return
				}
				// log.Printf("%s\n", msg)
				// log.Println("开始发送消息到远程主机: ", string(msg))
				// 过滤出 window change 事件
				var resize [2]int
				if err := json.Unmarshal(msg, &resize); err != nil {
					if _, err := stdin.Write(msg); err != nil {
						misc.Logger.Error().Err(err).Str("from", "terminal").Msg("session stdin write error")
						exitCh <- true
						return
					}
				} else {
					misc.Logger.Info().Str("from", "terminal").Msg("更改窗口大小")
					session.WindowChange(resize[1], resize[0])
				}
				// misc.Logger.Info().Str("from", "terminal").Msg("结束读取用户输入")
			}
		}()
		// 读取远程返回
		go func() {
			// 开始记录终端录屏
			now := time.Now()
			time_str, _ := times.GetTimeString(0, times.TimeOptions{Layout: "20060102150405"})
			record_dir := "./records/"
			record_file := fmt.Sprintf("%s%s-%s-%s-web.cast", record_dir, username, server.IP, time_str)
			record := models.Record{User: username, IP: server.IP, File: record_file}
			models.DB.Create(&record)
			var f *os.File
			f, _ = os.Create(record_file)
			defer f.Close()
			// 记录文件头
			timestamp, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float64(now.UnixNano())/float64(1e9)), 64)
			env := map[string]string{
				"SHELL": "/bin/bash",
				"TERM":  "xterm",
			}
			header := map[string]interface{}{
				"version":   2,
				"width":     80,
				"height":    24,
				"timestamp": timestamp,
				"env":       env,
			}
			headerbyte, _ := json.Marshal(header)
			f.WriteString(string(headerbyte) + "\n")
			// 开始读取远程
			for {
				buf := make([]byte, 1024)
				n, err := stdout.Read(buf[:])
				if err != nil {
					misc.Logger.Error().Err(err).Str("from", "terminal").Msg("session stdout read error")
					exitCh <- true
					return
				}
				// 记录终端内容
				nt, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float64(time.Now().UnixNano())/float64(1e9)), 64)
				iodata := []string{fmt.Sprintf("%.9f", nt-timestamp), "o", string(buf[:n])} // 指定读出多少，不然都是补0的多余数据
				iodatabyte, _ := json.Marshal(iodata)
				f.WriteString(string(iodatabyte) + "\n")
				// 输出到 websocket
				err = ws.WriteMessage(websocket.TextMessage, buf[:n])
				if err != nil {
					misc.Logger.Error().Err(err).Str("from", "terminal").Msg("write ws msg error")
					exitCh <- true
					return
				}
				// misc.Logger.Info().Str("from", "terminal").Msg("结束读取远程")
			}
		}()
		// 监控连接超时，空闲超时和退出信号
		for {
			select {
			case <-ctx.Done():
				misc.Logger.Error().Err(ctx.Err()).Str("from", "terminal").Msg("超时")
				if err := ws.WriteMessage(websocket.TextMessage, []byte("连接超时断开连接")); err != nil {
					misc.Logger.Error().Err(err).Str("from", "terminal").Msg("write ws msg error")
				}
				// 超时的话退出连接
				return nil
			case <-ch:
				// misc.Logger.Info().Str("from", "terminal").Msg("没有超时")
				continue
			case <-exitCh:
				misc.Logger.Info().Str("from", "terminal").Msg("收到退出信号")
				return nil
			case <-idletimer.C:
				misc.Logger.Info().Str("from", "terminal").Msg("收到空闲超时信号")
				if err := ws.WriteMessage(websocket.TextMessage, []byte("空闲超时断开连接")); err != nil {
					misc.Logger.Error().Err(err).Str("from", "terminal").Msg("write ws msg error")
				}
				return nil
			}
		}
		// session.Wait()
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("token nil"))
	}
	return nil
}
