# youtube-transcript-api-go

This is a Go version of the YouTube Transcript API library, which allows you to retrieve transcriptions (subtitles) for YouTube videos.
Based on https://github.com/jdepoix/youtube-transcript-api/.

## Installation

```
go get github.com/hightemp/youtube-transcript-api-go
```

## Usage

### As a library

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hightemp/youtube-transcript-api-go/api"
)

func main() {
	// Set proxy (if necessary)
	os.Setenv("HTTP_PROXY", "http://your-proxy-server:port")

	youtubeAPI := api.NewYouTubeTranscriptApi()

	videoID := "aBcDeFgHiJk" // Replace with the actual video ID
	languages := []string{"en", "ru"} // Language priority

	transcript, err := youtubeAPI.GetTranscript(videoID, languages)
	if err != nil {
		log.Fatalf("Error getting transcript: %v", err)
	}

	fmt.Printf("Transcript for video %s in language %s:\n", videoID, transcript.Language)
	for _, entry := range transcript.Entries {
		fmt.Printf("[%.2f - %.2f]: %s\n", entry.Start, entry.Start+entry.Duration, entry.Text)
	}
}
```

### Example of getting a transcript for a specific video

In the `examples/get_transcript` directory, there's an example demonstrating how to get a transcript for a video with ID "pPjDPe0duXc" in Russian. To run this example, execute the following commands:

```
cd examples/get_transcript
go run main.go
```

### As a CLI tool

```
HTTP_PROXY=http://your-proxy-server:port youtube-transcript-api-go -video=aBcDeFgHiJk -languages=en,ru -format=text
```

Available options:
- `-video`: YouTube video ID (required parameter)
- `-languages`: Language codes separated by commas (default "en")
- `-format`: Output format (text, json, srt) (default "text")

## License

This project is distributed under the MIT license. See the LICENSE file for details.