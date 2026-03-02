/**
 * [INPUT]: 依赖 cmd 包内的 runAppInit（包内白盒），os、path/filepath
 * [OUTPUT]: 覆盖 app init 子命令核心逻辑的单元测试
 * [POS]: cmd 模块 app_init.go 的配套测试
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunAppInit(t *testing.T) {
	t.Run("rejects unknown provider", func(t *testing.T) {
		if err := runAppInit(t.TempDir(), "unknown"); err == nil {
			t.Fatal("expected error for unknown provider")
		}
	})

	t.Run("fails if folder does not exist", func(t *testing.T) {
		if err := runAppInit("/nonexistent/path/xyz", "anthropic"); err == nil {
			t.Fatal("expected error for nonexistent folder")
		}
	})

	t.Run("fails if config file already exists", func(t *testing.T) {
		dir := t.TempDir()
		_ = runAppInit(dir, "anthropic")
		if err := runAppInit(dir, "anthropic"); err == nil {
			t.Fatal("expected error for existing config file")
		}
	})

	// ---------------------------------- 每个 provider 验证正确文件名 ----------------------------------
	providerCases := []struct {
		provider     string
		expectedFile string
	}{
		{"anthropic", "CLAUDE.md"},
		{"openai", "AGENTS.md"},
		{"google", "GEMINI.md"},
		{"cursor", ".cursorrules"},
	}

	for _, tc := range providerCases {
		tc := tc // capture
		t.Run("provider_"+tc.provider, func(t *testing.T) {
			dir := t.TempDir()
			if err := runAppInit(dir, tc.provider); err != nil {
				t.Fatalf("runAppInit(%q): %v", tc.provider, err)
			}
			target := filepath.Join(dir, tc.expectedFile)
			if _, err := os.Stat(target); err != nil {
				t.Errorf("expected file %q not found: %v", target, err)
			}
		})
	}
}
