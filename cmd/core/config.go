/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-07-28 20:41
 * @Description:
 */

package core

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/leafney/rose"
	"github.com/leafney/seine/config"
	"github.com/leafney/seine/global"
	"github.com/leafney/seine/global/vars"
	"io/fs"
	"os"
)

var k = koanf.New(".")

func InitConfig() {
	// 指定多个配置文件路径
	paths := []string{vars.DefConfigFile}

	// 如果配置文件不存在，则自动创建默认配置文件
	loadDefaultConfig(paths)

	for _, path := range paths {
		if err := k.Load(file.Provider(path), yaml.Parser()); err == nil {
			global.GXLog.Infof("[Koanf] Load config file [%v] success", path)
			break
		} else {
			global.GXLog.Errorf("[Koanf] Load config file [%v] error [%v]", path, err)
		}
	}

	if len(k.Keys()) == 0 {
		global.GXLog.Fatalf("[Koanf] Config file empty")
	}

	if err := k.Unmarshal("", &global.GConfig); err != nil {
		global.GXLog.Fatalf("[Koanf] Unmarshal config error [%v]", err)
	}

	global.GXLog.Infoln("[Koanf] Load successful")
}

// 检查配置文件是否存在，如果不存在则创建默认配置文件
func loadDefaultConfig(paths []string) {
	configPath := vars.DefConfigFile

	exist := false
	for _, path := range paths {
		if rose.FIsExist(path) {
			exist = true
			break
		}
	}

	// 初始化默认配置文件
	if !exist {
		// 保证配置文件所在目录存在
		if err := rose.DEnsurePathExist(configPath); err != nil {
			global.GXLog.Fatalf("[Koanf] Failed to create config directory: %v", err)
		}

		data, err := fs.ReadFile(config.DefaultConfig, "config-default.yaml")
		if err != nil {
			global.GXLog.Fatalf("[Koanf] Failed to read embedded config: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			global.GXLog.Fatalf("[Koanf] Failed to write default config file: %v", err)
		}

		global.GXLog.Infoln("[Koanf] Default config file created")
	}
}
