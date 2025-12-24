package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blitzh/go-autoupdater/pkg/updater"
)

type HTTPManifestSource struct {
	ManifestURL string
	Timeout     time.Duration
	UserAgent   string
}

func NewHTTPManifestSource(url string) *HTTPManifestSource {
	return &HTTPManifestSource{
		ManifestURL: url,
		Timeout:     30 * time.Second,
		UserAgent:   "portable-updater/1.0",
	}
}

func (s *HTTPManifestSource) Fetch(ctx context.Context) (*updater.Manifest, error) {
	client := &http.Client{Timeout: s.Timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.ManifestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("manifest http status: %s", resp.Status)
	}

	var m updater.Manifest
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
