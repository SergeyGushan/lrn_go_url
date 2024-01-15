package url

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
)

func CreateShortURL(longURL string) (string, error) {
	hash := md5.New()
	_, err := io.WriteString(hash, longURL)
	if err != nil {
		return "", fmt.Errorf("i was unable to create a short hash")
	}

	shortCode := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

	return fmt.Sprintf("%s", shortCode), nil
}
