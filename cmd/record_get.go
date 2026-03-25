package cmd

import "github.com/spf13/cobra"

func newRecordGetCmd() *cobra.Command {
	return &cobra.Command{Use: "get", Short: "Get a record by ID"}
}
