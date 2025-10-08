package ctx

const (
	RequestID = "request_id"

	TopRoute         = "top_route"           // like /openai
	SubPath          = "sub_path"            // like /v1/chat/completions
	SubPathWithQuery = "sub_path_with_query" // like /v1/chat/completions?key=value
	TargetEndpoint   = "target_endpoint"     // like https://api.openai.com
	TargetURL        = "target_url"          // like https://api.openai.com/v1/chat/completions
)
