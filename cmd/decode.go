package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var File string
var Dir string

func init() {
	rootCmd.AddCommand(decodeCmd)
	decodeCmd.Flags().StringVarP(&File, "FILE", "f", "", "decode file path")
	decodeCmd.Flags().StringVarP(&Dir, "DIR", "d", "", "decode dir path")
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "decode music file",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if File == "" && Dir == "" {
			return errors.New("requires a file path")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		decode(args)
	},
}

func decode(args []string) error {
	if len(args) <= 0 {
		fmt.Print("A command must be supplied to run1111")
		return fmt.Errorf("A command must be supplied to run")
	}
	return nil
}
