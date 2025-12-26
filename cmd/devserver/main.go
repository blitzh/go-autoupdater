package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Manifest struct {
	Product     string     `json:"product"`
	Channel     string     `json:"channel"`
	Version     string     `json:"version"`
	PublishedAt string     `json:"published_at"`
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

func main() {
	var (
		addr         = flag.String("addr", "127.0.0.1:8089", "listen address")
		root         = flag.String("root", "./testdata/updates", "root dir containing manifest.json + artifacts")
		manifestPath = flag.String("manifest", "manifest.json", "manifest file name under root")
		gen          = flag.Bool("gen", false, "auto-generate manifest from files under root (simple mode)")
		product      = flag.String("product", "agent", "product name used in generated manifest")
		channel      = flag.String("channel", "stable", "channel used in generated manifest")
		version      = flag.String("version", "0.0.0-dev", "version used in generated manifest")
		notes        = flag.String("notes", "local dev build", "notes used in generated manifest")
	)
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve manifest at /manifest.json
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")

		if *gen {
			m, err := generateManifest(absRoot, *product, *channel, *version, *notes, "http://"+*addr)
			if err != nil {
				http.Error(w, "generate manifest: "+err.Error(), 500)
				return
			}
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			_ = enc.Encode(m)
			return
		}

		p := filepath.Join(absRoot, *manifestPath)
		b, err := os.ReadFile(p)
		if err != nil {
			http.Error(w, "read manifest: "+err.Error(), 500)
			return
		}
		_, _ = w.Write(b)
	})

	// Serve artifacts under /files/<filename>
	mux.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/files/")
		name = filepath.Base(name) // prevent path traversal
		if name == "" || name == "." || name == "/" {
			http.Error(w, "bad file", 400)
			return
		}
		p := filepath.Join(absRoot, name)

		f, err := os.Open(p)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		// set content type by ext (exe)
		ext := filepath.Ext(p)
		if ext != "" {
			if ct := mime.TypeByExtension(ext); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		http.ServeContent(w, r, name, time.Now(), f)
	})

	// Simple index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "go-autoupdater devserver\n\n")
		fmt.Fprintf(w, "Manifest:  http://%s/manifest.json\n", *addr)
		fmt.Fprintf(w, "Artifacts: http://%s/files/<filename>\n\n", *addr)
		fmt.Fprintf(w, "Root dir:  %s\n", absRoot)
		fmt.Fprintf(w, "Mode:      %s\n", map[bool]string{true: "GEN", false: "STATIC"}[*gen])
	})

	s := &http.Server{
		Addr:              *addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("devserver listening on http://%s\n", *addr)
	log.Printf("manifest: http://%s/manifest.json\n", *addr)
	log.Fatal(s.ListenAndServe())
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// generateManifest: simple mode.
// It scans files in root and tries to infer os/arch from filename patterns like:
// agent_1.0.12_windows_amd64.exe
func generateManifest(root, product, channel, version, notes, baseURL string) (*Manifest, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	m := &Manifest{
		Product:     product,
		Channel:     channel,
		Version:     version,
		PublishedAt: time.Now().UTC().Format(time.RFC3339),
		Notes:       notes,
		Artifacts:   []Artifact{},
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "manifest.json" {
			continue
		}

		// try infer os/arch from filename: *_<os>_<arch>(.exe)
		osName, arch := inferOSArch(name)
		if osName == "" || arch == "" {
			continue
		}

		path := filepath.Join(root, name)
		sum, err := fileSHA256(path)
		if err != nil {
			return nil, err
		}

		m.Artifacts = append(m.Artifacts, Artifact{
			OS:     osName,
			Arch:   arch,
			Name:   name,
			URL:    fmt.Sprintf("%s/files/%s", baseURL, name),
			SHA256: sum,
		})
	}

	if len(m.Artifacts) == 0 {
		return nil, fmt.Errorf("no artifacts detected; ensure filenames contain _<os>_<arch>")
	}
	return m, nil
}

func inferOSArch(filename string) (string, string) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.Split(base, "_")
	if len(parts) < 3 {
		return "", ""
	}
	// last two tokens assumed os + arch
	osName := parts[len(parts)-2]
	arch := parts[len(parts)-1]
	// allow windows/linux/darwin
	switch osName {
	case "windows", "linux", "darwin":
	default:
		return "", ""
	}
	return osName, arch
}

func fileSHA256(path string) (string, error) {
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
