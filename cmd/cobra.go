package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"prince-x/cmd/api"
	"prince-x/cmd/migrate"
)

var rootCmd = &cobra.Command{
	Use:               "prince-x",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `prince-x`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usageStr := `prince-x 1.0.0 欢迎使用，可以是用 -h 查看命令`
		log.Printf("%s\n", usageStr)
	},
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
