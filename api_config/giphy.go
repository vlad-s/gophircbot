package api_config

type Giphy struct {
	ApiKey string `json:"api_key"`
	Limit  uint8  `json:"limit"`
}

type GiphyResponse struct {
	Data []struct {
		ID       string `json:"id"`
		Slug     string `json:"slug"`
		URL      string `json:"url"`
		ShortURL string `json:"bitly_url"`
		Rating   string `json:"rating"`
	} `json:"data"`
	Pagination struct {
		Total  int `json:"total_count"`
		Count  int `json:"count"`
		Offset int `json:"offset"`
	} `json:"pagination"`
	Meta struct {
		Status     int    `json:"status"`
		Message    string `json:"msg"`
		ResponseID string `json:"response_id"`
	} `json:"meta"`
}
