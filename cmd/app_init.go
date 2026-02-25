/**
 * [INPUT]: 依赖 fmt、os、path/filepath、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newAppInitCmd 函数
 * [POS]: cmd/app 的 init 子命令，在已有 Folder 内创建 provider 对应的配置文件
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// providerFile 每个 provider 对应的配置文件名
var providerFile = map[string]string{
	"anthropic": "CLAUDE.md",
	"openai":    "AGENTS.md",
	"google":    "GEMINI.md",
	"cursor":    ".cursorrules",
}

func newAppInitCmd() *cobra.Command {
	var provider string

	cmd := &cobra.Command{
		Use:          "init <folder>",
		Short:        "Initialize an app with provider config file",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppInit(args[0], provider)
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "anthropic", "AI provider (anthropic|openai|google|cursor)")
	return cmd
}

func runAppInit(folder, provider string) error {
	filename, ok := providerFile[provider]
	if !ok {
		return fmt.Errorf("unknown provider '%s', valid options: anthropic, openai, google, cursor", provider)
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return fmt.Errorf("'%s' does not exist, run 'app create %s' first", folder, folder)
	}

	target := filepath.Join(folder, filename)
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("'%s' already exists", target)
	}

	if err := os.WriteFile(target, []byte(""), 0644); err != nil {
		return err
	}

	fmt.Printf("Initialized '%s' with %s\n", folder, filename)
	return nil
}
