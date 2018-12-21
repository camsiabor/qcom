package qencode

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

/*
import "golang.org/x/text/encoding/"

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
*/

func Unicode2StringEx(form string, byteOrder binary.ByteOrder) (to string, err error) {

	form = strings.Replace(form, `\u`, ``, -1)
	bs, err := hex.DecodeString(form)
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, byteOrder, &r)
		to += string(r)
	}
	return
}

func Unicode2String(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func Md5Str(s string) string {
	var md5ctx = md5.New()
	md5ctx.Write([]byte(s))
	var sum = md5ctx.Sum(nil)
	return hex.EncodeToString(sum)
}
