package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AKAama/go-skeleton/config"
	"github.com/AKAama/go-skeleton/pkg/skeleton"
	"github.com/AKAama/go-skeleton/pkg/util"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "go-skeleton",
		Short: "go项目骨架生成器",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
			HiddenDefaultCmd:    true,
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       util.GetVersion().Version,
	}

	//--project api-backend -d $GOPATH$/src/go-awesome/test
	cmd.PersistentFlags().StringP("project", "p", "", "项目名称")
	_ = cmd.MarkPersistentFlagRequired("project")
	cmd.PersistentFlags().StringP("dir", "d", "", "项目所在文件夹路径 (default 当前路径)")

	cmd.PersistentFlags().StringSliceP("without", "w", []string{}, fmt.Sprintf("不启用%s框架", strings.Join(lo.Keys[string](config.DefaultModules), ",")))

	cmd.PersistentFlags().StringP("mod", "m", "", ".mod 文件的module路径 (default 项目名称)")
	_ = viper.BindPFlags(cmd.PersistentFlags())
	return cmd
}

func run() error {
	cfg, err := flagsHandler()
	if err != nil {
		return fmt.Errorf("参数处理错误: %w", err)
	}
	// 1. 检查是否安装go程序
	if err := skeleton.ShellGo(); err != nil {
		return fmt.Errorf("无法获取 go 执行命令: %w", err)
	}

	// 2. 检查创建项目目录，不是空目录，清除警告
	if err := skeleton.CreateProject(cfg); err != nil {
		return fmt.Errorf("创建项目目录错误: %w", err)
	}

	// 3. 执行go mod init 命令
	if err := skeleton.ShellModInit(cfg); err != nil {
		return fmt.Errorf("创建 go.mod 错误: %w", err)
	}

	if err := skeleton.GoMainFile(cfg); err != nil {
		return fmt.Errorf("创建 main.go 错误: %w", err)
	}

	if err := skeleton.GoRootFile(cfg); err != nil {
		return fmt.Errorf("创建 root.go 错误: %w", err)
	}

	if _, ok := cfg.Modules["gin"]; ok {
		if err := skeleton.GoSignalFile(cfg); err != nil {
			return fmt.Errorf("创建 signal.go 错误: %w", err)
		}
	}

	if err := skeleton.GoGlobalConfigFile(cfg); err != nil {
		return fmt.Errorf("创建 global_config.go 错误: %w", err)
	}
	if err := skeleton.GoConfigYamlFile(cfg); err != nil {
		return fmt.Errorf("创建 config.yaml 错误: %w", err)
	}
	if _, ok := cfg.Modules["gorm"]; ok {
		if err := skeleton.GoDBConfigFile(cfg); err != nil {
			return fmt.Errorf("创建 db_config.go 错误: %w", err)
		}
		if err := skeleton.GoDBFile(cfg); err != nil {
			return fmt.Errorf("创建 database.go 错误: %w", err)
		}
	}
	if _, ok := cfg.Modules["gin"]; ok {
		if err := skeleton.GoHttpFile(cfg); err != nil {
			return fmt.Errorf("创建 http.go 错误: %w", err)
		}
		if err := skeleton.GoRouteFile(cfg); err != nil {
			return fmt.Errorf("创建 route.go 错误: %w", err)
		}
	}
	if err := skeleton.GoUtilFile(cfg); err != nil {
		return fmt.Errorf("创建 util.go 错误: %w", err)
	}
	if err := skeleton.GoVersionFile(cfg); err != nil {
		return fmt.Errorf("创建 version.go 错误: %w", err)
	}
	if err := skeleton.MakeFile(cfg); err != nil {
		return fmt.Errorf("创建 Makefile 错误: %w", err)
	}
	if err := skeleton.DockerFile(cfg); err != nil {
		return fmt.Errorf("创建 Dockerfile 错误: %w", err)
	}
	if err := skeleton.GitIgnore(cfg); err != nil {
		return fmt.Errorf("创建 .gitignore 错误: %w", err)
	}
	if err := skeleton.ShellPinDependencies(cfg); err != nil {
		return fmt.Errorf("固定项目依赖版本错误: %w", err)
	}
	if err := skeleton.ShellGoModTidy(cfg); err != nil {
		return fmt.Errorf("执行 go mod tidy 命令错误: %w", err)
	}
	return nil
}

func flagsHandler() (*config.ProjectConfig, error) {
	dir := viper.GetString("dir")
	if dir == "" {
		dir, _ = os.Getwd()
	}
	name := viper.GetString("project")
	modulePath := viper.GetString("mod")
	if modulePath == "" {
		modulePath = name
	}

	opts := []config.ConfigOption{
		config.WithProjectName(name),
		config.WithProjectDir(dir),
		config.WithModulePath(modulePath),
	}

	out := viper.GetStringSlice("without")
	if len(out) > 0 {
		opts = append(opts, config.WithOutModules(out...))
	}

	cfg := config.NewProjectConfig(opts...)
	if errs := cfg.Validate(); len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}
