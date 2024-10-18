/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 00:56
 * @Description:
 */

package main

import (
	"github.com/leafney/seine/cmd"
	"log"
)

func main() {
	// 用于退出的通道
	quitChan := make(chan struct{})

	injector, callback, err := cmd.BuildInjector(quitChan)
	if err != nil {
		log.Fatalln(err)
	}
	defer callback()

	if err := injector.R.Init(); err != nil {
		injector.L.Fatalf("初始化异常 %v", err)
	}

	cmd.StartHttpServer(injector, quitChan)
}
