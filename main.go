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

	app, callback, err := cmd.InitApp()
	if err != nil {
		log.Fatalln(err)
	}

	defer callback()

	app.Start(quitChan)
}
