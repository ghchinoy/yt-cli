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

var mineInfo bool

var channelInfoCmd = &cobra.Command{
	Use:   "info [channel-id]",
	Short: "Display basic info about a channel",
	Long:  `Retrieves basic metadata about a YouTube channel using its ID. Use --mine to get info about your own channel.`,
	Example: `  yt-cli channel info UC_x5XG1OV2P6uZZ5FSM9Ttw
  yt-cli channel info --mine
  yt-cli channel info UC_x5XG1OV2P6uZZ5FSM9Ttw --json`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !mineInfo && len(args) == 0 {
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

		call := srv.Channels.List([]string{"snippet", "statistics"})
		if mineInfo {
			call = call.Mine(true)
		} else {
			call = call.Id(channelID)
		}
		
		resp, err := call.Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Verify the channel ID is correct or you are authenticated."))
			os.Exit(1)
		}

		if len(resp.Items) == 0 {
			fmt.Fprintln(os.Stderr, ui.FormatError(fmt.Errorf("channel not found"), ""))
			os.Exit(1)
		}

		ch := resp.Items[0]

		if viper.GetBool("json") {
			b, _ := json.MarshalIndent(ch, "", "  ")
			fmt.Println(string(b))
			return
		}

		fmt.Printf("%s %s\n", ui.Accent.Render("Channel Title:"), ch.Snippet.Title)
		fmt.Printf("%s %s\n", ui.Accent.Render("Channel ID:"), ui.ID.Render(ch.Id))
		if ch.Statistics != nil {
			fmt.Printf("%s %d\n", ui.Muted.Render("Subscribers:"), ch.Statistics.SubscriberCount)
			fmt.Printf("%s %d\n", ui.Muted.Render("Total Views:"), ch.Statistics.ViewCount)
		}
	},
}

func init() {
	channelCmd.AddCommand(channelInfoCmd)
	channelInfoCmd.Flags().BoolVarP(&mineInfo, "mine", "m", false, "Get info for the authenticated user's channel")
}
