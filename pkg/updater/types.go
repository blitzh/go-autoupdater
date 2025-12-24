package updater

import "time"

type Manifest struct {
	Product     string     `json:"product"`
	Channel     string     `json:"channel"`
	Version     string     `json:"version"`
	PublishedAt time.Time  `json:"published_at"`
	Notes       string     `json:"notes"`
	Artifacts   []Artifact `json:"artifacts"`
}

type Artifact struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

type CheckResult struct {
	CurrentVersion  string
	RemoteVersion   string
	UpdateAvailable bool
	Artifact        *Artifact
	Notes           string
}

type UpdateResult struct {
	DidUpdate     bool
	OldBackupPath string
	NewBinaryPath string
	RemoteVersion string
}
