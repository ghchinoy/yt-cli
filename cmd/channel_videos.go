// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/ghchinoy/yt-cli/auth"
	"github.com/ghchinoy/yt-cli/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/youtube/v3"
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
	Example: `  yt-cli channel videos UC_x5XG1OV2P6uZZ5FSM9Ttw
  yt-cli channel videos --mine --limit 5
  yt-cli channel videos --mine --detailed --filter public,unlisted
  yt-cli channel videos UC_x5XG1OV2P6uZZ5FSM9Ttw --json`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		channelID, err := resolveChannelTarget(args, mineVideos)
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Pass a channel ID positional argument or --mine."))
			_ = cmd.Usage()
			os.Exit(1)
		}

		ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		srv, err := auth.NewYouTubeService(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, authSetupHint))
			os.Exit(1)
		}

		chCall := srv.Channels.List([]string{"contentDetails"}).Id(channelID)
		if mineVideos {
			chCall = srv.Channels.List([]string{"contentDetails"}).Mine(true)
		}
		chResp, err := chCall.Context(ctx).Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Verify the channel ID is correct or you are authenticated."))
			os.Exit(1)
		}
		if len(chResp.Items) == 0 {
			fmt.Fprintln(os.Stderr, ui.FormatError(errors.New("channel not found"), ""))
			os.Exit(1)
		}

		ch := chResp.Items[0]
		if ch.ContentDetails == nil || ch.ContentDetails.RelatedPlaylists == nil || ch.ContentDetails.RelatedPlaylists.Uploads == "" {
			fmt.Fprintln(os.Stderr, ui.FormatError(errors.New("uploads playlist not found for channel"), ""))
			os.Exit(1)
		}
		uploadsPlaylistID := ch.ContentDetails.RelatedPlaylists.Uploads

		playlistCall := srv.PlaylistItems.List([]string{"snippet", "status"}).PlaylistId(uploadsPlaylistID).MaxResults(limit)
		playlistResp, err := playlistCall.Context(ctx).Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Could not retrieve playlist items. Note: filtering private videos requires owner authentication."))
			os.Exit(1)
		}

		items := filterPlaylistItems(playlistResp.Items, filter)

		if viper.GetBool("json") {
			if err := printJSON(cmd.OutOrStdout(), items); err != nil {
				fmt.Fprintln(os.Stderr, ui.FormatError(err, ""))
				os.Exit(1)
			}
			return
		}

		displayID := channelID
		if mineVideos {
			displayID = ch.Id
		}

		fmt.Printf("%s Videos for %s (Limit: %d)\n", ui.Accent.Render("Latest"), ui.ID.Render(displayID), limit)
		for _, item := range items {
			if item.Snippet == nil {
				continue
			}
			videoID := ""
			if item.Snippet.ResourceId != nil {
				videoID = item.Snippet.ResourceId.VideoId
			}
			fmt.Printf("- %s [%s]\n", ui.Command.Render(item.Snippet.Title), ui.ID.Render(videoID))
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

func filterPlaylistItems(items []*youtube.PlaylistItem, filterCSV string) []*youtube.PlaylistItem {
	if strings.TrimSpace(filterCSV) == "" {
		return items
	}
	statusMap := make(map[string]bool)
	for _, s := range strings.Split(filterCSV, ",") {
		if trimmed := strings.TrimSpace(s); trimmed != "" {
			statusMap[trimmed] = true
		}
	}

	var filtered []*youtube.PlaylistItem
	for _, item := range items {
		if item != nil && item.Status != nil && statusMap[item.Status.PrivacyStatus] {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func truncate(s string, maxRunes int) string {
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "..."
}

func init() {
	channelCmd.AddCommand(channelVideosCmd)
	channelVideosCmd.Flags().Int64VarP(&limit, "limit", "l", 10, "Maximum number of videos to return")
	channelVideosCmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed information")
	channelVideosCmd.Flags().StringVarP(&filter, "filter", "f", "", "Comma-separated list of privacy statuses to include (e.g., public,unlisted,private)")
	channelVideosCmd.Flags().BoolVarP(&mineVideos, "mine", "m", false, "Get videos for the authenticated user's channel")
}
