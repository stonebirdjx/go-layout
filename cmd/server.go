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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stonebirdjx/go-layout/internal/config"
	"github.com/stonebirdjx/go-layout/internal/logger"
	"github.com/stonebirdjx/go-layout/internal/manager"
	"github.com/stonebirdjx/go-layout/pkg/signals"
	"go.uber.org/zap"
)

var cfgFile string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the backend application server",
	Long:  `Start the backend application server with lifecycle management and dependency injection.`,
	RunE:  runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// Here you will define your flags and configuration settings.
	serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "configs/server.yaml", "config file path")
}

// runServer runs the server.
func runServer(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v", err)
		return err
	}

	// 校验配置
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "validate config error: %v", err)
		return err
	}

	// 初始化日志
	zapLogger := logger.Init(cfg.Log)
	defer zapLogger.Sync()

	// 创建管理器
	mgr := manager.NewManager(
		manager.WithLogger(zapLogger),
	)

	// 校验管理器
	if err := mgr.Validate(); err != nil {
		zapLogger.Error("validate manager error", zap.Error(err))
		return err
	}

	// 启动
	signalCtx := signals.SetupSignalHandler()
	if err := mgr.Start(signalCtx); err != nil {
		zapLogger.Error("server failed to start", zap.Error(err))
		return err
	}

	return nil
}
