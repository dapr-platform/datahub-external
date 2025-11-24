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
	PS         PSConfig
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

// PSConfig PS系统配置
type PSConfig struct {
	BaseURL             string
	AppKey              string
	AppSecret           string
	Stage               string
	EnableScheduler     bool
	FamilyMemberCron    string
	PositionIncCron     string
	PositionAllCron     string
	OrganizationIncCron string
	OrganizationAllCron string
	EmployeeIncCron     string
	EmployeeAllCron     string
	EmployeeHonorCron   string
	FamilyMainCron      string
	PageSize            int
	MaxPages            int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	LogLevel string // 日志级别: silent, error, warn, info
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
			ReservationCron:    getEnv("LVYUN_RESERVATION_CRON", "0 */30 * * * *"),     // 每30分钟
			RegistrationCron:   getEnv("LVYUN_REGISTRATION_CRON", "0 */30 * * * *"),    // 每30分钟
			CheckoutCron:       getEnv("LVYUN_CHECKOUT_CRON", "0 */30 * * * *"),        // 每30分钟
			BusinessReportCron: getEnv("LVYUN_BUSINESS_REPORT_CRON", "0 */30 * * * *"), // 每30分钟
			QueryDays:          getEnvInt("LVYUN_QUERY_DAYS", 7),
			BusinessReportDays: getEnvInt("LVYUN_BUSINESS_REPORT_DAYS", 1),
		},
		PS: PSConfig{
			BaseURL:             getEnv("PS_BASE_URL", "https://datadisclose.hdlapis.com"),
			AppKey:              getEnv("PS_APP_KEY", ""),
			AppSecret:           getEnv("PS_APP_SECRET", ""),
			Stage:               getEnv("PS_STAGE", "RELEASE"),
			EnableScheduler:     getEnvBool("PS_ENABLE_SCHEDULER", true),
			FamilyMemberCron:    getEnv("PS_FAMILY_MEMBER_CRON", "0 0 2 * * *"),     // 每天凌晨2点
			PositionIncCron:     getEnv("PS_POSITION_INC_CRON", "0 30 4 * * *"),     // 每天凌晨4:30
			PositionAllCron:     getEnv("PS_POSITION_ALL_CRON", "0 30 3 * * 0"),     // 每周日凌晨3:30
			OrganizationIncCron: getEnv("PS_ORGANIZATION_INC_CRON", "0 30 4 * * *"), // 每天凌晨4:30
			OrganizationAllCron: getEnv("PS_ORGANIZATION_ALL_CRON", "0 30 3 * * 0"), // 每周日凌晨3:30
			EmployeeIncCron:     getEnv("PS_EMPLOYEE_INC_CRON", "0 30 4 * * *"),     // 每天凌晨4:30
			EmployeeAllCron:     getEnv("PS_EMPLOYEE_ALL_CRON", "0 30 3 * * 0"),     // 每周日凌晨3:30
			EmployeeHonorCron:   getEnv("PS_EMPLOYEE_HONOR_CRON", "0 30 3 * * 0"),   // 每周日凌晨3:30
			FamilyMainCron:      getEnv("PS_FAMILY_MAIN_CRON", "0 30 3 * * 0"),      // 每周日凌晨3:30
			PageSize:            getEnvInt("PS_PAGE_SIZE", 2000),                    // API最大支持2000
			MaxPages:            getEnvInt("PS_MAX_PAGES", 1000),                    // 最多查询1000页
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "datahub_external"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			LogLevel: getEnv("DB_LOG_LEVEL", "warn"), // 默认warn级别，不打印SQL
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
