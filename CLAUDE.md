# makecli - MakeHQ 的命令行工具
Go + github.com/spf13/cobra

<directory>
cmd/            - Cobra 子命令层（root、version）
internal/build/ - 构建元数据（Version/Date，由 ldflags 注入）
</directory>

<config>
go.mod   - 模块声明，module github.com/MakeHQ/makecli
Makefile - 构建脚本，通过 ldflags 注入版本和日期
</config>
