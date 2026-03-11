/**
 * [INPUT]: 依赖 encoding/json、fmt、os
 * [OUTPUT]: 对外提供 list 命令通用的输出格式校验和 JSON 编码辅助函数
 * [POS]: cmd 模块的输出层辅助，当前用于 app list / entity list 的 table|json 双输出
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	outputTable = "table"
	outputJSON  = "json"
)

func validateOutputFormat(output string) error {
	switch output {
	case outputTable, outputJSON:
		return nil
	default:
		return fmt.Errorf("unsupported output format %q, valid options: %s, %s", output, outputTable, outputJSON)
	}
}

func writeJSON(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
