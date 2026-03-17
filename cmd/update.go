/**
 * [INPUT]: 依赖 github.com/spf13/cobra、internal/update、internal/build
 * [OUTPUT]: 对外提供 newUpdateCmd 函数
 * [POS]: cmd 模块的 update 子命令，从 GitHub Releases 自更新二进制
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeHQ/makecli/internal/build"
	"github.com/MakeHQ/makecli/internal/update"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update makecli to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd)
		},
	}
}

func runUpdate(cmd *cobra.Command) error {
	currentVersion := build.Version

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Checking for updates...\n")

	release, newer, err := update.CheckLatest(currentVersion)
	if err != nil {
		return err
	}

	if !newer {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Already up to date (%s)\n", release.TagName)
		return nil
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updating makecli: %s → %s\n",
		formatCurrentVersion(currentVersion), release.TagName)

	if err := update.Apply(release); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated makecli: %s → %s\n",
		formatCurrentVersion(currentVersion), release.TagName)

	return nil
}

// formatCurrentVersion 格式化当前版本号用于显示
func formatCurrentVersion(v string) string {
	v = strings.TrimPrefix(v, "v")
	if v == "DEV" {
		return v
	}
	return "v" + v
}
