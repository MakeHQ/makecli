/**
 * [INPUT]: 依赖 runtime/debug 的 ReadBuildInfo
 * [OUTPUT]: 对外提供 Version、Date 变量（可通过 ldflags 在构建时注入）
 * [POS]: internal/build 的版本元数据持有者，被 cmd/version.go 消费
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package build

import "runtime/debug"

// Version 在构建时通过 ldflags 注入，默认为 DEV
var Version = "DEV"

// Date 在构建时通过 ldflags 注入，格式 YYYY-MM-DD
var Date = ""

func init() {
	// 若未通过 ldflags 注入，尝试从 go module 构建信息中读取版本
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
