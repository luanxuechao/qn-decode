## qn-decode
A command tool for transfering  `qmcflac`|`qmc0`|`qmc3`|`ncm` to `mp3` or `flac`.

The repo is used for learning, if there is any infringement, please contact the author to delete


## Installing
### Homebrew
```
$ brew tap luanxuechao/tools
$ brew install qn-decode
```

### GO
Using `qn-decode` is easy. First, use go get to install the latest version of the library. This command will install the `qn-decode` generator executable along with the library and its dependencies:
```
go get github.com/luanxuechao/qn-decode
```

## Usage
```
A command tool for transfering  'qmcflac'
        |'qmc0'|'qmc3'|'ncm' to 'mp3' or 'flac'.

Usage:
  qn-decode [command]

Available Commands:
  decode      decode music file
  help        Help about any command
  version     Print the version number of qn-decode

Flags:
      --config string   config file (default is $HOME/.qn-decode.yaml)
  -h, --help            help for qn-decode
  -t, --toggle          Help message for toggle

Use "qn-decode [command] --help" for more information about a command.
```

### Reference
 - https://github.com/MBearo/qmcdump
 - https://github.com/yoki123/ncmdump
### Example
```
$  decode -d /Users/xuechaoluan/Downloads
```
```
$  decode -f /Users/xuechaoluan/Downloads/xxxx.qmc3
```
