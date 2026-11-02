package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	SD       SDConfig       `yaml:"sd"`
	Upload   UploadConfig   `yaml:"upload"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug | release | test
}

type LoggerConfig struct {
	Level     string `yaml:"level"`
	Format    string `yaml:"format"`
	AddSource bool   `yaml:"add_source"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SearchPath      string `yaml:"search_path"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

func (d *DatabaseConfig) DSN() string {
	if d.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			d.Host, d.Port, d.User, d.DBName)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.DBName)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Addr 返回 Redis 地址
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret          string `yaml:"secret"`
	AccessTokenTTL  int    `yaml:"access_token_ttl"`
	RefreshTokenTTL int    `yaml:"refresh_token_ttl"`
}

// SDConfig Stable Diffusion 配置
type SDConfig struct {
	Provider       string `yaml:"provider"`
	BaseURL        string `yaml:"base_url"`        // A1111 地址
	APIKey         string `yaml:"api_key"`         // Replicate API Key
	MaxConcurrent  int    `yaml:"max_concurrent"`  // 最大并发数
	RequestTimeout int    `yaml:"request_timeout"` // 请求超时(秒)

	MockImageDir string `yaml:"mock_image_dir"`
	MockDelayMs  int    `yaml:"mock_delay_ms"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Dir     string `yaml:"dir"`      // 图片存储目录
	MaxSize int64  `yaml:"max_size"` // 最大文件大小(MB)
}

var once sync.Once
var config = &Config{}

// Load 加载配置
func Load(path string) (*Config, error) {
	once.Do(func() {
		conf, err := os.ReadFile(path)
		if err != nil {
			log.Fatalln("read config file error:", err)
		}
		if err := yaml.Unmarshal(conf, config); err != nil {
			log.Fatalln("unmarshal config file error:", err)
		}
	})
	return config, nil
}

func Temple() []byte {
	conf := Config{}
	out, err := yaml.Marshal(conf)
	if err != nil {
		log.Printf("marshal config file error: %s", err)
		return nil
	}
	return out
}
