package uiuscraper

import (
	"crypto/sha256"
	"fmt"
)

func GenerateNoticeID(title, date string) string {
	hashInput := fmt.Sprintf("%s-%s", title, date)
	hash := sha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}
