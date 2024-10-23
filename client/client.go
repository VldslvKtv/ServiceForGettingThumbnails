package main

import (
	"context"
	_ "flag"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	thumbnail1 "getthumbnails/gen"
	"getthumbnails/internal/config"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const grpcHost = "localhost"

func main() {
	var serverAddr string
	var async bool

	var rootCmd = &cobra.Command{
		Use:   "thumbnail-client",
		Short: "A client for downloading YouTube thumbnails",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalf("Please provide video URLs")
			}

			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			client := thumbnail1.NewThumbnailClient(conn)

			if async {
				var wg sync.WaitGroup
				for _, url := range args {
					wg.Add(1)
					go func(videoURL string) {
						defer wg.Done()
						downloadThumbnail(client, videoURL)
					}(url)
				}
				wg.Wait()
			} else {
				for _, url := range args {
					downloadThumbnail(client, url)
				}
			}
		},
	}

	cfg := config.MustLoadPath("../config/local.yaml")

	rootCmd.PersistentFlags().StringVarP(&serverAddr, "server", "s", grpcAddres(cfg), "The server address in the format of host:port")
	rootCmd.PersistentFlags().BoolVarP(&async, "async", "a", false, "Download thumbnails asynchronously")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("cmd failed: %v", err)
	}
}

func grpcAddres(cfd *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfd.GRPC.Port))
}

func downloadThumbnail(client thumbnail1.ThumbnailClient, videoURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req := &thumbnail1.ThumbnailRequest{Url: videoURL}
	res, err := client.GetThumbnail(ctx, req)
	if err != nil {
		log.Printf("could not get thumbnail for %s: %v", videoURL, err)
		return
	}

	log.Printf("Thumbnail for %s received: %s\n", videoURL, res.Thumbnail)
}
