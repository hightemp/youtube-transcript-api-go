package errors

import (
	"fmt"
	"strings"
)

type TranscriptError struct {
	VideoID string
	Message string
}

func (e *TranscriptError) Error() string {
	return fmt.Sprintf("Transcript error for video %s: %s", e.VideoID, e.Message)
}

func NewYouTubeRequestFailed(err error, videoID string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: fmt.Sprintf("YouTube request error: %v", err),
	}
}

func NewTooManyRequests(videoID string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: "Too many requests",
	}
}

func NewVideoUnavailable(videoID string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: "Video unavailable",
	}
}

func NewTranscriptsDisabled(videoID string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: "Transcripts are disabled for this video",
	}
}

func NewNoTranscriptAvailable(videoID string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: "No transcripts available for this video",
	}
}

func NewFailedToCreateConsentCookie(videoID string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: "Failed to create consent cookie",
	}
}

func NewNoTranscriptFound(videoID string, languageCodes []string) error {
	return &TranscriptError{
		VideoID: videoID,
		Message: fmt.Sprintf("No transcript found for languages: %s", strings.Join(languageCodes, ", ")),
	}
}
