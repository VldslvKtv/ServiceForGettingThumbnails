package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"getthumbnails/internal/models"
	"getthumbnails/storage/storage_err"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite3.New"

	// путь до файла БД

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) GetCache(ctx context.Context,
	url string) (string, error) {

	const op = "storage.sqlite.GetCache"

	stmt, err := s.db.Prepare("SELECT id, url FROM thumbnails WHERE url = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, url)

	var thumbnail models.Thumbnail
	err = row.Scan(&thumbnail.ID, &thumbnail.Url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage_err.ErrNotFoundCache)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return thumbnail.Url, nil
}

func (s *Storage) NewThumbnail(ctx context.Context,
	url string) (int64, error) {
	const op = "storage.sqlite.NewThumbnail"

	stmt, err := s.db.Prepare("INSERT INTO thumbnails(url) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, url)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
