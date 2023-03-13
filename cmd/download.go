/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
		defer cancel()

		localStore, err := downloader.NewLocalStorage(time.Now().Format("2006-01-02"))
		if err != nil {
			log.Fatal(err)
		}

		streamClient := downloader.NewClient(
			downloader.WithStore(localStore),
		)

		urls := []string{
			"https://stream-us1-alfa.dropcam.com/nexus_aac/88090853f44849c892642d805c42ad9a/playlist.m3u8?public=MGKM8iptgQ", // gym
			// "https://stream-us1-charlie.dropcam.com/nexus_aac/e6163f08d5094f21a6524733ad2c7023/playlist.m3u8?public=zWixbR3lG9", // toybox
		}

		if err := streamClient.Download(ctx, urls); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
