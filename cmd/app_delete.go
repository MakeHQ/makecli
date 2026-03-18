/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、fmt、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newAppDeleteCmd 函数
 * [POS]: cmd/app 的 delete 子命令，调用 Meta Server API 删除指定 App
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAppDeleteCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:          "delete <name>",
		Short:        "Delete an app on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppDelete(args[0], profile)
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "default", "credentials profile to use")
	return cmd
}

func runAppDelete(name, profile string) error {
	client, err := newClientFromProfile(profile)
	if err != nil {
		return err
	}

	if err := client.DeleteApp(name); err != nil {
		return err
	}

	fmt.Printf("App '%s' deleted successfully\n", name)
	return nil
}
