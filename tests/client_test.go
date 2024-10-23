package tests

import (
	"context"
	_ "fmt"
	thumbnail1 "getthumbnails/gen"
	"log/slog"
	_ "net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestIntegrationDownloadThumbnail(t *testing.T) {
	serverAddr := "localhost:5001"

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := thumbnail1.NewThumbnailClient(conn)

	testCases := []struct {
		name      string
		videoURL  string
		expectErr bool
	}{
		{"Valid URL", "https://www.youtube.com/watch?v=P_SXTUiA-9Y", false},
		{"Valid URL 2", "https://www.youtube.com/watch?v=5ClH8EZu5Ug", false},
		{"Valid URL 3 into cache", "https://www.youtube.com/watch?v=5ClH8EZu5Ug", false},
		{"Invalid URL", "invalid-url", true},
		{"Invalid URL (empty)", "invalid-url-2", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			req := &thumbnail1.ThumbnailRequest{Url: tc.videoURL}
			res, err := client.GetThumbnail(ctx, req)

			if tc.expectErr {
				require.Error(t, err)
				slog.Warn("Test failed for URL: %s, error: %v", tc.videoURL, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, res.Thumbnail)
				slog.Info("Test passed for URL: %s, thumbnail: %s", tc.videoURL, res.Thumbnail)
			}
		})
	}
}
