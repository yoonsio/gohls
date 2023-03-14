/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/yoonsio/gohls/downloader"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download HLS streams",
	Long:  `Download HLS streams`,
	Run: func(cmd *cobra.Command, args []string) {

		duration, err := time.ParseDuration(fDuration)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			log.Fatal(errors.New("no urls supplied"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		localStore, err := downloader.NewLocalStorage(outPath + "/" + time.Now().Format("2006-01-02"))
		if err != nil {
			log.Fatal(err)
		}

		streamClient := downloader.NewClient(
			downloader.WithStore(localStore),
		)

		for _, arg := range args {
			if _, err := url.ParseRequestURI(arg); err != nil {
				log.Fatal(err)
			}
		}

		urls := args

		if err := streamClient.Download(ctx, urls); err != nil {
			log.Fatal(err)
		}
	},
}

var (
	fDuration string
	outPath   string
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&outPath, "out", "o", "", "output path")
	downloadCmd.MarkFlagRequired("out")
	downloadCmd.Flags().StringVarP(&fDuration, "duration", "d", "12h", "how long gohls should record before terminating")
}
