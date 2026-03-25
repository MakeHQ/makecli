package cmd

import "github.com/spf13/cobra"

func newRecordCreateCmd() *cobra.Command {
	return &cobra.Command{Use: "create", Short: "Create a new record in an entity"}
}
