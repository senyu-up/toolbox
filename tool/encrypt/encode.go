package encrypt

import "encoding/base64"

func Base64Encode(d []byte) string {
	s := base64.StdEncoding.EncodeToString(d)

	return s
}

func Base64Decode(s []byte) (b []byte, err error) {
	b, err = base64.StdEncoding.DecodeString(string(s))

	return
}

func Base64EncodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64DecodeString(s string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(s)

	return string(decodeString), err
}
