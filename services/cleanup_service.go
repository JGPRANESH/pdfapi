package services

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func StartCleanupWorker() {

	ticker := time.NewTicker(30 * time.Second)

	go func() {

		defer ticker.Stop()

		for range ticker.C {

			cleanupDirectory("tmp")
			cleanupDirectory("uploads")
		}
	}()
}

func cleanupDirectory(dir string) {

	files, err := os.ReadDir(dir)

	if err != nil {
		return
	}

	for _, file := range files {

		path := filepath.Join(dir, file.Name())

		info, err := file.Info()

		if err != nil {
			continue
		}

		// Delete files older than 30 seconds
		if time.Since(info.ModTime()) > 30*time.Second {

			err := os.RemoveAll(path)

			if err == nil {
				log.Printf("Deleted: %s", path)
			}
		}
	}
}
