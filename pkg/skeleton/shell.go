package skeleton

import (
	"os/exec"
	"path"

	"github.com/AKAama/go-skeleton/config"
	"github.com/AKAama/go-skeleton/pkg/util"
	"github.com/pkg/errors"
)

// ShellGo 检查是否可执行go命令
func ShellGo() error {
	cmd := exec.Command("go", "version")
	return cmd.Run()
}

// CreateProject 创建项目目录
func CreateProject(cfg *config.ProjectConfig) error {
	cfg.ProjectDir = path.Join(cfg.ProjectDir, cfg.ProjectName)
	exist, err := util.DirExist(cfg.ProjectDir)
	if err != nil {
		return err
	}
	if exist {
		isEmpty, err := util.DirIsEmpty(cfg.ProjectDir)
		if err != nil {
			return err
		}
		if !isEmpty {
			return errors.Errorf("项目目录[%q]包含其他文件，请选择空目录创建项目", cfg.ProjectDir)
		}
		return nil
	}
	return util.MKDir(cfg.ProjectDir)

}

// ShellModInit 初始化go.mod文件
func ShellModInit(cfg *config.ProjectConfig) error {
	exist, err := util.FileExist(path.Join(cfg.ProjectDir, "go.mod"))
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	cmd := exec.Command("go", "mod", "init", cfg.ModulePath)
	cmd.Dir = cfg.ProjectDir
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "mod", "edit", "-go=1.20")
	cmd.Dir = cfg.ProjectDir
	return cmd.Run()
}

// ShellPinDependencies pins versions that are compatible with the generated
// project's minimum Go version. This keeps generation reproducible and avoids
// go mod tidy unexpectedly raising the go directive.
func ShellPinDependencies(cfg *config.ProjectConfig) error {
	dependencies := []string{
		"github.com/pkg/errors@v0.9.1",
		"github.com/spf13/cast@v1.5.1",
		"github.com/spf13/cobra@v1.7.0",
		"github.com/spf13/viper@v1.16.0",
	}
	if _, ok := cfg.Modules["zap"]; ok {
		dependencies = append(dependencies, "go.uber.org/zap@v1.24.0")
	}
	if _, ok := cfg.Modules["gin"]; ok {
		dependencies = append(dependencies,
			"github.com/gin-contrib/cors@v1.4.0",
			"github.com/gin-gonic/gin@v1.9.1",
			"golang.org/x/sync@v0.2.0",
		)
	}
	if _, ok := cfg.Modules["gorm"]; ok {
		dependencies = append(dependencies,
			"gorm.io/driver/mysql@v1.5.1",
			"gorm.io/gorm@v1.25.2",
			"gorm.io/plugin/dbresolver@v1.5.0",
		)
	}
	cmd := exec.Command("go", append([]string{"get"}, dependencies...)...)
	cmd.Dir = cfg.ProjectDir
	return cmd.Run()
}

func ShellGoModTidy(cfg *config.ProjectConfig) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = cfg.ProjectDir
	return cmd.Run()
}
