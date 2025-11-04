package utils

import(
	"encoding/base64"
)

func Encode(cursor *string) (*string) {
	if cursor == nil {
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(*cursor))

	return &encoded
}

func Decode(cursor *string) (*string, error) {
	if cursor == nil {
		return nil, nil
	}

	decoded_byte, err := base64.StdEncoding.DecodeString(*cursor)
	if (err != nil) {
		return nil, err
	}

	decoded := string(decoded_byte)

	return &decoded, nil
}