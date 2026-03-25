package cmd

import "github.com/spf13/cobra"

func newRecordUpdateCmd() *cobra.Command {
	return &cobra.Command{Use: "update", Short: "Update one or more records"}
}
