/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:03
 * @Description:
 */

package config

type Config struct {
	DSN   string
	Redis Redis
}

type (
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
)

func NewConfig() (*Config, error) {

	return &Config{}, nil
}
