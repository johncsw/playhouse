package flow

import (
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"playhouse-server/env"
	"playhouse-server/repo"
	"playhouse-server/util"
	"time"
)

func DeleteAllDataEveryHour() {
	go func() {
		deleteAllData()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				deleteAllData()
			}
		}
	}()
}

func deleteAllData() {
	log.Printf("Start deleting all data. %s", time.Now().Format("2006-01-02 15:04:05"))
	dbErr := repo.NewTransaction(func(tx *gorm.DB) error {
		if err := repo.ChunkRepo().DeleteAll(tx); err != nil {
			return err
		}

		if err := repo.VideoRepo().DeleteAll(tx); err != nil {
			return err
		}

		if err := repo.SessionRepo().DeleteAll(tx); err != nil {
			return err
		}

		if err := repo.UserRepo().DeleteAll(tx); err != nil {
			return err
		}

		return nil
	})

	if dbErr != nil {
		util.LogError(dbErr, "")
		return
	}

	if binErr := deleteChunkBinary(); binErr != nil {
		util.LogError(binErr, "")
		return
	}
	log.Printf("Finish	deleting all data. %s", time.Now().Format("2006-01-02 15:04:05"))
}

func deleteChunkBinary() error {
	d, openErr := os.Open(env.CHUNK_STORAGE_PATH())
	if openErr != nil {
		return openErr
	}
	defer func() {
		if closeErr := d.Close(); closeErr != nil {
			util.LogError(closeErr, "")
		}
	}()

	names, readErr := d.Readdirnames(-1)
	if readErr != nil {
		return readErr
	}

	for _, name := range names {
		if name == ".gitkeep" {
			continue
		}
		if removeErr := os.RemoveAll(filepath.Join(env.CHUNK_STORAGE_PATH(), name)); removeErr != nil {
			return removeErr
		}
	}

	return nil
}
