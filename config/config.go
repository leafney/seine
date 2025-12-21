package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Cache    CacheConfig    `yaml:"cache"`
	JWT      JWTConfig      `yaml:"jwt"`
	Crypto   CryptoConfig   `yaml:"crypto"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	XEnable bool   `yaml:"xEnable"` // xlog 是否启用
	XDebug  bool   `yaml:"xDebug"`  // xlog 调试模式
	XLevel  string `yaml:"xLevel"`  // xlog 日志级别
	ZEnable bool   `yaml:"zEnable"` // zap 是否启用
	ZLevel  string `yaml:"zLevel"`  // zap 日志级别
	ZCaller bool   `yaml:"zCaller"` // zap 是否显示调用者
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string       `yaml:"driver"` // mysql 或 sqlite
	MySQL  MySQLConfig  `yaml:"mysql"`
	SQLite SQLiteConfig `yaml:"sqlite"`
	Debug  bool         `yaml:"debug"`
}

type MySQLConfig struct {
	DSN      string `yaml:"dsn"`
	IdleOpen int    `yaml:"idleOpen"`
	MaxOpen  int    `yaml:"maxOpen"`
}

type SQLiteConfig struct {
	DSN string `yaml:"dsn"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	DB       int    `yaml:"db"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Minutes int64 `yaml:"minutes"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	SigningKey        string `yaml:"signingKey"`
	AccessTokenExpire int    `yaml:"accessTokenExpire"` // 秒
}

// CryptoConfig 加密配置
type CryptoConfig struct {
	Key string `yaml:"key"` // AES-256 密钥（32字节）
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// 实现 configx.CacheConfig 接口
func (c CacheConfig) GetCacheMinutes() int64 {
	return c.Minutes
}

// 实现 configx.LogConfig 接口
func (c LogConfig) GetLogLevel() string {
	return c.ZLevel
}

func (c LogConfig) GetLogFormat() string {
	return "json"
}

func (c LogConfig) IsDebugMode() bool {
	return c.XDebug
}

// 实现 configx.RedisConfig 接口
func (c RedisConfig) GetRedisAddr() string {
	return c.Host
}

func (c RedisConfig) GetRedisPassword() string {
	return c.Password
}

func (c RedisConfig) GetRedisDB() int {
	return c.DB
}

// 实现 configx.JWTConfig 接口
func (c JWTConfig) GetJWTSigningKey() string {
	return c.SigningKey
}

func (c JWTConfig) GetJWTExpireHours() int {
	return c.AccessTokenExpire / 3600
}

// 实现 configx.BadgerDBConfig 接口
func (c DatabaseConfig) GetBadgerDBPath() string {
	return "data/.cache/badger"
}

// 实现 configx.SQLiteConfig 接口
func (c DatabaseConfig) GetSQLiteDatabase() string {
	return c.SQLite.DSN
}
