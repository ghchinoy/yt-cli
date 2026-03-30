package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/youtube/v3"
	"github.com/ghchinoy/yt-cli/auth"
	"github.com/ghchinoy/yt-cli/ui"
)

var (
	uploadTitle       string
	uploadDescription string
	dryRun            bool
)

var videoUploadCmd = &cobra.Command{
	Use:   "upload <filepath>",
	Short: "Upload a video to your YouTube channel",
	Long:  `Uploads an mp4 video file to the authenticated YouTube channel. Requires a title. Use --dry-run for mutative safety testing.`,
	Example: `  yt-cli video upload ./my-video.mp4 --title "My Title" --dry-run
  yt-cli video upload ./my-video.mp4 --title "My Title" --description "Some desc"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		if dryRun {
			if viper.GetBool("json") {
				out := map[string]interface{}{
					"status": "dry-run",
					"file":   filename,
					"snippet": map[string]string{
						"title":       uploadTitle,
						"description": uploadDescription,
					},
				}
				b, _ := json.MarshalIndent(out, "", "  ")
				fmt.Println(string(b))
				return
			}
			fmt.Printf("%s\n", ui.Warn.Render("DRY RUN MODE ENABLED"))
			fmt.Printf("Would upload file: %s\n", ui.Command.Render(filename))
			fmt.Printf("Title: %s\n", uploadTitle)
			fmt.Printf("Description: %s\n", uploadDescription)
			return
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Verify the file path and permissions."))
			os.Exit(1)
		}
		defer file.Close()

		srv, err := auth.GetYouTubeService()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, ""))
			os.Exit(1)
		}

		upload := &youtube.Video{
			Snippet: &youtube.VideoSnippet{
				Title:       uploadTitle,
				Description: uploadDescription,
				CategoryId:  "22", // Default to People & Blogs
			},
			Status: &youtube.VideoStatus{PrivacyStatus: "private"}, // Default to private for safety
		}

		call := srv.Videos.Insert([]string{"snippet", "status"}, upload)
		resp, err := call.Media(file).Do()
		if err != nil {
			fmt.Fprintln(os.Stderr, ui.FormatError(err, "Check your internet connection and API quota."))
			os.Exit(1)
		}

		if viper.GetBool("json") {
			b, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(b))
			return
		}

		fmt.Printf("%s Video uploaded successfully!\n", ui.Pass.Render("Success:"))
		fmt.Printf("Video ID: %s\n", ui.ID.Render(resp.Id))
	},
}

func init() {
	videoCmd.AddCommand(videoUploadCmd)
	videoUploadCmd.Flags().StringVarP(&uploadTitle, "title", "t", "", "Title of the video (required)")
	videoUploadCmd.MarkFlagRequired("title")
	videoUploadCmd.Flags().StringVarP(&uploadDescription, "description", "d", "", "Description of the video")
	videoUploadCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate parameters without actually uploading")
}
