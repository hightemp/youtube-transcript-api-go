package main

import (
	"fmt"
	"log"

	"github.com/hightemp/youtube-transcript-api-go/api"
)

func main() {
	youtubeAPI := api.NewYouTubeTranscriptApi()

	videoID := "pPjDPe0duXc"
	languages := []string{"en", "ru"}

	transcript, err := youtubeAPI.GetTranscript(videoID, languages)
	if err != nil {
		log.Fatalf("Error getting transcript: %v", err)
	}

	fmt.Printf("Transcript for video %s in language %s:\n", videoID, transcript.Language)
	for _, entry := range transcript.Entries {
		fmt.Printf("[%.2f - %.2f]: %s\n", entry.Start, entry.Start+entry.Duration, entry.Text)
	}
}
