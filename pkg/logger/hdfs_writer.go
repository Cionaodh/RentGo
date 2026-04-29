package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/colinmarc/hdfs/v2"
)

type HDFSRotatingWriter struct {
	client      *hdfs.Client
	serviceName string
	level       string

	mu          sync.Mutex
	currentHour string
	currentFile *hdfs.FileWriter // ← держим файл открытым
}

func NewHDFSRotatingWriter(client *hdfs.Client, serviceName, level string) *HDFSRotatingWriter {
	return &HDFSRotatingWriter{
		client:      client,
		serviceName: serviceName,
		level:       level,
	}
}

func (w *HDFSRotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UTC()
	hourKey := now.Format("20060102-15")

	if w.currentHour != hourKey {
		if err := w.rotate(now); err != nil {
			fmt.Fprintf(os.Stderr, "[hdfs] rotate error [%s]: %v\n", w.level, err)
			return len(p), nil
		}
		w.currentHour = hourKey
	}

	n, err := w.currentFile.Write(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[hdfs] write error [%s]: %v\n", w.level, err)
		return len(p), nil
	}

	// ↓ Сбрасываем буфер — данные сразу видны в HDFS
	if err := w.currentFile.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "[hdfs] flush error [%s]: %v\n", w.level, err)
	}

	return n, nil
}

// rotate — закрывает старый файл и открывает новый
func (w *HDFSRotatingWriter) rotate(t time.Time) error {
	// Закрываем предыдущий если был открыт
	if w.currentFile != nil {
		if err := w.currentFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "[hdfs] warn: close old file [%s]: %v\n", w.level, err)
		}
		w.currentFile = nil
	}

	dir := fmt.Sprintf("/logs/%s/%s", w.serviceName, w.level)
	if err := w.client.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	path := fmt.Sprintf("%s/%s-%s.log", dir, w.serviceName, t.Format("20060102-15"))

	var f *hdfs.FileWriter

	_, statErr := w.client.Stat(path)
	if statErr != nil {
		// Файл не существует — создаём
		created, err := w.client.Create(path)
		if err != nil {
			return fmt.Errorf("create %s: %w", path, err)
		}
		f = created
		fmt.Fprintf(os.Stderr, "[hdfs] new file: %s\n", path)
	} else {
		// Файл существует — открываем для Append
		appended, err := w.client.Append(path)
		if err != nil {
			return fmt.Errorf("append open %s: %w", path, err)
		}
		f = appended
		fmt.Fprintf(os.Stderr, "[hdfs] append to existing: %s\n", path)
	}

	w.currentFile = f
	return nil
}

// Close — вызвать при завершении приложения
func (w *HDFSRotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}
