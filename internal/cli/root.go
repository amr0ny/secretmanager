package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secretmanager",
		Short: "A CLI tool to generate secrets and manage them",
	}
	cmd.AddCommand(NewGenerateCommand())
	return cmd
}
