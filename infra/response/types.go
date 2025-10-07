package response

type ErrorDetail struct {
	RequestID string `json:"request_id,omitempty"`
	Note      string `json:"note,omitempty"`
}

type ErrorInfo struct {
	Message string       `json:"message"`
	Type    string       `json:"type"`
	Source  string       `json:"source"`
	Details *ErrorDetail `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error ErrorInfo `json:"error"`
}
