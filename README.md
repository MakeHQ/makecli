## 说明
makecli 是 make 平台管理的命令行工具

## 安装
```bash
brew tap MakeHQ/makecli
brew install makecli
```
## 功能

### 配置凭证
```bash
# 配置默认 profile
makecli configure

# 配置指定 profile
makecli configure --profile todo
```

交互示例：
```
Configuring profile [default]
MakeHQ Access Token [****YDUW]:
Credentials saved to ~/.make/credentials
```

凭证保存在 `~/.make/credentials`，格式：
```ini
[default]
access_token = AKIAUXFQEUPWGEXEYDUW

[todo]
access_token = AKIAUXFQEUPWGEXEYDUW
```

### 查看版本
```bash
makecli version
```
