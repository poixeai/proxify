package util

import "strings"

// ExtractRoute splits "/openai/v1/chat" → "openai", "v1/chat"
func ExtractRoute(path string) (string, string) {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	top := parts[0]
	if len(parts) == 2 {
		return top, "/" + parts[1]
	}
	return top, ""
}
