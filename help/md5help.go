package help

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Getmd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	s := hex.EncodeToString(h.Sum(nil))
	fmt.Println(s)
	return s
}
