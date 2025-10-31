package service

import (
	"datahub-external/service/config"
	"datahub-external/service/proxy"
	"datahub-external/service/proxy/lvyun"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	// 加载配置
	cfg := config.LoadConfig()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
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
}
