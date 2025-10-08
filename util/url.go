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

// JoinURL joins base and sub path, ensuring there is only one "/" between them
func JoinURL(base, sub string) string {
	base = strings.TrimRight(base, "/")
	sub = strings.TrimLeft(sub, "/")
	return base + "/" + sub
}
