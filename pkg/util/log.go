package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	mu   sync.Mutex
	file string
}

func NewLogger(filePath string) *Logger {
	return &Logger{file: filePath}
}

func (l *Logger) Printf(format string, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	line := time.Now().Format("2006-01-02 15:04:05") + " " + fmt.Sprintf(format, a...) + "\n"
	fmt.Print(line)

	if l.file == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(l.file), 0755)
	f, err := os.OpenFile(l.file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		_, _ = f.WriteString(line)
		_ = f.Close()
	}
}
