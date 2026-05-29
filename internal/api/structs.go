package api

// SearchParameters is a struct to handle the POST parameters sent to get the query.
type SearchParameters struct {
	Query string `json:"query"`
}

// SearchResponse is a struct to format the data for the response.
type SearchResponse struct {
	Result  string `json:"result"`
	Content string `json:"content"`
}
