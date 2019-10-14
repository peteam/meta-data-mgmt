package response

/*
Standard response body structure to adhere
Reference: https://quickplay.atlassian.net/wiki/spaces/IASD/pages/2392065/HTTP+API+Specification+Standard
Author: Jerald M
*/

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ContentResponse struct {
	Header *Header     `json:"header"`
	Data   interface{} `json:"data,omitempty"`
}

type Data struct {
	ItemId           string `json:"itemId,omitempty"`
	UpdatedTimestamp int64  `json:"updatedTimestamp,omitempty"`
	Count            *int   `json:"count,omitempty"`
	ResourceType     string `json:"resourceType,omitempty"`
	URN              string `json:"urn,omitempty"`
	Schema           string `json:"schema,omitempty"`
}

type ResourceData struct {
	ResItem []Data `json:"items"`
}

type Header struct {
	Source     string  `json:"source"`
	Code       string  `json:"code"`
	Message    string  `json:"message"`
	SystemTime int64   `json:"systemtime"`
	Errors     []Error `json:"errors,omitempty"`
	Start      *int    `json:"start,omitempty"`
	Rows       *int    `json:"rows,omitempty"`
	Count      *int    `json:"count,omitempty"`
}

type Response struct {
	Header       *Header       `json:"header"`
	Data         *Data         `json:"data,omitempty"`
	ResourceData *ResourceData `json:"items,omitempty"`
}

type MultiResponse struct {
	Header       *Header       `json:"header"`
	ResourceData *ResourceData `json:"data,omitempty"`
}
