package utils

import (
	"net/url"
	"path/filepath"
)

func CanonicalizeURL(urlstr string) (string, error) {
	parsed, err := url.Parse(urlstr)
	if err != nil {
		return "", err
	}

	if parsed.Scheme == "" {
		path, err := filepath.Abs(parsed.Path)
		if err != nil {
			return "", err
		}

		parsed.Scheme = "file"
		parsed.Path = path
	}

	return parsed.String(), nil
}

func Base58Encode(src []byte) string {
	alphabet := []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	bytes := make([]byte, 0)

	leadingZeros := 0
	for _, b := range src {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	for _, b := range src {
		carry := int(b)
		for j := 0; carry != 0 || j < len(bytes); j++ {
			if j == len(bytes) {
				carry += 0
			} else {
				carry += int(bytes[j]) << 8
			}

			if j == len(bytes) {
				bytes = append(bytes, byte(carry%58))
			} else {
				bytes[j] = byte(carry % 58)
			}

			carry /= 58
		}
	}

	str := ""
	for i := 0; i < leadingZeros+len(bytes); i++ {
		if i < leadingZeros {
			str += string(alphabet[0])
		} else {
			str += string(alphabet[int(bytes[len(bytes)+leadingZeros-i-1])])
		}
	}

	return str
}
