/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、fmt、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newAppCreateCmd 函数
 * [POS]: cmd/app 的 create 子命令，调用 Meta Server API 创建 App，支持 --description 选项
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAppCreateCmd() *cobra.Command {
	var profile string
	var description string

	cmd := &cobra.Command{
		Use:          "create <name>",
		Short:        "Create a new app on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppCreate(args[0], description, profile)
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "default", "credentials profile to use")
	cmd.Flags().StringVar(&description, "description", "", "app description")
	return cmd
}

func runAppCreate(name, description, profile string) error {
	client, err := newClientFromProfile(profile)
	if err != nil {
		return err
	}

	props := map[string]any{}
	if description != "" {
		props["description"] = description
	}

	if apiErr := client.CreateApp(name, props); apiErr != nil {
		return apiErr
	}

	fmt.Printf("App '%s' created successfully\n", name)
	return nil
}
