package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/luanxuechao/qn-decode/util"
	"github.com/nu11ptr/cmpb"
	"github.com/spf13/cobra"
)

// FilePath file path
var FilePath string
var dirname string

func init() {
	rootCmd.AddCommand(decodeCmd)
	decodeCmd.Flags().StringVarP(&FilePath, "FILE", "f", "", "decode file path")
	decodeCmd.Flags().StringVarP(&dirname, "DIR", "d", "", "decode dir path")
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "decode music file",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if FilePath == "" && dirname == "" {
			return errors.New("Require a file path")
		}
		if FilePath != "" {
			_, err := os.Lstat(FilePath)
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

func decodeFile(filePath string, p *cmpb.Progress) error {
	var strIndex int = strings.LastIndex(filePath, ".")
	var fileFormat = filePath[strIndex+1 : len(filePath)]
	_, fileName := filepath.Split(filePath)

	switch fileFormat {
	case "qmcflac":
		util.DecodeQmcFlac(filePath, fileName, p)
		break
	case "qmc0", "qmc3":
		util.DecodeQmc0OrQmc3(filePath, fileName, p)
		break
	case "ncm":
		util.Dump(filePath, fileName, p)
		break
	default:
		return errors.New("The file not support")
	}
	return nil
}
func decodeDir(p *cmpb.Progress) error {
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
	var FilePathList []string = make([]string, 0)
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
		FilePathList = append(FilePathList, dirname+"/"+fi.Name())
	}

	for _, filePath := range FilePathList {
		decodeFile(filePath, p)
	}
	return nil
}
func decode(args []string) error {
	p := cmpb.New()
	colors := new(cmpb.BarColors)

	colors.Post, colors.KeyDiv, colors.LBracket, colors.RBracket =
		color.HiCyanString, color.HiCyanString, color.HiCyanString, color.HiCyanString

	colors.Key = color.HiBlueString
	colors.Msg, colors.Empty = color.HiYellowString, color.HiYellowString
	colors.Full = color.HiGreenString
	colors.Curr = color.GreenString
	colors.PreBar, colors.PostBar = color.HiMagentaString, color.HiMagentaString
	p.SetColors(colors)
	p.Start()
	p.Wait()
	if FilePath != "" {
		return decodeFile(FilePath, p)
	}
	return decodeDir(p)
}
