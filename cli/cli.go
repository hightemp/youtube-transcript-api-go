package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hightemp/youtube-transcript-api-go/api"
	"github.com/hightemp/youtube-transcript-api-go/formatters"
)

// CLI represents the command-line interface
type CLI struct {
	api *api.YouTubeTranscriptApi
}

// NewCLI creates a new instance of CLI
func NewCLI() *CLI {
	return &CLI{
		api: api.NewYouTubeTranscriptApi(),
	}
}

// Run executes the CLI
func (c *CLI) Run() error {
	videoID := flag.String("video", "", "YouTube video ID")
	languages := flag.String("languages", "en", "Language codes separated by commas (e.g., 'en,ru,fr')")
	format := flag.String("format", "text", "Output format (text, json, srt)")
	flag.Parse()

	if *videoID == "" {
		return fmt.Errorf("video ID must be specified using the -video flag")
	}

	languageCodes := strings.Split(*languages, ",")
	transcript, err := c.api.GetTranscript(*videoID, languageCodes)
	if err != nil {
		return fmt.Errorf("error getting transcript: %w", err)
	}

	var formatter formatters.Formatter
	switch *format {
	case "json":
		formatter = &formatters.JSONFormatter{}
	case "srt":
		formatter = &formatters.SRTFormatter{}
	default:
		formatter = &formatters.TextFormatter{}
	}

	output, err := formatter.Format(transcript)
	if err != nil {
		return fmt.Errorf("error formatting transcript: %w", err)
	}

	fmt.Println(output)
	return nil
}

// Execute runs the CLI and handles errors
func Execute() {
	cli := NewCLI()
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
