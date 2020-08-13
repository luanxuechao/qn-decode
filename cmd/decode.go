package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
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

func decodeFile(fileName string) error {
	var strIndex int = strings.LastIndex(fileName, ".")
	var fileFormat = fileName[strIndex+1 : len(fileName)]
	fmt.Println(fileName)
	switch fileFormat {
	case "qmcflac":
		util.DecodeQmcFlac(filename)
		break
	case "qmc0", "qmc3":
		util.DecodeQmc0OrQmc3(filename)
		break
	case "ncm":
		util.Dump(filename)
		break
	default:
		return errors.New("The file not support")
	}
	return nil
}
func decodeDir() error {
	s, err := os.Stat(dirname)
	if err != nil {
		return errors.New("The dir not found")
	}
	if !s.IsDir() {
		return errors.New("The dir is not a folder")
	}
	rd, err := ioutil.ReadDir(dirname)
	if err != nil {
		return errors.New(err.Error())
	}
	var filenameList []string = make([]string, 0)
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		var strIndex int = strings.LastIndex(name, ".")
		var fileFormat = name[strIndex+1 : len(name)]
		if fileFormat != "qmcflac" && fileFormat != "qmc0" && fileFormat != "qmc3" && fileFormat != "ncm" {
			continue
		}
		fmt.Println(len(dirname + "/" + fi.Name()))
		fmt.Println(len("/Users/xuechaoluan/Downloads/OurSong.ncm"))
		filenameList = append(filenameList, dirname+"/"+fi.Name())
	}

	for _, fileName := range filenameList {
		decodeFile(fileName)
	}
	return nil
}
func decode(args []string) error {
	if filename != "" {
		return decodeFile(filename)
	}
	return decodeDir()
}
