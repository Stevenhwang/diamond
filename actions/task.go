package actions

import (
	"diamond/models"
	"fmt"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
)

func getTasks(c echo.Context) error {
	var total int64
	tasks := models.Tasks{}
	if res := models.DB.Scopes(models.AnyFilter(models.Task{}, c)).Model(&models.Task{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c), models.AnyFilter(models.Task{}, c)).Find(&tasks); res.Error != nil {
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
	if res := models.DB.Model(&models.TaskHistory{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c)).Order("created_at desc").Omit("content").Find(&taskhists); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": taskhists, "total": total})
}

func getTasksHistDetail(c echo.Context) error {
	taskHist := models.TaskHistory{}
	if result := models.DB.First(&taskHist, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": taskHist})
}

func invokeTask(c echo.Context) error {
	task := models.Task{}
	if res := models.DB.First(&task, c.Param("id")); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	script := models.Script{}
	if res := models.DB.First(&script, task.ScriptID); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	// create temp script file
	f, err := os.CreateTemp("", "tempscript")
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if _, err := f.WriteString(script.Content); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// ansible localhost -m script -a "/tmp/test.sh arg1 arg2"
	scriptArgs := fmt.Sprintf("%s %s", f.Name(), task.Args)
	cmdArgs := []string{task.Target, "-m", "script", "-a", scriptArgs}
	cmd := exec.Command("ansible", cmdArgs...)
	username := c.Get("username").(string)
	fromip := c.RealIP()
	go func() {
		defer os.Remove(f.Name()) // ensure temp script file is deleted
		output, err := cmd.CombinedOutput()
		var (
			success bool
			content string
		)
		if err != nil {
			success = false
			content = string(output) + "\n" + err.Error()
		} else {
			success = true
			content = string(output)
		}
		hist := models.TaskHistory{TaskName: task.Name, User: username, FromIP: fromip, Success: success, Content: content}
		models.DB.Create(&hist)
	}()
	return c.JSON(200, echo.Map{"success": true})
}
