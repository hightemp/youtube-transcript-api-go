package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hightemp/youtube-transcript-api-go/errors"
)

const (
	WATCH_URL = "https://www.youtube.com/watch?v=%s"
)

type YouTubeTranscriptApi struct {
	httpClient *http.Client
}

type Transcript struct {
	VideoID              string
	Language             string
	LanguageCode         string
	IsGenerated          bool
	TranslationLanguages []TranslationLanguage
	Entries              []TranscriptEntry
	url                  string
}

type TranscriptEntry struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	Duration float64 `json:"duration"`
}

type TranslationLanguage struct {
	Language     string `json:"language"`
	LanguageCode string `json:"language_code"`
}

type TranscriptList struct {
	VideoID                    string
	ManuallyCreatedTranscripts map[string]*Transcript
	GeneratedTranscripts       map[string]*Transcript
	TranslationLanguages       []TranslationLanguage
}

func NewYouTubeTranscriptApi() *YouTubeTranscriptApi {
	httpClient := &http.Client{}

	if proxyURL, exists := os.LookupEnv("HTTP_PROXY"); exists {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err == nil {
			httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURLParsed),
			}
		}
	}

	return &YouTubeTranscriptApi{
		httpClient: httpClient,
	}
}

func (api *YouTubeTranscriptApi) ListTranscripts(videoID string) (*TranscriptList, error) {
	html, err := api.fetchVideoHTML(videoID)
	if err != nil {
		return nil, err
	}

	captionsJSON, err := api.extractCaptionsJSON(html, videoID)
	if err != nil {
		return nil, err
	}

	return buildTranscriptList(api.httpClient, videoID, captionsJSON)
}

func (api *YouTubeTranscriptApi) GetTranscript(videoID string, languages []string) (*Transcript, error) {
	transcriptList, err := api.ListTranscripts(videoID)
	if err != nil {
		return nil, err
	}

	transcript, err := transcriptList.FindTranscript(languages)
	if err != nil {
		return nil, err
	}

	err = transcript.Fetch(api.httpClient)
	if err != nil {
		return nil, err
	}

	return transcript, nil
}

func (api *YouTubeTranscriptApi) fetchVideoHTML(videoID string) (string, error) {
	resp, err := api.httpClient.Get(fmt.Sprintf(WATCH_URL, videoID))
	if err != nil {
		return "", errors.NewYouTubeRequestFailed(err, videoID)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.NewYouTubeRequestFailed(err, videoID)
	}

	html := string(body)

	if strings.Contains(html, `action="https://consent.youtube.com/s"`) {
		err = api.createConsentCookie(html, videoID)
		if err != nil {
			return "", err
		}
		return api.fetchVideoHTML(videoID)
	}

	return html, nil
}

func (api *YouTubeTranscriptApi) extractCaptionsJSON(html, videoID string) (map[string]interface{}, error) {
	parts := strings.Split(html, `"captions":`)
	if len(parts) <= 1 {
		if strings.Contains(html, `class="g-recaptcha"`) {
			return nil, errors.NewTooManyRequests(videoID)
		}
		if !strings.Contains(html, `"playabilityStatus":`) {
			return nil, errors.NewVideoUnavailable(videoID)
		}
		return nil, errors.NewTranscriptsDisabled(videoID)
	}

	jsonPart := strings.Split(parts[1], `,"videoDetails`)[0]
	var captionsJSON map[string]interface{}
	err := json.Unmarshal([]byte(jsonPart), &captionsJSON)
	if err != nil {
		return nil, errors.NewYouTubeRequestFailed(err, videoID)
	}

	playerCaptionsTracklistRenderer, ok := captionsJSON["playerCaptionsTracklistRenderer"].(map[string]interface{})
	if !ok {
		return nil, errors.NewTranscriptsDisabled(videoID)
	}

	if _, ok := playerCaptionsTracklistRenderer["captionTracks"]; !ok {
		return nil, errors.NewNoTranscriptAvailable(videoID)
	}

	return playerCaptionsTracklistRenderer, nil
}

func (api *YouTubeTranscriptApi) createConsentCookie(html, videoID string) error {
	re := regexp.MustCompile(`name="v" value="(.*?)"`)
	match := re.FindStringSubmatch(html)
	if match == nil {
		return errors.NewFailedToCreateConsentCookie(videoID)
	}

	youtubeURL, _ := url.Parse("https://www.youtube.com")
	cookies := []*http.Cookie{
		{
			Name:   "CONSENT",
			Value:  "YES+" + match[1],
			Domain: ".youtube.com",
		},
	}
	api.httpClient.Jar.SetCookies(youtubeURL, cookies)

	return nil
}

func buildTranscriptList(httpClient *http.Client, videoID string, captionsJSON map[string]interface{}) (*TranscriptList, error) {
	translationLanguages := []TranslationLanguage{}
	if translationLanguagesJSON, ok := captionsJSON["translationLanguages"].([]interface{}); ok {
		for _, lang := range translationLanguagesJSON {
			if langMap, ok := lang.(map[string]interface{}); ok {
				translationLanguages = append(translationLanguages, TranslationLanguage{
					Language:     langMap["languageName"].(map[string]interface{})["simpleText"].(string),
					LanguageCode: langMap["languageCode"].(string),
				})
			}
		}
	}

	manuallyCreatedTranscripts := make(map[string]*Transcript)
	generatedTranscripts := make(map[string]*Transcript)

	captionTracks := captionsJSON["captionTracks"].([]interface{})
	for _, captionTrack := range captionTracks {
		caption := captionTrack.(map[string]interface{})
		isGenerated := caption["kind"] == "asr"
		transcript := &Transcript{
			VideoID:              videoID,
			Language:             caption["name"].(map[string]interface{})["simpleText"].(string),
			LanguageCode:         caption["languageCode"].(string),
			IsGenerated:          isGenerated,
			url:                  caption["baseUrl"].(string),
			TranslationLanguages: translationLanguages,
		}

		if isGenerated {
			generatedTranscripts[caption["languageCode"].(string)] = transcript
		} else {
			manuallyCreatedTranscripts[caption["languageCode"].(string)] = transcript
		}
	}

	return &TranscriptList{
		VideoID:                    videoID,
		ManuallyCreatedTranscripts: manuallyCreatedTranscripts,
		GeneratedTranscripts:       generatedTranscripts,
		TranslationLanguages:       translationLanguages,
	}, nil
}

func (tl *TranscriptList) FindTranscript(languageCodes []string) (*Transcript, error) {
	for _, langCode := range languageCodes {
		if transcript, ok := tl.ManuallyCreatedTranscripts[langCode]; ok {
			return transcript, nil
		}
		if transcript, ok := tl.GeneratedTranscripts[langCode]; ok {
			return transcript, nil
		}
	}
	return nil, errors.NewNoTranscriptFound(tl.VideoID, languageCodes)
}

func (t *Transcript) Fetch(httpClient *http.Client) error {
	resp, err := httpClient.Get(t.url)
	if err != nil {
		return errors.NewYouTubeRequestFailed(err, t.VideoID)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.NewYouTubeRequestFailed(err, t.VideoID)
	}

	var xmlData struct {
		Text []struct {
			Start string `xml:"start,attr"`
			Dur   string `xml:"dur,attr"`
			Value string `xml:",chardata"`
		} `xml:"text"`
	}

	err = xml.Unmarshal(body, &xmlData)
	if err != nil {
		return errors.NewYouTubeRequestFailed(err, t.VideoID)
	}

	t.Entries = make([]TranscriptEntry, len(xmlData.Text))
	for i, text := range xmlData.Text {
		start, _ := strconv.ParseFloat(text.Start, 64)
		duration, _ := strconv.ParseFloat(text.Dur, 64)
		t.Entries[i] = TranscriptEntry{
			Text:     text.Value,
			Start:    start,
			Duration: duration,
		}
	}

	return nil
}
