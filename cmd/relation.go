/**
 * [INPUT]: 依赖 github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newRelationCmd 函数
 * [POS]: cmd 模块的 relation 命令组，挂载 create / update / delete / list 子命令，--app 参数为子命令继承
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import "github.com/spf13/cobra"

func newRelationCmd() *cobra.Command {
	var app string

	cmd := &cobra.Command{
		Use:   "relation",
		Short: "Manage relations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if app == "" {
				return cmd.Usage()
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&app, "app", "", "app name (required)")
	_ = cmd.MarkPersistentFlagRequired("app")

	cmd.AddCommand(newRelationCreateCmd())
	cmd.AddCommand(newRelationUpdateCmd())
	cmd.AddCommand(newRelationDeleteCmd())
	cmd.AddCommand(newRelationListCmd())
	return cmd
}
