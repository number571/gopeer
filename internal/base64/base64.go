package base64

import (
	"encoding/base64"
)

func GetSizeInBase64(pBytesNum uint64) (uint64, error) {
	if pBytesNum < 2 {
		return 0, ErrBase64Size
	}
	// base64 encoding bytes with add 1/4 bytes of original
	// (-2) is a '=' characters in the suffix of encoding bytes
	return uint64(base64.StdEncoding.DecodedLen(int(pBytesNum))) - 2, nil
}
