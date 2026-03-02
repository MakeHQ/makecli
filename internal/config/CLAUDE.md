# internal/config/
> L2 | 父级: /CLAUDE.md

## 成员清单
credentials.go:      读写 ~/.make/credentials（INI 格式），提供 Load/Save/CredentialsPath，Credentials/Profile 类型
credentials_test.go: 覆盖 parseINI（白盒）+ Load/Save 全路径测试，用 t.Setenv("HOME",...) 隔离文件系统

[PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
