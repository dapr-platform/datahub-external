package config

import "os"

// Config 应用配置
type Config struct {
	ListenPort int
	ApiKey     string
	Lvyun      LvyunConfig
}

// LvyunConfig 绿云配置
type LvyunConfig struct {
	BaseURL        string
	HotelGroupCode string
	UserCode       string
	Password       string
	AppKey         string
	AppSecret      string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		ApiKey: getEnv("API_KEY", "default-api-key-please-change-in-production"),
		Lvyun: LvyunConfig{
			BaseURL:        getEnv("LVYUN_BASE_URL", ""),
			HotelGroupCode: getEnv("LVYUN_HOTEL_GROUP_CODE", ""),
			UserCode:       getEnv("LVYUN_USER_CODE", ""),
			Password:       getEnv("LVYUN_PASSWORD", ""),
			AppKey:         getEnv("LVYUN_APP_KEY", ""),
			AppSecret:      getEnv("LVYUN_APP_SECRET", ""),
		},
	}
}

// getEnv 获取环境变量,如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}



