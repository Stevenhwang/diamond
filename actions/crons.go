package actions

// import (
// 	"context"
// 	"sync"
// 	"diamond/misc"
// 	"diamond/models"
// 	"time"

// 	"github.com/Stevenhwang/gommon/ssh"
// 	"github.com/Stevenhwang/gommon/tools"
// 	"github.com/robfig/cron/v3"
// 	gossh "golang.org/x/crypto/ssh"
// )

// func init() {
// 	crons := cron.New(cron.WithSeconds(),
// 		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger),
// 			cron.Recover(cron.DefaultLogger)))
// 	crons.AddFunc("@every 3600s", gatherTask)
// 	crons.Start()
// }

// type remoteRes struct {
// 	// Cpu          int
// 	// Memory       int
// 	// Disk         int
// 	InstanceType string
// 	Err          error
// }

// func remoteSSH(server models.Server) remoteRes {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5秒超时
// 	defer cancel()
// 	resChan := make(chan remoteRes)
// 	go func(ctx context.Context, ch chan remoteRes) {
// 		credential := models.Credential{}
// 		if res := models.DB.First(&credential, server.CredentialID); res.Error != nil {
// 			ch <- remoteRes{Err: res.Error}
// 			return
// 		}
// 		var client *gossh.Client
// 		spass := tools.AesDecrypt(credential.AuthContent, "0123456789012345")
// 		if credential.AuthType == 1 {
// 			cli, err := ssh.GetSSHClientByPassword(server.IP, spass, ssh.SSHOptions{Port: server.Port, User: credential.AuthUser})
// 			if err != nil {
// 				ch <- remoteRes{Err: err}
// 				return
// 			}
// 			client = cli
// 		}
// 		if credential.AuthType == 2 {
// 			cli, err := ssh.GetSSHClientByKey(server.IP, []byte(spass), ssh.SSHOptions{Port: server.Port, User: credential.AuthUser})
// 			if err != nil {
// 				ch <- remoteRes{Err: err}
// 				return
// 			}
// 			client = cli
// 		}
// 		defer client.Close()
// 		// cpustr, err := ssh.SSHExec(`echo $[100-$(vmstat 1 2|tail -1|awk '{print $15}')]`, client)
// 		// if err != nil {
// 		// 	ch <- remoteRes{Err: err}
// 		// 	return
// 		// }
// 		// memstr, err := ssh.SSHExec(`echo $(free | grep Mem | awk '{printf "%d", $3/$2 * 100.0}')`, client)
// 		// if err != nil {
// 		// 	ch <- remoteRes{Err: err}
// 		// 	return
// 		// }
// 		// diskstr, err := ssh.SSHExec(`df -Th / | tail -1 | awk '{print $6}' | sed s/%//g`, client)
// 		// if err != nil {
// 		// 	ch <- remoteRes{Err: err}
// 		// 	return
// 		// }
// 		instance_type, err := ssh.SSHExec("curl -s http://169.254.169.254/latest/meta-data/instance-type", client)
// 		if err != nil {
// 			ch <- remoteRes{Err: err}
// 			return
// 		}
// 		// cpu, _ := strconv.Atoi(strings.TrimRight(cpustr, "\n"))
// 		// mem, _ := strconv.Atoi(strings.TrimRight(memstr, "\n"))
// 		// disk, _ := strconv.Atoi(strings.TrimRight(diskstr, "\n"))
// 		res := remoteRes{InstanceType: instance_type, Err: nil}
// 		ch <- res
// 	}(ctx, resChan)
// 	// 检测channel
// 	select {
// 	case x := <-resChan:
// 		return x
// 	case <-ctx.Done():
// 		return remoteRes{Err: ctx.Err()}
// 	}
// }

// func gatherTask() {
// 	servers := models.Servers{}
// 	if res := models.DB.Find(&servers); res.Error != nil {
// 		misc.Logger.Error().Err(res.Error).Str("from", "crons").Msg("获取服务器列表失败")
// 		return
// 	}
// 	var wg sync.WaitGroup
// 	for _, server := range servers {
// 		wg.Add(1)
// 		x := server
// 		go func() {
// 			defer wg.Done()
// 			res := remoteSSH(x)
// 			if res.Err != nil {
// 				x.InstanceType = res.Err.Error()
// 			} else {
// 				x.InstanceType = res.InstanceType
// 			}
// 			models.DB.Save(&x)
// 		}()
// 	}
// 	wg.Wait()
// }
