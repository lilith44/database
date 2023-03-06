package database

import (
	"fmt"
	"time"
)

type (
	// ConnectionConfig 数据库连接池配置
	ConnectionConfig struct {
		// 最大打开连接数
		MaxOpen int `default:"150" yaml:"maxOpen"`
		// 最大休眠连接数
		MaxIdle int `default:"30" yaml:"maxIdle"`
		// 每个连接最大存活时间
		MaxLifetime time.Duration `default:"5s" yaml:"maxLifetime"`
	}

	// Config 数据库配置
	Config struct {
		// 类型
		Type string `default:"MYSQL" yaml:"type" validate:"required"`
		// 地址，带端口
		Address string `default:"127.0.0.1:3306" yaml:"address" validate:"required"`
		// 用户名
		Username string `yaml:"username"`
		// 密码
		Password string `yaml:"password"`
		// 协议
		Protocol string `default:"tcp" yaml:"protocol"`

		// 数据库名
		Schema string `yaml:"schema"`
		// 额外参数
		Parameters string `yaml:"parameters"`
		// 是否打印SQL语句
		ShowSQL bool `yaml:"showSQL"`

		// 连接池配置
		Connection ConnectionConfig `yaml:"connection"`
	}
)

func (c Config) DSN() string {
	dsn := fmt.Sprintf("%s:%s@%s(%s)", c.Username, c.Password, c.Protocol, c.Address)
	if c.Schema != "" {
		dsn = fmt.Sprintf("%s/%s", dsn, c.Schema)
	}
	if c.Parameters != "" {
		dsn = fmt.Sprintf("%s?%s", dsn, c.Parameters)
	}

	return dsn
}
