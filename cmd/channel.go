package cmd

import (
	"github.com/spf13/cobra"
)

var channelCmd = &cobra.Command{
	Use:     "channel",
	GroupID: "Info",
	Short:   "Manage and query YouTube channels",
	Long:    `Group of commands to query information and retrieve video lists for YouTube channels.`,
}

func init() {
	rootCmd.AddCommand(channelCmd)
}
