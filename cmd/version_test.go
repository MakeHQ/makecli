/**
 * [INPUT]: 依赖 cmd 包内的 formatVersion、changelogURL（包内白盒）
 * [OUTPUT]: 覆盖版本格式化逻辑的单元测试
 * [POS]: cmd 模块 version.go 的配套测试
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"strings"
	"testing"
)

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		version   string
		buildDate string
		wantIn    []string
	}{
		{
			version:   "1.2.3",
			buildDate: "2024-01-15",
			wantIn:    []string{"makecli version 1.2.3 (2024-01-15)", "releases/tag/v1.2.3"},
		},
		{
			// v 前缀应被剥离
			version:   "v1.2.3",
			buildDate: "2024-01-15",
			wantIn:    []string{"makecli version 1.2.3", "releases/tag/v1.2.3"},
		},
		{
			// 无 buildDate 时不显示括号
			version:   "1.0.0",
			buildDate: "",
			wantIn:    []string{"makecli version 1.0.0\n"},
		},
		{
			// 非 semver 指向 latest
			version:   "DEV",
			buildDate: "",
			wantIn:    []string{"makecli version DEV", "releases/latest"},
		},
		{
			version:   "abc123",
			buildDate: "2024-06-01",
			wantIn:    []string{"makecli version abc123 (2024-06-01)", "releases/latest"},
		},
	}

	for _, tt := range tests {
		got := formatVersion(tt.version, tt.buildDate)
		for _, want := range tt.wantIn {
			if !strings.Contains(got, want) {
				t.Errorf("formatVersion(%q, %q) = %q, want to contain %q",
					tt.version, tt.buildDate, got, want)
			}
		}
	}
}

func TestChangelogURL(t *testing.T) {
	tests := []struct {
		version string
		wantIn  string
	}{
		{"1.2.3", "releases/tag/v1.2.3"},
		{"1.0.0-beta.1", "releases/tag/v1.0.0-beta.1"},
		{"0.0.1", "releases/tag/v0.0.1"},
		{"DEV", "releases/latest"},
		{"abc123", "releases/latest"},
		{"1.2", "releases/latest"},         // 不完整 semver
		{"1.2.3.4", "releases/latest"},     // 多余段
	}

	for _, tt := range tests {
		got := changelogURL(tt.version)
		if !strings.Contains(got, tt.wantIn) {
			t.Errorf("changelogURL(%q) = %q, want to contain %q", tt.version, got, tt.wantIn)
		}
	}
}
