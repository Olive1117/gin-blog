package utils

import (
	"encoding/binary"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func EncodeByOBID(id int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))
	var start int
	for start = 0; start < len(buf); start++ {
		if buf[start] != 0 {
			break
		}
	}
	return "OB" + base58.Encode(buf[start:])
}
func DecodeByOBID(obID string) int64 {
	cleanCode := strings.TrimPrefix(obID, "OB")
	buf := base58.Decode(cleanCode)
	if len(buf) == 0 {
		return 0
	}
	fullBuf := make([]byte, 8)
	copy(fullBuf[8-len(buf):], buf)

	return int64(binary.BigEndian.Uint64(fullBuf))
}
