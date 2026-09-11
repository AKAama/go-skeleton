# go-skeleton

Go 项目骨架生成器。

## 安装

```bash
go install github.com/AKAama/go-skeleton@latest
```

Go 会把可执行文件安装到 `GOBIN`；未设置 `GOBIN` 时，默认安装到
`$(go env GOPATH)/bin`。如果安装后执行 `go-skeleton` 提示
`command not found`，说明该目录尚未加入 `PATH`。

先查看当前配置：

```bash
go env GOBIN GOPATH
```

### macOS / Linux（Zsh）

```zsh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### macOS / Linux（Bash）

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Windows（PowerShell）

```powershell
$goBin = "$(go env GOPATH)\bin"
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
[Environment]::SetEnvironmentVariable("Path", "$userPath;$goBin", "User")
```

修改后重新打开终端，并验证：

```bash
go-skeleton --version
```

如果设置了自定义 `GOBIN`，请将 `go env GOBIN` 输出的目录加入 `PATH`，
不要再使用默认的 `GOPATH/bin`。

## 使用

```bash
go-skeleton -p <项目名称>
```

## 参数

```text
-d, --dir string        项目所在文件夹路径（默认当前路径）
-m, --mod string        go.mod 的 module 路径（默认使用项目名称）
-p, --project string    项目名称
-v, --version           显示版本
-w, --without strings   不启用指定框架，可选 gin、zap、gorm
```
