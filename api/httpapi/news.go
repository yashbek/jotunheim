package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	URL         string `json:"url"`
	Date        string `json:"date"`
	Source      string `json:"source"`
}

type SerpApiResponse struct {
	NewsResults []struct {
		Title     string `json:"title"`
		Link      string `json:"link"`
		Snippet   string `json:"snippet"`
		Source    string `json:"source"`
		Date      string `json:"date"`
		Thumbnail string `json:"thumbnail"`
	} `json:"news_results"`
}

func HandleNews(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := "00fd6278c7f1d2e0aa4560643e9ad80dc7762f5db682d7e422ff2d40ee425eec"

	baseURL := "https://serpapi.com/search.json"
	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("engine", "google")
	params.Add("q", "hnefatafl OR tafl game OR viking board game")
	params.Add("tbm", "nws")
	params.Add("num", "10")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := client.Get(fullURL)
	if err != nil {
		http.Error(w, "Failed to fetch news", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("SerpApi error: %s", string(body)), http.StatusInternalServerError)
		return
	}

	var serpResponse SerpApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&serpResponse); err != nil {
		http.Error(w, "Failed to parse news data", http.StatusInternalServerError)
		return
	}

	news := make([]NewsArticle, 0)
	for _, result := range serpResponse.NewsResults {
		if result.Thumbnail == "" {
			continue
		}

		news = append(news, NewsArticle{
			Title:       result.Title,
			Description: result.Snippet,
			ImageURL:    result.Thumbnail,
			URL:         result.Link,
			Date:        result.Date,
			Source:      result.Source,
		})
	}

	if err := json.NewEncoder(w).Encode(news); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func HandleMockNews(w http.ResponseWriter, r *http.Request) {
	mockNews := []NewsArticle{
		{
			Title:       "Ancient Board Game Championships Features Hnefatafl",
			Description: "The annual Viking Games Festival showcases traditional Norse board games, with Hnefatafl taking center stage in this year's competition.",
			ImageURL:    "/api/placeholder/400/300",
			URL:         "https://example.com/article1",
			Date:        "2 days ago",
			Source:      "Viking Heritage Magazine",
		},
		{
			Title:       "New Archaeological Find Reveals Complete Hnefatafl Set",
			Description: "Archaeologists in Northern Scotland have uncovered a perfectly preserved Hnefatafl board game set, dating back to the 9th century.",
			ImageURL:    "/api/placeholder/400/300",
			URL:         "https://example.com/article2",
			Date:        "5 days ago",
			Source:      "Archaeological Digest",
		},
	}

	json.NewEncoder(w).Encode(mockNews)
}
