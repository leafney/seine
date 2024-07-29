/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-29 17:06
 * @Description:
 */

package core

import (
	"github.com/go-co-op/gocron"
	"github.com/leafney/seine/global"
	"time"
)

func InitCron(stop chan struct{}) {
	cron := gocron.NewScheduler(time.Local)

	// -------------------

	//cron.
	//	Cron("*/5 * * * *").
	//	SingletonMode().
	//	WaitForSchedule().
	//	Do(service.Project.LoadAllProject)

	//cron.Every(1).Day().At("11:09").Do(func() {
	//	log.Println("hello")
	//})

	//cron.
	//	Cron("*/3 * * * *").
	//	SingletonMode().
	//	WaitForSchedule().
	//	Do(dao.Node.RunTask, int64(3), "9527")

	//job.LoadJobList(cron)

	// -------------------

	go func() {
		<-stop // 等待停止信号
		cron.Stop()
		global.GXLog.Infoln("[Cron] Exit successful")
	}()

	// 异步启动
	cron.StartAsync()
	global.GXLog.Infoln("[Cron] Load successful")
}
