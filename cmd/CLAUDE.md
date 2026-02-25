# cmd/
> L2 | 父级: /CLAUDE.md

## 成员清单
root.go:       根命令入口，挂载所有子命令，对外暴露 Execute(version, date)
version.go:    version 子命令，格式化版本输出（参考 GitHub CLI 模式）
configure.go:  configure 子命令，交互式写入 ~/.make/credentials，支持 --profile

[PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
