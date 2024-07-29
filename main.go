/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 11:43
 * @Description:
 */

package main

import (
	"fmt"
	"github.com/leafney/seine/cmd/core"
	"github.com/leafney/seine/cmd/run"
	flag "github.com/spf13/pflag"
	"runtime"
)

var (
	v         bool
	h         bool
	Version   = "v0.1.0"
	GitBranch = ""
	GitCommit = ""
	BuildTime = "2024-07-29 17:18:30"
)

func main() {
	flag.BoolVarP(&h, "help", "h", false, "help")
	flag.BoolVarP(&v, "version", "v", false, "version")
	flag.Parse()

	if h {
		flag.PrintDefaults()
	} else if v {
		// 输出版本信息
		fmt.Println("Version:      " + Version)
		fmt.Println("Git branch:   " + GitBranch)
		fmt.Println("Git commit:   " + GitCommit)
		fmt.Println("Built time:   " + BuildTime)
		fmt.Println("Go version:   " + runtime.Version())
		fmt.Println("OS/Arch:      " + runtime.GOOS + "/" + runtime.GOARCH)
	} else {
		// 基础服务
		core.InitXLog()
		core.InitConfig()
		//core.InitMsqQueue()

		// 用于退出的通道
		quitChan := make(chan struct{})
		// 相关服务
		core.InitZLog(quitChan)
		//core.InitMongo(quitChan)
		//core.InitRedis(quitChan)
		//core.InitCron(quitChan)
		//core.InitCache(quitChan)
		//core.InitNotify()
		// 异步任务
		//core.InitQueue()
		// web服务
		run.Start(quitChan)
	}
}
