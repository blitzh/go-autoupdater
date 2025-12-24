package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadToFile(ctx context.Context, url, dst string, userAgent string, minBytes int64) error {
	tmp := dst + ".part"
	_ = os.Remove(tmp)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download http status: %s", resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	_ = f.Sync()

	if minBytes > 0 && n < minBytes {
		return fmt.Errorf("downloaded file too small: %d bytes", n)
	}

	return RenameWithRetry(tmp, dst, 30, 250*time.Millisecond)
}

func RenameWithRetry(from, to string, retries int, delay time.Duration) error {
	var last error
	for i := 0; i < retries; i++ {
		last = os.Rename(from, to)
		if last == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("rename failed: %w (from=%s to=%s)", last, from, to)
}

func RemoveWithRetry(path string, retries int, delay time.Duration) error {
	var last error
	for i := 0; i < retries; i++ {
		err := os.Remove(path)
		if err == nil || errors.Is(err, os.ErrNotExist) {
			return nil
		}
		last = err
		time.Sleep(delay)
	}
	return fmt.Errorf("remove failed: %w (path=%s)", last, path)
}
