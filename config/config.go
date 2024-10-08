/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:03
 * @Description:
 */

package config

type Config struct {
	DSN     string
	Redis   Redis
	Cache   Cache
	Log     Log
	LevelDB LevelDB
	Mongo   Mongo
}

type (
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}

	Cache struct {
		Enable  bool
		Minutes int64
	}

	Log struct {
		XEnable bool
		XDebug  bool
		XLevel  string
		ZEnable bool
		ZLevel  string
	}

	LevelDB struct {
		Path string
	}

	Mongo struct {
		Addr string
		DB   string
	}
)

func NewConfig() (*Config, error) {

	return &Config{}, nil
}
