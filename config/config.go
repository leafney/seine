/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-29 16:53
 * @Description:
 */

package config

import "embed"

//go:embed config-default.yaml
var DefaultConfig embed.FS

type Config struct {
	Port    string  `koanf:"port"`
	Mongo   Mongo   `koanf:"mongo"`
	Redis   Redis   `koanf:"redis"`
	LevelDB string  `koanf:"leveldb"`
	Jwt     JWT     `koanf:"jwt"`
	Captcha Captcha `koanf:"captcha"`
	Email   Email   `koanf:"email"`
	Notify  Notify  `koanf:"notify"`
}

type (
	Mongo struct {
		Addr string `koanf:"addr"` // mongodb://user:password@127.0.0.1:27017/?authSource=admin
		Db   string `koanf:"db"`
	}

	Redis struct {
		Addr string `koanf:"addr"`
		Pwd  string `koanf:"pwd"`
		Db   int    `koanf:"db"`
	}

	JWT struct {
		SigningKey       string `koanf:"signing_key"`
		LoginTokenExpire int64  `koanf:"login_token_expire"`
		ApiTokenExpire   int64  `koanf:"api_token_expire"`
		ApiTokenEncode   bool   `koanf:"api_token_encode"`
	}

	Captcha struct {
		Debug     bool   `koanf:"debug"`      // 调试模式
		DebugCode string `koanf:"debug_code"` // 调试验证码
		MinCode   int    `koanf:"min_code"`   // 最小值
		MaxCode   int    `koanf:"max_code"`   // 最大值
		ExpireSec int64  `koanf:"expire_sec"` // 验证码过期时间
		DelaySec  int64  `koanf:"delay_sec"`  // 获取验证码间隔时间
		DayTimes  int    `koanf:"day_times"`  // 每日次数限制
	}

	Email struct {
		Enable     bool     `koanf:"enable"`      // 可用状态
		SuperMode  bool     `koanf:"super_mode"`  // 特权模式
		SuperEmail []string `koanf:"super_email"` // 特权邮箱列表
		Host       string   `koanf:"host"`        // smtp地址
		Port       int      `koanf:"port"`        // 端口
		NickName   string   `koanf:"nick_name"`   // 发送方昵称
		UserName   string   `koanf:"user_name"`   // 用户名
		PassWord   string   `koanf:"pass_word"`   // 密码
	}

	Notify struct {
		WoChatToken string `koanf:"wochat_token"`
	}
)
