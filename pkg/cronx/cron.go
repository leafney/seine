/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 20:18
 * @Description:
 */

package cronx

import (
	"github.com/go-co-op/gocron"
	"github.com/leafney/seine/config"
	"log"
	"time"
)

type CronSvc struct {
	*gocron.Scheduler
}

func NewCronSvc(cfg *config.Config, stop chan struct{}) *CronSvc {

	cron := gocron.NewScheduler(time.Local)

	go func() {
		<-stop // 等待停止信号
		cron.Stop()
		log.Println("[Cron] Exit successful")
	}()

	cron.StartAsync()

	log.Println("[Cron] Load successful")

	return &CronSvc{cron}
}
