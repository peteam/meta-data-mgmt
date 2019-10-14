package response

// ResponseResponse output json
type ResponseResponse struct {
	Header *Header      `json:"header"`
	Data   *interface{} `json:"data,omitempty"`
}
