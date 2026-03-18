/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、fmt、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newEntityDeleteCmd 函数
 * [POS]: cmd/entity 的 delete 子命令，调用 Meta Server API 删除指定 Entity
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newEntityDeleteCmd() *cobra.Command {
	var profile string
	var server string

	cmd := &cobra.Command{
		Use:          "delete <name>",
		Short:        "Delete an entity on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, _ := cmd.Parent().Flags().GetString("app")
			return runEntityDelete(args[0], app, profile, server)
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "default", "credentials profile to use")
	cmd.Flags().StringVar(&server, "server", defaultMetaServer, "Meta Server base URL")
	return cmd
}

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
