package pdf

import "strings"

const (
	targetWordsPerChunk = 2000
	maxWordsPerChunk    = 3000
)

// ChunkText splits a large text into word-safe chunks for LLM processing.
func ChunkText(text string) []string {
	paragraphs := strings.Split(text, "\n\n")

	var chunks []string
	var current strings.Builder
	wordCount := 0

	flush := func() {
		s := strings.TrimSpace(current.String())
		if s != "" {
			chunks = append(chunks, s)
		}
		current.Reset()
		wordCount = 0
	}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		words := strings.Fields(para)
		if wordCount+len(words) > maxWordsPerChunk && wordCount > 0 {
			flush()
		}
		if current.Len() > 0 {
			current.WriteString("\n\n")
		}
		current.WriteString(para)
		wordCount += len(words)
		if wordCount >= targetWordsPerChunk {
			flush()
		}
	}
	flush()

	if len(chunks) == 0 && strings.TrimSpace(text) != "" {
		chunks = append(chunks, strings.TrimSpace(text))
	}
	return chunks
}
