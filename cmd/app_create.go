/**
 * [INPUT]: 依赖 fmt、os、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newAppCreateCmd 函数
 * [POS]: cmd/app 的 create 子命令，在当前目录创建 <name> 文件夹
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newAppCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "create <name>",
		Short:        "Create a new app",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppCreate(args[0])
		},
	}
}

func runAppCreate(name string) error {
	if err := os.Mkdir(name, 0755); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("'%s' already exists", name)
		}
		return err
	}
	fmt.Printf("Created app '%s'\n", name)
	return nil
}
