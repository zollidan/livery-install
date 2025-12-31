/*
Copyright © 2025 ZOLLIDAN zollidan@aol.com
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [zip-file-full-path]",
	Short: "A brief description of your command",
	Long: `A longer description `,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		info, err := os.Stat(args[0])
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("Error: File does not exist, %s", err.Error())
				return
			}
			fmt.Printf("Error accessing the file: %s", err.Error())
			return
		}

		fmt.Println(info.Name())
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
