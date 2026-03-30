package cmd

import (
	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:     "video",
	GroupID: "Info",
	Short:   "Manage and query YouTube videos",
	Long:    `Group of commands to query information, look up channel relationships, and upload YouTube videos.`,
}

func init() {
	rootCmd.AddCommand(videoCmd)
}
