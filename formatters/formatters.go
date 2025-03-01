package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hightemp/youtube-transcript-api-go/api"
)

type Formatter interface {
	Format(transcript *api.Transcript) (string, error)
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(transcript *api.Transcript) (string, error) {
	jsonData, err := json.MarshalIndent(transcript.Entries, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %w", err)
	}
	return string(jsonData), nil
}

type TextFormatter struct{}

func (f *TextFormatter) Format(transcript *api.Transcript) (string, error) {
	var sb strings.Builder
	for _, entry := range transcript.Entries {
		sb.WriteString(fmt.Sprintf("[%.2f - %.2f]: %s\n", entry.Start, entry.Start+entry.Duration, entry.Text))
	}
	return sb.String(), nil
}

type SRTFormatter struct{}

func (f *SRTFormatter) Format(transcript *api.Transcript) (string, error) {
	var sb strings.Builder
	for i, entry := range transcript.Entries {
		startTime := formatSRTTime(entry.Start)
		endTime := formatSRTTime(entry.Start + entry.Duration)
		sb.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, entry.Text))
	}
	return sb.String(), nil
}

func formatSRTTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	milliseconds := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, milliseconds)
}
