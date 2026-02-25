# cmd/
> L2 | 父级: /CLAUDE.md

## 成员清单
root.go:        根命令入口，挂载所有子命令，对外暴露 Execute(version, date)
version.go:     version 子命令，格式化版本输出（参考 GitHub CLI 模式）
configure.go:   configure 子命令，交互式写入 ~/.make/credentials，支持 --profile
app.go:         app 命令组，挂载 app 相关子命令
app_create.go:  app create 子命令，在当前目录创建 <name> 文件夹，已存在则报错
app_init.go:    app init 子命令，在已有 Folder 内创建 provider 对应配置文件（anthropic→CLAUDE.md / openai→AGENTS.md / google→GEMINI.md / cursor→.cursorrules）

[PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
