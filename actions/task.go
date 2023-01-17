package actions

import (
	"diamond/misc"
	"diamond/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Stevenhwang/gommon/times"
	"github.com/labstack/echo/v4"
)

func getTasks(c echo.Context) error {
	var total int64
	tasks := models.Tasks{}
	if res := models.DB.Model(&models.Task{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c)).Find(&tasks); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": tasks, "total": total})
}

func createTask(c echo.Context) error {
	task := models.Task{}
	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&task); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Create(&task); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func updateTask(c echo.Context) error {
	task := models.Task{}
	if result := models.DB.First(&task, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if err := c.Bind(&task); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&task); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Save(&task); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func deleteTask(c echo.Context) error {
	task := models.Task{}
	if res := models.DB.Delete(&task, c.Param("id")); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func getTasksHist(c echo.Context) error {
	var total int64
	taskhists := models.TaskHistorys{}
	if res := models.DB.Model(&models.Task{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c)).Find(&taskhists); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": taskhists, "total": total})
}

func invokeTask(c echo.Context) error {
	task := models.Task{}
	if result := models.DB.First(&task, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	cmd := exec.Command("/bin/bash", "-c", task.Command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	username := c.Get("username").(string)
	fromip := c.RealIP()
	go func(stdout io.ReadCloser, user string, ip string) {
		now := time.Now()
		time_str, _ := times.GetTimeString(0, times.TimeOptions{Layout: "20060102150405"})
		hist_dir := "./taskhist/"
		hist_file := fmt.Sprintf("%s%s-%s.cast", hist_dir, username, time_str)
		hist := models.TaskHistory{TaskName: task.Name, User: user, FromIP: ip, File: hist_file}
		models.DB.Create(&hist)
		var f *os.File
		f, _ = os.Create(hist_file)
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
				misc.Logger.Error().Err(err).Str("from", "task").Msg("task stdout read error")
				return
			}
			// 记录终端内容
			nt, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float64(time.Now().UnixNano())/float64(1e9)), 64)
			iodata := []string{fmt.Sprintf("%.9f", nt-timestamp), "o", string(buf[:n])}
			iodatabyte, _ := json.Marshal(iodata)
			f.WriteString(string(iodatabyte) + "\n")
			// // 输出到 websocket
			// err = ws.WriteMessage(websocket.TextMessage, buf[:n])
			// if err != nil {
			// 	misc.Logger.Error().Err(err).Str("from", "terminal").Msg("write ws msg error")
			// 	exitCh <- true
			// 	return
			// }
		}
	}(stdout, username, fromip)
	return c.JSON(200, echo.Map{"success": true})
}
