/**
 * [INPUT]: 依赖 github.com/spf13/cobra、fmt、regexp、strings
 * [OUTPUT]: 对外提供 newVersionCmd 函数
 * [POS]: cmd 模块的 version 子命令，参考 GitHub CLI 的 pkg/cmd/version/version.go 实现
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func newVersionCmd(version, buildDate string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(formatVersion(version, buildDate))
		},
	}
}

// formatVersion 生成版本输出字符串，格式参考 GitHub CLI
func formatVersion(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	dateStr := ""
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}

	return fmt.Sprintf("makecli version %s%s\n%s\n", version, dateStr, changelogURL(version))
}

// changelogURL 根据版本号是否符合 semver 返回对应的 release 页面
func changelogURL(version string) string {
	semver := regexp.MustCompile(`^\d+\.\d+\.\d+(-[\w.]+)?$`)
	if semver.MatchString(version) {
		return fmt.Sprintf("https://github.com/MakeHQ/makecli/releases/tag/v%s", version)
	}
	return "https://github.com/MakeHQ/makecli/releases/latest"
}
