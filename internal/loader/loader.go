package loader

import (
	"context"
	"errors"
	"fmt"
	"getthumbnails/internal/youtube"
	"getthumbnails/storage/storage_err"
	"log/slog"
)

type Load struct {
	log   *slog.Logger
	cache Cache
}

func New(log *slog.Logger,
	cache Cache) *Load {
	return &Load{
		log:   log,
		cache: cache,
	}
}

type Cache interface {
	GetCache(
		ctx context.Context,
		url string,
	) (urlThumb string, err error)

	NewThumbnail(ctx context.Context,
		url string,
	) (urlID int64, err error)
}

var (
	ErrCache     = errors.New("load into YouTube")
	ErrLoadVideo = errors.New("failed to load video")
)

func (l *Load) LoadThumbnail(ctx context.Context,
	url string) (string, error) {

	const op = "loader.LoadThumbnail"

	log := l.log.With(
		slog.String("op", op),
		slog.String("url", url),
	)

	_, err := l.cache.GetCache(ctx, url)
	if err != nil {
		if errors.Is(err, storage_err.ErrNotFoundCache) {
			_, err := l.cache.NewThumbnail(ctx, url)
			if err != nil {
				fmt.Println(err)
				return "", fmt.Errorf("%s: %w", op, ErrCache)
			}

			thunmbnailUrl, err := youtube.GetThumbnail(url)
			if err != nil {
				fmt.Println(err)
				return "", fmt.Errorf("%s: %w", op, ErrLoadVideo)
			}
			log.Info("thumbnail loaded")

			return thunmbnailUrl, nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	thunmbnailUrl, err := youtube.GetThumbnail(url)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("%s: %w", op, ErrLoadVideo)
	}

	log.Info("thumbnail loaded into cache")

	return thunmbnailUrl, nil

}