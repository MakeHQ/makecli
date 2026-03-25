package cmd

import "github.com/spf13/cobra"

func newRecordDeleteCmd() *cobra.Command {
	return &cobra.Command{Use: "delete", Short: "Delete one or more records"}
}
