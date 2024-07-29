/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-29 17:00
 * @Description:
 */

package core

import (
	"github.com/leafney/rose-notify/notify"
	"github.com/leafney/rose-notify/wochat"
	"github.com/leafney/seine/global"
)

func InitNotify() {
	cfgNotifyWoChatToken := global.GConfig.Notify.WoChatToken
	// 企业微信
	woChat := wochat.NewWoChat(cfgNotifyWoChatToken)
	n := notify.NewNotify(woChat)

	global.GNotify = n
	global.GXLog.Infoln("[Notify] Load successful")
}
