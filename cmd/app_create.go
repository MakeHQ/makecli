/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、fmt、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newAppCreateCmd 函数
 * [POS]: cmd/app 的 create 子命令，调用 Meta Server API 创建 App，支持 --code 选项
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const defaultMetaServer = "https://dev-make.qtech.cn/api/make"

func newAppCreateCmd() *cobra.Command {
	var profile string
	var server string
	var code string

	cmd := &cobra.Command{
		Use:          "create <name>",
		Short:        "Create a new app on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppCreate(args[0], code, profile, server)
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "default", "credentials profile to use")
	cmd.Flags().StringVar(&server, "server", defaultMetaServer, "Meta Server base URL")
	cmd.Flags().StringVar(&code, "code", "", "app code (defaults to name)")
	return cmd
}

func runAppCreate(name, code, profile, server string) error {
	client, err := newClientFromProfile(profile, server)
	if err != nil {
		return err
	}

	var apiErr error
	if code == "" {
		apiErr = client.CreateApp(name) // code 默认为 name
	} else {
		apiErr = client.CreateAppWithCode(name, code)
	}
	if apiErr != nil {
		return apiErr
	}

	fmt.Printf("App '%s' created successfully\n", name)
	return nil
}
