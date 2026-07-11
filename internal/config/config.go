/*
Copyright 2026 stonebirdjx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/stonebirdjx/go-layout/pkg/consts"
)

// Config represents the global configuration structure.
type Config struct {
	Log LogOptions `yaml:"log"`
}

type LogOptions struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	// 日志文件路径，为空则只输出到 stdout
	FilePath string `yaml:"file_path"`
	// 单个日志文件最大大小（MB）
	MaxSize int `yaml:"max_size"`
	// 保留旧日志文件最大数量
	MaxBackups int `yaml:"max_backups"`
	// 保留旧日志文件最大天数
	MaxAge int `yaml:"max_age"`
	// 是否压缩旧日志文件
	Compress bool `yaml:"compress"`
}

// Validate checks all configuration fields.
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if err := c.Log.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate checks that Level, Format and rotation parameters are valid.
func (l LogOptions) Validate() error {
	switch l.Level {
	case consts.LoggerLevelDebug, consts.LoggerLevelInfo, consts.LoggerLevelWarn, consts.LoggerLevelError:
	default:
		return fmt.Errorf("invalid log level: %q (must be one of: %s, %s, %s, %s)",
			l.Level, consts.LoggerLevelDebug, consts.LoggerLevelInfo, consts.LoggerLevelWarn, consts.LoggerLevelError)
	}

	switch l.Format {
	case consts.LoggerFormatJSON, consts.LoggerFormatConsole:
	default:
		return fmt.Errorf("invalid log format: %q (must be one of: %s, %s)",
			l.Format, consts.LoggerFormatJSON, consts.LoggerFormatConsole)
	}

	// 当指定了日志文件路径时，校验轮转参数
	if l.FilePath != "" {
		if l.MaxSize <= 0 {
			return fmt.Errorf("invalid log max_size: %d (must be > 0 when file_path is set)", l.MaxSize)
		}
		if l.MaxBackups < 0 {
			return fmt.Errorf("invalid log max_backups: %d (must be >= 0)", l.MaxBackups)
		}
		if l.MaxAge < 0 {
			return fmt.Errorf("invalid log max_age: %d (must be >= 0)", l.MaxAge)
		}
	}

	return nil
}

// Load loads the configuration from the given file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 不推荐写法
	// 避免心智负担：
	// 如果写成 yaml.Unmarshal(data, cfg)，读者在看到这一行时，
	// 无法直接确认 cfg 会不会被修改，
	// 必须向上回看 cfg 的定义才能知道它是一个指针（*Config）还是一个普通结构体。
	// 这种写法略微增加了阅读代码时的心智负担。

	// var cfg *Config
	// if err := yaml.Unmarshal(data, &cfg); err != nil {
	// 	return nil, err
	// }

	// 推荐写法
	// 调用处的“显式可变性”
	// Go 社区及标准库（如 json.Unmarshal）
	// 在绝大多数情况下都更推崇 “声明值类型变量，并在调用时显式取地址 & 传递” 的写法。
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
