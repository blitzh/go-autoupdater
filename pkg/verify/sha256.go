package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func FileSHA256Hex(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func VerifyFileSHA256(path string, expectedHex string) error {
	expectedHex = strings.ToLower(strings.TrimSpace(expectedHex))
	expectedHex = strings.TrimPrefix(expectedHex, "sha256:")
	if expectedHex == "" {
		return fmt.Errorf("expected sha256 is empty")
	}
	got, err := FileSHA256Hex(path)
	if err != nil {
		return err
	}
	if got != expectedHex {
		return fmt.Errorf("sha256 mismatch got=%s expected=%s", got, expectedHex)
	}
	return nil
}
