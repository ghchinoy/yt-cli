package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/youtube/v3"
	"github.com/ghchinoy/yt-cli/auth"
	"github.com/ghchinoy/yt-cli/ui"
)

var (
	limit      int64
	detailed   bool
	filter     string
	mineVideos bool
)

var channelVideosCmd = &cobra.Command{
	Use:   "videos [channel-id]",
	Short: "List videos for a channel",
	Long:  `Retrieves a list of videos for a specified YouTube channel. You can limit the results, get detailed views, and filter videos by privacy status (if authenticated as the channel owner). Use --mine to get videos for your own channel.`,
	Example: `  yt-cli channel videos UC...
  yt-cli channel videos --mine --limit 5
  yt-cli channel videos --mine --detailed --filter public,unlisted
  yt-cli channel videos UC... --json`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !mineVideos && len(args) == 0 {
			fmt.Fprintln(os.Stderr, ui.FormatError(fmt.Errorf("requires either a channel ID or the --mine flag"), ""))
			cmd.Usage()
			os.Exit(1)
		}

		channelID := ""
		if len(args) > 0 {
			channelID = args[0]
		}

		srv, err := auth.GetYouTubeService()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, ""))
			os.Exit(1)
		}

		// First, get the "uploads" playlist ID for the channel
		chCall := srv.Channels.List([]string{"contentDetails"})
		if mineVideos {
			chCall = chCall.Mine(true)
		} else {
			chCall = chCall.Id(channelID)
		}
		chResp, err := chCall.Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Verify the channel ID is correct or you are authenticated."))
			os.Exit(1)
		}
		if len(chResp.Items) == 0 {
			fmt.Fprintln(os.Stderr, ui.FormatError(fmt.Errorf("channel not found"), ""))
			os.Exit(1)
		}

		uploadsPlaylistID := chResp.Items[0].ContentDetails.RelatedPlaylists.Uploads

		// Now fetch the playlist items
		playlistCall := srv.PlaylistItems.List([]string{"snippet", "status"}).PlaylistId(uploadsPlaylistID).MaxResults(limit)
		playlistResp, err := playlistCall.Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Could not retrieve playlist items. Note: filtering private videos requires owner authentication."))
			os.Exit(1)
		}

		// Filter by privacy if requested
		var items = playlistResp.Items
		if filter != "" {
			allowedStatuses := strings.Split(filter, ",")
			statusMap := make(map[string]bool)
			for _, s := range allowedStatuses {
				statusMap[strings.TrimSpace(s)] = true
			}

			var filtered []*youtube.PlaylistItem
			for _, item := range items {
				if item.Status != nil && statusMap[item.Status.PrivacyStatus] {
					filtered = append(filtered, item)
				}
			}
			items = filtered
		}

		if viper.GetBool("json") {
			b, _ := json.MarshalIndent(items, "", "  ")
			fmt.Println(string(b))
			return
		}

		displayID := channelID
		if mineVideos {
			displayID = chResp.Items[0].Id
		}

		fmt.Printf("%s Videos for %s (Limit: %d)\n", ui.Accent.Render("Latest"), ui.ID.Render(displayID), limit)
		for _, item := range items {
			fmt.Printf("- %s [%s]\n", ui.Command.Render(item.Snippet.Title), ui.ID.Render(item.Snippet.ResourceId.VideoId))
			if detailed {
				fmt.Printf("  %s %s\n", ui.Muted.Render("Published At:"), item.Snippet.PublishedAt)
				privacy := "unknown"
				if item.Status != nil {
					privacy = item.Status.PrivacyStatus
				}
				fmt.Printf("  %s %s\n", ui.Muted.Render("Privacy:"), privacy)
				fmt.Printf("  %s %s\n", ui.Muted.Render("Desc:"), truncate(item.Snippet.Description, 50))
			}
		}
	},
}

func truncate(s string, l int) string {
	if len(s) > l {
		return s[:l] + "..."
	}
	return s
}

func init() {
	channelCmd.AddCommand(channelVideosCmd)
	channelVideosCmd.Flags().Int64VarP(&limit, "limit", "l", 10, "Maximum number of videos to return")
	channelVideosCmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed information")
	channelVideosCmd.Flags().StringVarP(&filter, "filter", "f", "", "Comma-separated list of privacy statuses to include (e.g., public,unlisted,private)")
	channelVideosCmd.Flags().BoolVarP(&mineVideos, "mine", "m", false, "Get videos for the authenticated user's channel")
}
