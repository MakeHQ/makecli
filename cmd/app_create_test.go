/**
 * [INPUT]: 依赖 cmd 包内的 runAppCreate（包内白盒），os、path/filepath
 * [OUTPUT]: 覆盖 app create 子命令核心逻辑的单元测试
 * [POS]: cmd 模块 app_create.go 的配套测试
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunAppCreate(t *testing.T) {
	t.Run("creates directory", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "myapp")
		if err := runAppCreate(name); err != nil {
			t.Fatalf("runAppCreate: %v", err)
		}
		if _, err := os.Stat(name); err != nil {
			t.Errorf("expected directory %q to exist: %v", name, err)
		}
	})

	t.Run("fails if directory already exists", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "myapp")
		_ = runAppCreate(name) // 先建一次
		if err := runAppCreate(name); err == nil {
			t.Fatal("expected error for existing directory, got nil")
		}
	})
}
