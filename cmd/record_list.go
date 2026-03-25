package cmd

import "github.com/spf13/cobra"

func newRecordListCmd() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List records in an entity"}
}
