# Config File Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `~/.make/config` support with `x-tenant-id` and `operator-id` per profile, injected as HTTP headers on every API request.

**Architecture:** New `internal/config/config.go` mirrors credentials pattern. `api.Client` gains functional options for custom headers. `cmd/client.go` extracts shared credential+config loading. `configure` becomes a command group with `token`/`config`/`set`/`get` subcommands.

**Tech Stack:** Go, cobra, INI format, httptest

**Spec:** `docs/superpowers/specs/2026-03-18-config-file-support-design.md`

---

### Task 1: `internal/config/config.go` — Config 读写

**Files:**
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`

- [ ] **Step 1: Write failing tests for config.go**

In `internal/config/config_test.go`, add tests after existing content:

```go
// ---------------------------------- Config ----------------------------------

func TestConfigPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	want := filepath.Join(home, ".make", "config")
	if path != want {
		t.Errorf("ConfigPath = %q, want %q", path, want)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(cfg) != 0 {
		t.Errorf("expected empty config, got %v", cfg)
	}
}

func TestSaveConfigAndLoad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	original := Config{
		"default": {XTenantID: "tenant-1", OperatorID: "op-1"},
		"staging": {XTenantID: "tenant-2", OperatorID: ""},
	}

	if err := SaveConfig(original); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	path, _ := ConfigPath()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("config file not found: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %v, want 0600", info.Mode().Perm())
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	for profile, want := range original {
		got := loaded[profile]
		if got.XTenantID != want.XTenantID {
			t.Errorf("profile %q: XTenantID = %q, want %q", profile, got.XTenantID, want.XTenantID)
		}
		if got.OperatorID != want.OperatorID {
			t.Errorf("profile %q: OperatorID = %q, want %q", profile, got.OperatorID, want.OperatorID)
		}
	}
}

func TestParseConfigINI(t *testing.T) {
	t.Run("both keys", func(t *testing.T) {
		f := writeTempINI(t, "[default]\nx-tenant-id = t1\noperator-id = o1\n")
		defer func() { _ = f.Close() }()

		cfg, err := parseConfigINI(f)
		if err != nil {
			t.Fatalf("parseConfigINI: %v", err)
		}
		if cfg["default"].XTenantID != "t1" || cfg["default"].OperatorID != "o1" {
			t.Errorf("unexpected config: %+v", cfg["default"])
		}
	})

	t.Run("partial keys", func(t *testing.T) {
		f := writeTempINI(t, "[default]\nx-tenant-id = t1\n")
		defer func() { _ = f.Close() }()

		cfg, err := parseConfigINI(f)
		if err != nil {
			t.Fatalf("parseConfigINI: %v", err)
		}
		if cfg["default"].XTenantID != "t1" {
			t.Errorf("XTenantID = %q, want %q", cfg["default"].XTenantID, "t1")
		}
		if cfg["default"].OperatorID != "" {
			t.Errorf("OperatorID = %q, want empty", cfg["default"].OperatorID)
		}
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Volumes/Coding/make/repos/makecli && go test ./internal/config/ -run "TestConfig|TestParseConfig|TestSaveConfig" -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement config.go**

Create `internal/config/config.go`:

```go
/**
 * [INPUT]: 依赖 os、bufio、fmt、strings、path/filepath
 * [OUTPUT]: 对外提供 LoadConfig、SaveConfig、ConfigPath 函数，Config/ConfigProfile 类型
 * [POS]: internal/config 的 config 文件管理，读写 ~/.make/config（INI 格式）
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------- 数据结构 ----------------------------------

type ConfigProfile struct {
	XTenantID  string
	OperatorID string
}

type Config map[string]ConfigProfile

// ---------------------------------- 路径 ----------------------------------

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取 home 目录: %w", err)
	}
	return filepath.Join(home, ".make", "config"), nil
}

// ---------------------------------- 读取 ----------------------------------

func LoadConfig() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return Config{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取 config 失败: %w", err)
	}
	defer func() { _ = f.Close() }()

	return parseConfigINI(f)
}

func parseConfigINI(f *os.File) (Config, error) {
	cfg := Config{}
	current := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = strings.TrimSpace(line[1 : len(line)-1])
			if _, ok := cfg[current]; !ok {
				cfg[current] = ConfigProfile{}
			}
			continue
		}

		if current == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		p := cfg[current]
		switch key {
		case "x-tenant-id":
			p.XTenantID = val
		case "operator-id":
			p.OperatorID = val
		}
		cfg[current] = p
	}

	return cfg, scanner.Err()
}

// ---------------------------------- 写入 ----------------------------------

func SaveConfig(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("创建 ~/.make 目录失败: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("写入 config 失败: %w", err)
	}
	defer func() { _ = f.Close() }()

	order := []string{}
	if _, ok := cfg["default"]; ok {
		order = append(order, "default")
	}
	for name := range cfg {
		if name != "default" {
			order = append(order, name)
		}
	}

	w := bufio.NewWriter(f)
	for i, name := range order {
		if i > 0 {
			_, _ = fmt.Fprintln(w)
		}
		_, _ = fmt.Fprintf(w, "[%s]\n", name)
		p := cfg[name]
		if p.XTenantID != "" {
			_, _ = fmt.Fprintf(w, "x-tenant-id = %s\n", p.XTenantID)
		}
		if p.OperatorID != "" {
			_, _ = fmt.Fprintf(w, "operator-id = %s\n", p.OperatorID)
		}
	}

	return w.Flush()
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Volumes/Coding/make/repos/makecli && go test ./internal/config/ -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add ~/.make/config read/write support with ConfigProfile"
```

---

### Task 2: API Functional Options + Client Helper + Migrate All Commands (Atomic)

> **IMPORTANT:** This task modifies `api.New` signature, creates the helper, and migrates all 8 callers in one atomic commit. Splitting would break compilation.

**Files:**
- Modify: `internal/api/client.go`
- Modify: `internal/api/client_test.go`
- Create: `cmd/client.go`
- Modify: `cmd/app_create.go`, `cmd/app_list.go`, `cmd/app_delete.go`
- Modify: `cmd/entity_create.go`, `cmd/entity_list.go`, `cmd/entity_delete.go`
- Modify: `cmd/apply.go`, `cmd/diff.go`

- [ ] **Step 1: Write failing tests for WithHeaders**

Append to `internal/api/client_test.go`:

```go
func TestWithHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-tenant-id"); got != "tenant-abc" {
			t.Errorf("x-tenant-id = %q, want %q", got, "tenant-abc")
		}
		if got := r.Header.Get("operator-id"); got != "op-123" {
			t.Errorf("operator-id = %q, want %q", got, "op-123")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"code": 200, "msg": "ok"})
	}))
	defer srv.Close()

	headers := map[string]string{
		"x-tenant-id": "tenant-abc",
		"operator-id": "op-123",
	}
	client := New(srv.URL, "test-token", WithHeaders(headers))
	if err := client.CreateApp("test"); err != nil {
		t.Fatalf("CreateApp with headers: %v", err)
	}
}

func TestWithDebugOption(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"code": 200, "msg": "ok"})
	}))
	defer srv.Close()

	client := New(srv.URL, "test-token", WithDebug(true))
	if err := client.CreateApp("test"); err != nil {
		t.Fatalf("CreateApp with debug: %v", err)
	}
}
```

- [ ] **Step 2: Modify client.go — functional options + headers**

In `internal/api/client.go`:

Add `headers` field to Client struct:
```go
type Client struct {
	baseURL    string
	token      string
	headers    map[string]string
	httpClient *http.Client
	debug      bool
}
```

Replace `New` function with functional options:
```go
type Option func(*Client)

func WithDebug(on bool) Option {
	return func(c *Client) { c.debug = on }
}

func WithHeaders(h map[string]string) Option {
	return func(c *Client) { c.headers = h }
}

func New(baseURL, token string, opts ...Option) *Client {
	c := &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
```

In `do()` method, after `req.Header.Set("X-Make-Target", target)`, add:
```go
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
```

In `do()` debug output, after the `X-Make-Target` curl line, add:
```go
		for k, v := range c.headers {
			fmt.Fprintf(os.Stderr, "  -H '%s: %s' \\\n", k, v)
		}
```

- [ ] **Step 3: Create cmd/client.go**

```go
/**
 * [INPUT]: 依赖 internal/config（Load/LoadConfig）、internal/api（New/WithDebug/WithHeaders）、fmt
 * [OUTPUT]: 对外提供 newClientFromProfile 函数
 * [POS]: cmd 模块的公共 helper，统一「凭证 + 配置 → API 客户端」的构建逻辑
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/qfeius/makecli/internal/api"
	"github.com/qfeius/makecli/internal/config"
)

func newClientFromProfile(profile, server string) (*api.Client, error) {
	creds, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("加载凭证失败: %w", err)
	}

	p, ok := creds[profile]
	if !ok || p.AccessToken == "" {
		return nil, fmt.Errorf("profile '%s' 未配置，请先运行: makecli configure --profile %s", profile, profile)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	headers := map[string]string{}
	if cp, ok := cfg[profile]; ok {
		if cp.XTenantID != "" {
			headers["x-tenant-id"] = cp.XTenantID
		}
		if cp.OperatorID != "" {
			headers["operator-id"] = cp.OperatorID
		}
	}

	return api.New(server, p.AccessToken, api.WithDebug(DebugMode), api.WithHeaders(headers)), nil
}
```

- [ ] **Step 4: Migrate all 8 command files**

Each file: replace `config.Load()` + profile check + `api.New(server, p.AccessToken, DebugMode)` with `newClientFromProfile(profile, server)`. Remove unused imports.

**app_create.go:**
```go
func runAppCreate(name, code, profile, server string) error {
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
	if code == "" {
		err = client.CreateApp(name)
	} else {
		err = client.CreateAppWithCode(name, code)
	}
	if err != nil {
		return err
	}
	fmt.Printf("App '%s' created successfully\n", name)
	return nil
}
```
Remove: `"github.com/qfeius/makecli/internal/api"`, `"github.com/qfeius/makecli/internal/config"`.

**app_delete.go:**
```go
func runAppDelete(name, profile, server string) error {
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
	if err := client.DeleteApp(name); err != nil {
		return err
	}
	fmt.Printf("App '%s' deleted successfully\n", name)
	return nil
}
```
Remove: `"github.com/qfeius/makecli/internal/api"`, `"github.com/qfeius/makecli/internal/config"`.

**app_list.go** — replace credential block in `runAppList`:
```go
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
	apps, total, err := client.ListApps(page, size)
```
Remove: `"github.com/qfeius/makecli/internal/api"`, `"github.com/qfeius/makecli/internal/config"`.

**entity_create.go** — replace credential block:
```go
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
```
Then `client.CreateEntity(...)`. Remove `"github.com/qfeius/makecli/internal/config"`. Keep `"github.com/qfeius/makecli/internal/api"` for `api.Field`.

**entity_list.go** — replace credential block:
```go
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
	if entityName != "" {
		return showEntity(client, app, entityName, output)
	}
	return listEntities(client, app, page, size, output)
```
Remove `"github.com/qfeius/makecli/internal/config"`. Keep `"github.com/qfeius/makecli/internal/api"` for `api.Client`.

**entity_delete.go:**
```go
func runEntityDelete(name, app, profile, server string) error {
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
	if err := client.DeleteEntity(name, app); err != nil {
		return err
	}
	fmt.Printf("Entity '%s' deleted successfully from app '%s'\n", name, app)
	return nil
}
```
Remove: `"github.com/qfeius/makecli/internal/api"`, `"github.com/qfeius/makecli/internal/config"`.

**apply.go** — replace credential block + update `applyResources` signature:
```go
func runAppApply(path, profile, server string) error {
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
	// ... file loading unchanged ...
	if err := applyResources(resources, client); err != nil {
```
Change `applyResources` from `(resources, server, token)` to `(resources, client *api.Client)`. Remove inner `api.New()`. Remove `"github.com/qfeius/makecli/internal/config"`.

**diff.go** — replace credential block:
```go
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}
```
Remove later `client := api.New(...)`. Remove `"github.com/qfeius/makecli/internal/config"`.

- [ ] **Step 5: Run all tests**

Run: `cd /Volumes/Coding/make/repos/makecli && go vet ./... && go test ./... -v`
Expected: ALL PASS

- [ ] **Step 6: Commit**

```bash
git add internal/api/client.go internal/api/client_test.go cmd/client.go \
  cmd/app_create.go cmd/app_list.go cmd/app_delete.go \
  cmd/entity_create.go cmd/entity_list.go cmd/entity_delete.go \
  cmd/apply.go cmd/diff.go
git commit -m "feat: functional options for api.Client, extract newClientFromProfile, migrate all commands"
```

---

### Task 3: Restructure `configure` Command Group + `token` Subcommand

**Files:**
- Modify: `cmd/configure.go`

- [ ] **Step 1: Refactor configure.go — command group with token subcommand**

Replace `newConfigureCmd` and rename `runConfigure` → `runConfigureToken`:

```go
func newConfigureCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:          "configure",
		Short:        "Configure MakeCLI credentials and settings",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureToken(profile)
		},
	}

	cmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	cmd.AddCommand(newConfigureTokenCmd(&profile))

	return cmd
}

func newConfigureTokenCmd(profile *string) *cobra.Command {
	return &cobra.Command{
		Use:          "token",
		Short:        "Configure access token (writes to ~/.make/credentials)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureToken(*profile)
		},
	}
}

func runConfigureToken(profile string) error {
	creds, err := config.Load()
	if err != nil {
		return err
	}

	current := creds[profile]
	fmt.Printf("Configuring profile [%s]\n", profile)

	token, err := prompt("MakeCLI Access Token", current.AccessToken)
	if err != nil {
		return err
	}
	if token != "" {
		if err := validateJWT(token); err != nil {
			return err
		}
		current.AccessToken = token
	}

	creds[profile] = current
	if err := config.Save(creds); err != nil {
		return err
	}

	path, _ := config.CredentialsPath()
	fmt.Printf("\nCredentials saved to %s\n", path)
	return nil
}
```

- [ ] **Step 2: Run existing tests + verify build**

Run: `cd /Volumes/Coding/make/repos/makecli && go build ./... && go test ./cmd/ -run "TestMask|TestValidateJWT" -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add cmd/configure.go
git commit -m "refactor: restructure configure as command group with token subcommand"
```

---

### Task 4: `configure config` Subcommand (Interactive)

**Files:**
- Modify: `cmd/configure.go`

- [ ] **Step 1: Add configure config subcommand**

In `newConfigureCmd`, add `cmd.AddCommand(newConfigureConfigCmd(&profile))`.

```go
func newConfigureConfigCmd(profile *string) *cobra.Command {
	return &cobra.Command{
		Use:          "config",
		Short:        "Configure custom headers (writes to ~/.make/config)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureConfig(*profile)
		},
	}
}

func runConfigureConfig(profile string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	current := cfg[profile]
	fmt.Printf("Configuring config profile [%s]\n", profile)

	tenantID, err := prompt("x-tenant-id", current.XTenantID)
	if err != nil {
		return err
	}
	if tenantID != "" {
		current.XTenantID = tenantID
	}

	operatorID, err := prompt("operator-id", current.OperatorID)
	if err != nil {
		return err
	}
	if operatorID != "" {
		current.OperatorID = operatorID
	}

	cfg[profile] = current
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	path, _ := config.ConfigPath()
	fmt.Printf("\nConfig saved to %s\n", path)
	return nil
}
```

- [ ] **Step 2: Verify build**

Run: `cd /Volumes/Coding/make/repos/makecli && go build ./...`
Expected: SUCCESS

- [ ] **Step 3: Commit**

```bash
git add cmd/configure.go
git commit -m "feat: add configure config subcommand for x-tenant-id and operator-id"
```

---

### Task 5: `configure set` / `configure get` Subcommands

**Files:**
- Modify: `cmd/configure.go`
- Modify: `cmd/configure_test.go`

- [ ] **Step 1: Write failing tests for key validation**

Append to `cmd/configure_test.go`:

```go
func TestValidConfigKeys(t *testing.T) {
	if err := validateConfigKey("x-tenant-id"); err != nil {
		t.Errorf("x-tenant-id should be valid: %v", err)
	}
	if err := validateConfigKey("operator-id"); err != nil {
		t.Errorf("operator-id should be valid: %v", err)
	}
	if err := validateConfigKey("bad-key"); err == nil {
		t.Error("bad-key should be invalid")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Volumes/Coding/make/repos/makecli && go test ./cmd/ -run TestValidConfigKeys -v`
Expected: FAIL

- [ ] **Step 3: Implement set/get subcommands**

In `newConfigureCmd`, add:
```go
	cmd.AddCommand(newConfigureSetCmd(&profile))
	cmd.AddCommand(newConfigureGetCmd(&profile))
```

Add to `cmd/configure.go`:

```go
var validConfigKeys = []string{"x-tenant-id", "operator-id"}

func validateConfigKey(key string) error {
	for _, k := range validConfigKeys {
		if key == k {
			return nil
		}
	}
	return fmt.Errorf("unknown config key '%s', valid keys: %s", key, strings.Join(validConfigKeys, ", "))
}

func newConfigureSetCmd(profile *string) *cobra.Command {
	return &cobra.Command{
		Use:          "set <key> <value>",
		Short:        "Set a config value (writes to ~/.make/config)",
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureSet(*profile, args[0], args[1])
		},
	}
}

func runConfigureSet(profile, key, value string) error {
	if err := validateConfigKey(key); err != nil {
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	p := cfg[profile]
	switch key {
	case "x-tenant-id":
		p.XTenantID = value
	case "operator-id":
		p.OperatorID = value
	}
	cfg[profile] = p

	return config.SaveConfig(cfg)
}

func newConfigureGetCmd(profile *string) *cobra.Command {
	return &cobra.Command{
		Use:          "get <key>",
		Short:        "Get a config value from ~/.make/config",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureGet(*profile, args[0])
		},
	}
}

func runConfigureGet(profile, key string) error {
	if err := validateConfigKey(key); err != nil {
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	p := cfg[profile]
	switch key {
	case "x-tenant-id":
		fmt.Println(p.XTenantID)
	case "operator-id":
		fmt.Println(p.OperatorID)
	}
	return nil
}
```

Add `"strings"` to imports.

- [ ] **Step 4: Run all tests**

Run: `cd /Volumes/Coding/make/repos/makecli && go build ./... && go test ./... -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/configure.go cmd/configure_test.go
git commit -m "feat: add configure set/get subcommands for config key-value access"
```

---

### Task 6: Update L2/L3 Documentation

**Files:**
- Modify: `internal/config/CLAUDE.md`
- Modify: `internal/api/CLAUDE.md`
- Modify: `cmd/CLAUDE.md`
- Modify: `CLAUDE.md` (L1)

- [ ] **Step 1: Update internal/config/CLAUDE.md**

Add `config.go` and `config_test.go` to member list.

- [ ] **Step 2: Update internal/api/CLAUDE.md**

Update `client.go` description: functional options (`Option`, `WithDebug`, `WithHeaders`), `headers` field.

- [ ] **Step 3: Update cmd/CLAUDE.md**

Add `client.go` to member list. Update `configure.go` description: command group with `token`/`config`/`set`/`get`.

- [ ] **Step 4: Update L1 CLAUDE.md**

Update `internal/config/` description to mention config file support.

- [ ] **Step 5: Commit**

```bash
git add internal/config/CLAUDE.md internal/api/CLAUDE.md cmd/CLAUDE.md CLAUDE.md
git commit -m "docs: update L1/L2 documentation for config file support"
```

---

### Task 7: Final Verification

- [ ] **Step 1: Run full test suite**

Run: `cd /Volumes/Coding/make/repos/makecli && go vet ./... && go test ./... -v`
Expected: ALL PASS, no vet issues

- [ ] **Step 2: Verify build**

Run: `cd /Volumes/Coding/make/repos/makecli && go build -o /dev/null .`
Expected: SUCCESS

- [ ] **Step 3: Smoke test commands exist**

```bash
cd /Volumes/Coding/make/repos/makecli
go run . configure --help
go run . configure token --help
go run . configure config --help
go run . configure set --help
go run . configure get --help
```
Expected: Each prints help text with correct usage.
