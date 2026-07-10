package project

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func GenerateGatewayAPIKey() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	key := strings.ToLower(
		base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b),
	)
	return "gw_" + key, nil
}
