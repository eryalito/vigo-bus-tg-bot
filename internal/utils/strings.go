package utils

import "strings"

// SplitLongMessage splits a long message into an array of strings, each of which is a segment
// of the original string split at breaklines and ensuring that no segment exceeds 4000 characters.
func SplitLongMessage(sendText string) []string {
	const maxMessageLength = 4000

	lines := strings.Split(sendText, "\n")
	var sb strings.Builder
	var messages []string

	for _, line := range lines {
		if sb.Len()+len(line)+1 > maxMessageLength {
			messages = append(messages, sb.String())
			sb.Reset()
		}
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
	}

	if sb.Len() > 0 {
		messages = append(messages, sb.String())
	}

	return messages
}
