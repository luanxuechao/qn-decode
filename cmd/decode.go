package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/luanxuechao/qn-decode/util"
	"github.com/spf13/cobra"
)

var filename string
var dirname string

func init() {
	rootCmd.AddCommand(decodeCmd)
	decodeCmd.Flags().StringVarP(&filename, "FILE", "f", "", "decode file path")
	decodeCmd.Flags().StringVarP(&dirname, "DIR", "d", "", "decode dir path")
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "decode music file",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if filename == "" && dirname == "" {
			return errors.New("Require a file path")
		}
		if filename != "" {
			_, err := os.Lstat(filename)
			if os.IsNotExist(err) {
				return errors.New("File not found")
			}
		}
		if dirname != "" {
			_, err := os.Lstat(dirname)
			if os.IsNotExist(err) {
				return errors.New("Dir not found")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		decode(args)
	},
}

func decode(args []string) error {
	var strIndex int = strings.LastIndex(filename, ".")
	var fileFormat = filename[strIndex+1 : len(filename)]
	switch fileFormat {
	case "qmcflac":
		util.DecodeQmcFlac(filename)
	case "qmc0", "qmc3":
		util.DecodeQmc0OrQmc3(filename)
	default:
		return errors.New("The file not support")
	}
	return nil
}
