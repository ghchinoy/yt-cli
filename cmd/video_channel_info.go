package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ghchinoy/yt-cli/auth"
	"github.com/ghchinoy/yt-cli/ui"
)

var videoChannelInfoCmd = &cobra.Command{
	Use:   "channel-info <video-id>",
	Short: "Get channel info from a video ID",
	Long:  `Retrieves the basic channel information (ID and title) for the channel that uploaded the specified video.`,
	Example: `  yt-cli video channel-info dQw4w9WgXcQ
  yt-cli video channel-info dQw4w9WgXcQ --json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		videoID := args[0]

		srv, err := auth.GetYouTubeService()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, ""))
			os.Exit(1)
		}

		call := srv.Videos.List([]string{"snippet"}).Id(videoID)
		resp, err := call.Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Verify the video ID is correct."))
			os.Exit(1)
		}

		if len(resp.Items) == 0 {
			fmt.Fprintln(os.Stderr, ui.FormatError(fmt.Errorf("video not found"), ""))
			os.Exit(1)
		}

		vid := resp.Items[0]

		if viper.GetBool("json") {
			b, _ := json.MarshalIndent(map[string]string{
				"channelId":    vid.Snippet.ChannelId,
				"channelTitle": vid.Snippet.ChannelTitle,
			}, "", "  ")
			fmt.Println(string(b))
			return
		}

		fmt.Printf("%s %s\n", ui.Accent.Render("Video Title:"), vid.Snippet.Title)
		fmt.Printf("%s %s\n", ui.Accent.Render("Channel Title:"), vid.Snippet.ChannelTitle)
		fmt.Printf("%s %s\n", ui.Accent.Render("Channel ID:"), ui.ID.Render(vid.Snippet.ChannelId))
	},
}

func init() {
	videoCmd.AddCommand(videoChannelInfoCmd)
}
