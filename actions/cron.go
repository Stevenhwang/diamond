package actions

import (
	"diamond/misc"
	"diamond/models"
	"fmt"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
	ocron "github.com/robfig/cron/v3"
)

var mycron *ocron.Cron

func init() {
	mycron = ocron.New(ocron.WithSeconds(),
		ocron.WithChain(ocron.SkipIfStillRunning(ocron.DefaultLogger),
			ocron.Recover(ocron.DefaultLogger)))
}

func CronStart() {
	// 初始化的时候如果数据库有定时任务，要添加进去，同时要更新entryid
	crons := models.Crons{}
	if result := models.DB.Find(&crons); result.Error != nil {
		misc.Logger.Error().Err(result.Error).Str("from", "cron").Msg("")
	}
	for _, c := range crons {
		// 添加定时任务
		entryid, _ := mycron.AddFunc(c.Spec, execCron(c.Name, c.Target, c.ScriptID, c.Args))
		// 更新entryid
		c.EntryID = int(entryid)
		models.DB.Save(&c)
	}
	mycron.Start()
}

// cron 执行函数
func execCron(name string, target string, scriptID uint, args string) func() {
	return func() {
		script := models.Script{}
		if res := models.DB.First(&script, scriptID); res.Error != nil {
			misc.Logger.Error().Err(res.Error).Str("from", "cron").Msg("find script error")
			return
		}
		// create temp script file
		f, err := os.CreateTemp("", "crontempscript")
		if err != nil {
			misc.Logger.Error().Err(err).Str("from", "cron").Msg("create temp file error")
			return
		}
		defer os.Remove(f.Name()) // ensure temp script file is deleted
		scriptArgs := fmt.Sprintf("%s %s", f.Name(), args)
		cmdArgs := []string{target, "-m", "script", "-a", scriptArgs}
		cmd := exec.Command("ansible", cmdArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			msg := fmt.Sprintf("%s任务执行失败", name)
			misc.Logger.Error().Err(err).Str("from", "cron").Msg(msg)
		} else {
			msg := fmt.Sprintf("%s任务执行成功", name)
			misc.Logger.Info().Str("from", "cron").Msg(msg)
		}
		misc.Logger.Info().Str("from", "cron").Msg(string(output))
	}
}

/*
******************
定时任务actions
******************
*/
func createCron(c echo.Context) error {
	cron := models.Cron{}
	if err := c.Bind(&cron); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&cron); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	entryid, err := mycron.AddFunc(cron.Spec, execCron(cron.Name, cron.Target, cron.ScriptID, cron.Args))
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	cron.EntryID = int(entryid)
	if res := models.DB.Create(&cron); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func updateCron(c echo.Context) error {
	cron := models.Cron{}
	if result := models.DB.First(&cron, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if err := c.Bind(&cron); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&cron); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 先删除cron
	mycron.Remove(ocron.EntryID(cron.EntryID))
	// 再添加cron
	entryid, err := mycron.AddFunc(cron.Spec, execCron(cron.Name, cron.Target, cron.ScriptID, cron.Args))
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 更新entryid和其他数据到数据库
	cron.EntryID = int(entryid)
	if result := models.DB.Save(&cron); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func deleteCron(c echo.Context) error {
	cron := models.Cron{}
	if result := models.DB.First(&cron, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	mycron.Remove(ocron.EntryID(cron.EntryID))
	if res := models.DB.Delete(&cron); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func getCrons(c echo.Context) error {
	var total int64
	crons := models.Crons{}
	if res := models.DB.Model(&models.Cron{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c), models.AnyFilter(models.Cron{}, c)).Find(&crons); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": crons, "total": total})
}

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
