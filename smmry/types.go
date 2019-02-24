package smmry

//Summary represents a smmry API response
type Summary struct {
	Message  string   `json:"sm_api_message"`
	Title    string   `json:"sm_api_title"`
	Content  string   `json:"sm_api_content"`
	Keywords []string `json:"sm_api_keyword_array"`
	Error    int      `json:"sm_api_error"`
}

// Params represents a smmry request, which will be mapped into URL params
type Params struct {
	URL    string
	Length int
}

//Client represents the smmry client, initialized with the API token
type Client struct {
	Token string
}
