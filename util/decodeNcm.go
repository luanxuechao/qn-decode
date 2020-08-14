package util

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

// Album model
type Album struct {
	ID       float64 `json:"albumId"`
	Name     string  `json:"album"`
	CoverURL string  `json:"albumPic"`
}

// Meta model
type Meta struct {
	ID       float64 `json:"musicId"`
	Name     string  `json:"musicName"`
	Album    *Album  `json:"-"`
	BitRate  float64 `json:"bitrate"`
	Duration float64 `json:"duration"`
	Format   string  `json:"format"`
	Comment  string  `json:"-"`
}

var (
	aesCoreKey   = []byte{0x68, 0x7A, 0x48, 0x52, 0x41, 0x6D, 0x73, 0x6F, 0x35, 0x6B, 0x49, 0x6E, 0x62, 0x61, 0x78, 0x57}
	aesModifyKey = []byte{0x23, 0x31, 0x34, 0x6C, 0x6A, 0x6B, 0x5F, 0x21, 0x5C, 0x5D, 0x26, 0x30, 0x55, 0x3C, 0x27, 0x28}
)

func buildKeyBox(key []byte) []byte {
	box := make([]byte, 256)
	for i := 0; i < 256; i++ {
		box[i] = byte(i)
	}
	keyLen := byte(len(key))
	var c, lastByte, keyOffset byte
	for i := 0; i < 256; i++ {
		c = (box[i] + lastByte + key[keyOffset]) & 0xff
		keyOffset++
		if keyOffset >= keyLen {
			keyOffset = 0
		}
		box[i], box[c] = box[c], box[i]
		lastByte = c
	}
	return box
}

func fixBlockSize(src []byte) []byte {
	return src[:len(src)/aes.BlockSize*aes.BlockSize]
}

func _PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func decryptAes128Ecb(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	dataLen := len(data)
	decrypted := make([]byte, dataLen)
	bs := block.BlockSize()
	for i := 0; i <= dataLen-bs; i += bs {
		block.Decrypt(decrypted[i:i+bs], data[i:i+bs])
	}
	return _PKCS7UnPadding(decrypted), nil
}

func readUint32(rBuf []byte, fp *os.File) (uint32, error) {
	if _, err := fp.Read(rBuf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(rBuf), nil
}

// NCMFile check file
func NCMFile(fp *os.File) (bool, error) {
	// Jump to begin of file
	if _, err := fp.Seek(0, io.SeekStart); err != nil {
		return false, err
	}

	var rBuf = make([]byte, 4)
	uLen, _ := readUint32(rBuf, fp)

	if uLen != 0x4e455443 {
		return false, fmt.Errorf("%s isn't netease cloud music copyright file", fp.Name())
	}

	uLen, _ = readUint32(rBuf, fp)
	if uLen != 0x4d414446 {
		return false, fmt.Errorf("%s isn't netease cloud music copyright file", fp.Name())
	}

	return true, nil
}

// DumpMeta dump meta info
func DumpMeta(fp *os.File) (Meta, error) {
	// detect whether ncm file
	if result, err := NCMFile(fp); err != nil || !result {
		return Meta{}, err
	}

	// jump over the magic head(4*2) and the gap(2).
	if _, err := fp.Seek(4*2+2, io.SeekStart); err != nil {
		return Meta{}, err
	}

	// whether decode key is successful
	if _, err := Decode(fp); err != nil {
		return Meta{}, err
	}

	var rBuf = make([]byte, 4)
	uLen, err := readUint32(rBuf, fp)
	if err != nil {
		return Meta{}, err
	}

	if uLen <= 0 {
		format := "flac"
		if info, err := fp.Stat(); err != nil && info.Size() < int64(math.Pow(1024, 2)*16) {
			format = "mp3"
		}
		return Meta{
			Format: format,
		}, nil
	}

	var modifyData = make([]byte, uLen)
	if _, err = fp.Read(modifyData); err != nil {
		return Meta{}, err
	}

	for i := range modifyData {
		modifyData[i] ^= 0x63
	}

	// 22 = len(`163 key(Don't modify):`)
	deModifyData := make([]byte, base64.StdEncoding.DecodedLen(len(modifyData)-22))
	if _, err = base64.StdEncoding.Decode(deModifyData, modifyData[22:]); err != nil {
		return Meta{}, err
	}

	deData, err := decryptAes128Ecb(aesModifyKey, fixBlockSize(deModifyData))
	if err != nil {
		return Meta{}, err
	}

	// 6 = len("music:")
	deData = deData[6:]

	var album Album
	if err := json.Unmarshal(deData, &album); err != nil {
		return Meta{}, err
	}

	var meta Meta
	if err := json.Unmarshal(deData, &meta); err != nil {
		return Meta{}, err
	}

	meta.Album = &album
	meta.Comment = string(modifyData)
	return meta, nil
}

// DumpCover dump cover info
func DumpCover(fp *os.File) ([]byte, error) {
	if result, err := NCMFile(fp); !result || err != nil {
		return nil, err
	}

	if _, err := DumpMeta(fp); err != nil {
		return nil, err
	}

	// jump over crc32 check
	if _, err := fp.Seek(9, io.SeekCurrent); err != nil {
		return nil, err
	}

	var rBuf = make([]byte, 4)
	if imgLen, err := readUint32(rBuf, fp); err != nil {
		return nil, err
	} else {
		var imgData = make([]byte, imgLen)
		if _, err = fp.Read(imgData); err != nil {
			return nil, err
		}

		return imgData, nil
	}
}

// Decode  info
func Decode(fp *os.File) ([]byte, error) {
	// detect whether ncm file
	if result, err := NCMFile(fp); err != nil || !result {
		return nil, err
	}

	// jump over the magic head(4*2) and the gap(2).
	if _, err := fp.Seek(4*2+2, io.SeekStart); err != nil {
		return nil, err
	}

	var rBuf = make([]byte, 4)

	uLen, err := readUint32(rBuf, fp)
	if err != nil {
		return nil, err
	}

	var keyData = make([]byte, uLen)
	if _, err := fp.Read(keyData); err != nil {
		return nil, err
	}

	for i := range keyData {
		keyData[i] ^= 0x64
	}

	deKeyData, err := decryptAes128Ecb(aesCoreKey, fixBlockSize(keyData))
	if err != nil {
		return nil, err
	}

	// 17 = len("neteasecloudmusic")
	return deKeyData[17:], nil
}
func isFlac(fp *os.File) (bool, error) {
	if result, err := NCMFile(fp); err != nil || !result {
		return false, err
	}

	// jump over the magic head(4*2) and the gap(2).
	if _, err := fp.Seek(4*2+2, io.SeekStart); err != nil {
		return false, err
	}

	var rBuf = make([]byte, 4)
	uLen, err := readUint32(rBuf, fp)
	if err != nil {
		return false, err
	}
	if uLen <= 0 {
		if info, err := fp.Stat(); err != nil && info.Size() < int64(math.Pow(1024, 2)*16) {
			return false, nil
		}
	}
	return true, nil
}

// Dump  info
func Dump(filename string) ([]byte, error) {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, errors.New("The file not support")
	}
	if result, err := NCMFile(fp); !result || err != nil {
		return nil, err
	}

	// whether decode key is successful
	deKeyData, err := Decode(fp)
	if err != nil {
		return nil, err
	}
	if _, err := DumpCover(fp); err != nil {
		return nil, err
	}

	box := buildKeyBox(deKeyData)
	n := 0x8000
	var writer bytes.Buffer

	var tb = make([]byte, n)
	for {
		if _, err := fp.Read(tb); err != nil {
			break // read EOF
		}

		for i := 0; i < n; i++ {
			j := byte((i + 1) & 0xff)
			tb[i] ^= box[(box[j]+box[(box[j]+j)&0xff])&0xff]
		}

		writer.Write(tb) // write to memory
	}
	strIndex := strings.LastIndex(filename, ".")

	if strIndex == -1 {
		return nil, errors.New("file not expected")
	}
	isFlacFormat, err := isFlac(fp)
	newFile := ""
	if isFlacFormat {
		newFile = filename[0:strIndex] + ".flac"
	} else {
		newFile = filename[0:strIndex] + ".mp3"
	}
	err2 := ioutil.WriteFile(newFile, writer.Bytes(), 0666)
	fmt.Println(newFile)
	if err2 != nil {
		fmt.Printf(err2.Error())
	}

	return nil, nil
}
