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
)

// Config represents the global configuration structure.
type Config struct {
	Log LogOptions `yaml:"log"`
}

type LogOptions struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
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

type logLevel string

const (
	logLevelDebug logLevel = "debug"
	logLevelInfo  logLevel = "info"
	logLevelWarn  logLevel = "warn"
	logLevelError logLevel = "error"
)

type logFormat string

const (
	logFormatJSON    logFormat = "json"
	logFormatConsole logFormat = "console"
)

// Validate checks that Level and Format are valid values.
func (l LogOptions) Validate() error {
	switch logLevel(l.Level) {
	case logLevelDebug, logLevelInfo, logLevelWarn, logLevelError:
	default:
		return fmt.Errorf("invalid log level: %q (must be one of: %s, %s, %s, %s)",
			l.Level, logLevelDebug, logLevelInfo, logLevelWarn, logLevelError)
	}

	switch logFormat(l.Format) {
	case logFormatJSON, logFormatConsole:
	default:
		return fmt.Errorf("invalid log format: %q (must be one of: %s, %s)",
			l.Format, logFormatJSON, logFormatConsole)
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
