package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	ListenPort int
	ApiKey     string
	Lvyun      LvyunConfig
	Database   DatabaseConfig
}

// LvyunConfig 绿云配置
type LvyunConfig struct {
	BaseURL            string
	HotelGroupCode     string
	HotelCode          string
	UserCode           string
	Password           string
	AppKey             string
	AppSecret          string
	EnableScheduler    bool
	ReservationCron    string
	RegistrationCron   string
	CheckoutCron       string
	BusinessReportCron string
	QueryDays          int
	BusinessReportDays int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		ApiKey: getEnv("API_KEY", "default-api-key-please-change-in-production"),
		Lvyun: LvyunConfig{
			BaseURL:            getEnv("LVYUN_BASE_URL", ""),
			HotelGroupCode:     getEnv("LVYUN_HOTEL_GROUP_CODE", ""),
			HotelCode:          getEnv("LVYUN_HOTEL_CODE", ""),
			UserCode:           getEnv("LVYUN_USER_CODE", ""),
			Password:           getEnv("LVYUN_PASSWORD", ""),
			AppKey:             getEnv("LVYUN_APP_KEY", ""),
			AppSecret:          getEnv("LVYUN_APP_SECRET", ""),
			EnableScheduler:    getEnvBool("LVYUN_ENABLE_SCHEDULER", true),
			ReservationCron:    getEnv("LVYUN_RESERVATION_CRON", "0 */30 * * * *"),    // 每30分钟
			RegistrationCron:   getEnv("LVYUN_REGISTRATION_CRON", "0 */30 * * * *"),   // 每30分钟
			CheckoutCron:       getEnv("LVYUN_CHECKOUT_CRON", "0 */30 * * * *"),       // 每30分钟
			BusinessReportCron: getEnv("LVYUN_BUSINESS_REPORT_CRON", "0 */30 * * * *"), // 每30分钟
			QueryDays:          getEnvInt("LVYUN_QUERY_DAYS", 7),
			BusinessReportDays: getEnvInt("LVYUN_BUSINESS_REPORT_DAYS", 1),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "datahub_external"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
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

// getEnvInt 获取整型环境变量
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvBool 获取布尔型环境变量
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}



