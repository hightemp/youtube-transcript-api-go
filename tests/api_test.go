package tests

import (
	"testing"

	"github.com/hightemp/youtube-transcript-api-go/api"
)

func TestGetTranscript(t *testing.T) {
	youtubeAPI := api.NewYouTubeTranscriptApi()

	testCases := []struct {
		name      string
		videoID   string
		languages []string
		wantErr   bool
	}{
		{
			name:      "Valid video ID",
			videoID:   "pPjDPe0duXc",
			languages: []string{"ru"},
			wantErr:   false,
		},
		{
			name:      "Invalid video ID",
			videoID:   "invalid_id",
			languages: []string{"en"},
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transcript, err := youtubeAPI.GetTranscript(tc.videoID, tc.languages)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetTranscript() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if transcript == nil {
					t.Errorf("GetTranscript() returned nil transcript for valid video ID")
					return
				}

				if transcript.VideoID != tc.videoID {
					t.Errorf("GetTranscript() returned transcript with incorrect VideoID, got %s, want %s", transcript.VideoID, tc.videoID)
				}

				if len(transcript.Entries) == 0 {
					t.Errorf("GetTranscript() returned empty transcript entries")
				}

				if transcript.Language == "" {
					t.Errorf("GetTranscript() returned empty Language")
				}

				if transcript.LanguageCode == "" {
					t.Errorf("GetTranscript() returned empty LanguageCode")
				}

				t.Logf("First entries of the transcript for video %s:", tc.videoID)
				for i, entry := range transcript.Entries {
					if i >= 5 {
						break
					}
					t.Logf("[%.2f - %.2f]: %s", entry.Start, entry.Start+entry.Duration, entry.Text)
				}
			}
		})
	}
}

func TestListTranscripts(t *testing.T) {
	youtubeAPI := api.NewYouTubeTranscriptApi()

	testCases := []struct {
		name    string
		videoID string
		wantErr bool
	}{
		{
			name:    "Valid video ID",
			videoID: "pPjDPe0duXc",
			wantErr: false,
		},
		{
			name:    "Invalid video ID",
			videoID: "invalid_id",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transcriptList, err := youtubeAPI.ListTranscripts(tc.videoID)

			if (err != nil) != tc.wantErr {
				t.Errorf("ListTranscripts() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if transcriptList == nil {
					t.Errorf("ListTranscripts() returned nil TranscriptList for valid video ID")
					return
				}

				if transcriptList.VideoID != tc.videoID {
					t.Errorf("ListTranscripts() returned TranscriptList with incorrect VideoID, got %s, want %s", transcriptList.VideoID, tc.videoID)
				}

				if len(transcriptList.ManuallyCreatedTranscripts) == 0 && len(transcriptList.GeneratedTranscripts) == 0 {
					t.Errorf("ListTranscripts() returned empty transcripts")
				}
			}
		})
	}
}
