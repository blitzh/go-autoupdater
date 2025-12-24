package updater

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/blitzh/go-autoupdater/pkg/apply"
	"github.com/blitzh/go-autoupdater/pkg/service"
	"github.com/blitzh/go-autoupdater/pkg/util"
	"github.com/blitzh/go-autoupdater/pkg/verify"
)

type Source interface {
	Fetch(ctx context.Context) (*Manifest, error)
}

type Config struct {
	CurrentVersion string

	InstallDir string
	ExeName    string // agent.exe / agent
	Channel    string
	Product    string

	Source  Source
	Service service.Controller
	Applier apply.Applier

	// download behavior
	UserAgent string
	MinBytes  int64

	// logging
	Logger  *util.Logger
	LogFile string
}

type Updater struct {
	cfg Config
}

func New(cfg Config) *Updater {
	if cfg.UserAgent == "" {
		cfg.UserAgent = "portable-updater/1.0"
	}
	if cfg.MinBytes == 0 {
		cfg.MinBytes = 32 * 1024
	}
	if cfg.Service == nil {
		cfg.Service = service.NoopController{}
	}
	if cfg.Logger == nil && cfg.LogFile != "" {
		cfg.Logger = util.NewLogger(cfg.LogFile)
	}
	return &Updater{cfg: cfg}
}

func (u *Updater) logf(format string, a ...any) {
	if u.cfg.Logger != nil {
		u.cfg.Logger.Printf(format, a...)
	}
}

func (u *Updater) currentPath() string {
	return filepath.Join(u.cfg.InstallDir, u.cfg.ExeName)
}

func (u *Updater) stagingPaths() (newPath, oldPath string) {
	cur := u.currentPath()
	trim := strings.TrimSuffix(cur, ".exe")
	if runtime.GOOS == "windows" {
		return trim + ".new.exe", trim + ".old.exe"
	}
	return cur + ".new", cur + ".old"
}

func (u *Updater) Check(ctx context.Context) (*CheckResult, error) {
	if u.cfg.Source == nil {
		return nil, errors.New("Source is nil")
	}
	m, err := u.cfg.Source.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	if u.cfg.Product != "" && m.Product != "" && u.cfg.Product != m.Product {
		// not fatal; just warn in notes
	}
	if u.cfg.Channel != "" && m.Channel != "" && u.cfg.Channel != m.Channel {
		// still ok; but likely user should point to matching manifest
	}

	a := selectArtifact(m, runtime.GOOS, runtime.GOARCH)
	res := &CheckResult{
		CurrentVersion:  u.cfg.CurrentVersion,
		RemoteVersion:   m.Version,
		Notes:           m.Notes,
		Artifact:        a,
		UpdateAvailable: false,
	}
	if a == nil {
		return res, fmt.Errorf("no artifact for os=%s arch=%s", runtime.GOOS, runtime.GOARCH)
	}

	// If no current version provided, always say update available (caller can decide)
	if strings.TrimSpace(u.cfg.CurrentVersion) == "" {
		res.UpdateAvailable = true
		return res, nil
	}

	if CompareVersion(u.cfg.CurrentVersion, m.Version) < 0 {
		res.UpdateAvailable = true
	}
	return res, nil
}

func selectArtifact(m *Manifest, osName, arch string) *Artifact {
	for i := range m.Artifacts {
		a := m.Artifacts[i]
		if strings.EqualFold(a.OS, osName) && strings.EqualFold(a.Arch, arch) {
			return &a
		}
	}
	return nil
}

func (u *Updater) Update(ctx context.Context) (*UpdateResult, error) {
	chk, err := u.Check(ctx)
	if err != nil {
		return nil, err
	}
	if !chk.UpdateAvailable {
		return &UpdateResult{DidUpdate: false, RemoteVersion: chk.RemoteVersion}, nil
	}
	if chk.Artifact == nil {
		return nil, fmt.Errorf("artifact is nil")
	}

	newPath, oldPath := u.stagingPaths()
	curPath := u.currentPath()

	u.logf("update available: %s -> %s", chk.CurrentVersion, chk.RemoteVersion)
	u.logf("downloading: %s", chk.Artifact.URL)

	// Download to staging newPath
	if err := util.DownloadToFile(ctx, chk.Artifact.URL, newPath, u.cfg.UserAgent, u.cfg.MinBytes); err != nil {
		return nil, err
	}
	u.logf("downloaded to: %s", newPath)

	// Verify SHA256 (required)
	if err := verify.VerifyFileSHA256(newPath, chk.Artifact.SHA256); err != nil {
		return nil, err
	}
	u.logf("sha256 verified")

	// Apply swap
	if u.cfg.Applier == nil {
		return nil, fmt.Errorf("Applier is nil")
	}

	oldBackup, err := u.cfg.Applier.Apply(ctx, u.cfg.Service, curPath, newPath, oldPath)
	if err != nil {
		return nil, err
	}

	u.logf("apply ok, old backup: %s", oldBackup)

	return &UpdateResult{
		DidUpdate:     true,
		OldBackupPath: oldBackup,
		NewBinaryPath: curPath,
		RemoteVersion: chk.RemoteVersion,
	}, nil
}

// Helper: convenience update with a hard deadline
func (u *Updater) UpdateWithTimeout(timeout time.Duration) (*UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return u.Update(ctx)
}
