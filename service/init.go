package service

import (
	"datahub-external/service/config"
	"datahub-external/service/proxy"
	"datahub-external/service/proxy/lvyun"
	"datahub-external/service/proxy/ps"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	// 设置时区为 Asia/Shanghai
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		slog.Warn("加载时区失败,使用系统默认时区", "error", err)
	} else {
		time.Local = loc
		slog.Info("时区设置成功", "timezone", "Asia/Shanghai")
	}

	// 加载配置
	cfg := config.LoadConfig()

	// 配置 slog，使用本地时区输出时间
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 将时间字段转换为本地时区
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.Attr{
						Key:   a.Key,
						Value: slog.StringValue(t.Local().Format(time.RFC3339)),
					}
				}
			}
			return a
		},
	})))

	// 初始化数据库连接
	var db *gorm.DB
	if cfg.Database.Host != "" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			slog.Error("数据库连接失败", "error", err)
		} else {
			slog.Info("数据库连接成功",
				"host", cfg.Database.Host,
				"database", cfg.Database.DBName)
		}
	} else {
		slog.Warn("数据库配置未设置,将不使用数据持久化功能")
	}

	// 初始化绿云客户端(如果配置了)
	if cfg.Lvyun.BaseURL != "" {
		lvyunClient := lvyun.NewLvyunClient(cfg.Lvyun, db)
		if err := proxy.GetGlobalRegistry().Register(lvyunClient); err != nil {
			slog.Error("注册绿云客户端失败", "error", err)
		} else {
			slog.Info("绿云客户端注册成功", "base_url", cfg.Lvyun.BaseURL)
		}
	} else {
		slog.Warn("绿云配置未设置,跳过初始化")
	}

	// 初始化PS客户端(如果配置了)
	if cfg.PS.AppKey != "" && cfg.PS.AppSecret != "" {
		psConfig := ps.PSConfig{
			AppKey:           cfg.PS.AppKey,
			AppSecret:        cfg.PS.AppSecret,
			Stage:            cfg.PS.Stage,
			BaseURL:          cfg.PS.BaseURL,
			EnableScheduler:  cfg.PS.EnableScheduler,
			FamilyMemberCron: cfg.PS.FamilyMemberCron,
			PageSize:         cfg.PS.PageSize,
			MaxPages:         cfg.PS.MaxPages,
		}
		psClient := ps.NewPSClientWithConfig(psConfig, db)
		if err := proxy.GetGlobalRegistry().Register(psClient); err != nil {
			slog.Error("注册PS客户端失败", "error", err)
		} else {
			slog.Info("PS客户端注册成功",
				"base_url", cfg.PS.BaseURL,
				"stage", cfg.PS.Stage,
				"scheduler_enabled", cfg.PS.EnableScheduler)
		}
	} else {
		slog.Warn("PS配置未设置,跳过初始化")
	}
}
