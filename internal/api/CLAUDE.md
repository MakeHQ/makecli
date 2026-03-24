# internal/api/
> L2 | 父级: /CLAUDE.md

## 成员清单
client.go:      Make Meta Service 的 HTTP 客户端，提供 Client 类型（含 debug/headers 字段）、Option 函数选项类型、WithDebug/WithHeaders 选项、New(baseURL, token, ...Option) 构造函数、App / Field / Entity / EntityProperties / RelationEnd / RelationProperties / Relation 类型、CreateApp / CreateAppWithCode / ListApps / DeleteApp / GetApp / CreateEntity / ListEntities / GetEntity / UpdateEntity / DeleteEntity / CreateRelation / UpdateRelation / ListRelations / GetRelation / DeleteRelation 方法；do() 为底层 POST，支持 debug 输出 curl 命令 + 自定义 headers 注入，post() 处理写操作
client_test.go: 覆盖 CreateApp / DeleteApp / ListApps / WithHeaders / WithDebug 的单元测试（成功/API错误/空列表/格式错误/自定义头/调试模式），用 httptest 隔离网络

[PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
