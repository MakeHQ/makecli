/**
 * [INPUT]: 依赖 internal/config，bufio、encoding/base64、fmt、os、strings、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newConfigureCmd 函数
 * [POS]: cmd 模块的 configure 子命令，交互式写入 ~/.make/credentials
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/MakeHQ/makecli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigureCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:          "configure",
		Short:        "Configure MakeHQ credentials",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigure(profile)
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "default", "credentials profile name")
	return cmd
}

func runConfigure(profile string) error {
	creds, err := config.Load()
	if err != nil {
		return err
	}

	current := creds[profile]

	fmt.Printf("Configuring profile [%s]\n", profile)

	token, err := prompt("MakeHQ Access Token", current.AccessToken)
	if err != nil {
		return err
	}
	if token != "" {
		if err := validateJWT(token); err != nil {
			return err
		}
		current.AccessToken = token
	}

	creds[profile] = current
	if err := config.Save(creds); err != nil {
		return err
	}

	path, _ := config.CredentialsPath()
	fmt.Printf("\nCredentials saved to %s\n", path)
	return nil
}

// prompt 打印提示行（已有值则遮掩末尾4位显示），读取用户输入
// 用户直接回车表示保留当前值，返回空字符串
func prompt(label, current string) (string, error) {
	hint := "None"
	if current != "" {
		hint = mask(current)
	}
	fmt.Printf("%s [%s]: ", label, hint)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// mask 保留末尾4位，其余替换为 *
// 短于4位则全部遮掩
func mask(s string) string {
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}

// validateJWT 校验 token 是否符合 JWT 格式（三段 base64url，不验证签名）
func validateJWT(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid token format: expected JWT (3 base64url segments separated by '.')")
	}
	for i, part := range parts {
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			return fmt.Errorf("invalid token format: segment %d is not valid base64url", i+1)
		}
	}
	return nil
}
