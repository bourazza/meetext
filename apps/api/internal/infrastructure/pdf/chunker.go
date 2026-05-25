package pdf

import (
	"strings"
)

const (
	// Optimized for llama3.1:8b with 8k-16k context
	targetWordsPerChunk = 2000
	maxWordsPerChunk    = 2500
	overlapWords        = 250
)

type Chunk struct {
	Index   int
	Content string
	Words   int
}

// ChunkText splits text into semantic chunks with overlap for context preservation.
func ChunkText(text string) []Chunk {
	paragraphs := strings.Split(text, "\n\n")
	
	var chunks []Chunk
	var currentChunk strings.Builder
	currentWords := 0
	chunkIndex := 0

	// Track last N words for overlap
	var overlapBuffer []string

	flush := func() {
		content := strings.TrimSpace(currentChunk.String())
		if content == "" {
			return
		}
		
		chunks = append(chunks, Chunk{
			Index:   chunkIndex,
			Content: content,
			Words:   currentWords,
		})
		chunkIndex++

		// Prepare overlap for next chunk
		words := strings.Fields(content)
		if len(words) > overlapWords {
			overlapBuffer = words[len(words)-overlapWords:]
		} else {
			overlapBuffer = words
		}

		currentChunk.Reset()
		currentWords = 0

		// Add overlap to start of next chunk
		if len(overlapBuffer) > 0 {
			currentChunk.WriteString(strings.Join(overlapBuffer, " "))
			currentChunk.WriteString("\n\n")
			currentWords = len(overlapBuffer)
		}
	}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		words := strings.Fields(para)
		paraWordCount := len(words)

		// If adding this paragraph exceeds max, flush current chunk
		if currentWords > 0 && currentWords+paraWordCount > maxWordsPerChunk {
			flush()
		}

		// Add paragraph to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(para)
		currentWords += paraWordCount

		// If we've reached target size, flush
		if currentWords >= targetWordsPerChunk {
			flush()
		}
	}

	// Flush remaining content
	flush()

	// If no chunks were created but we have text, create one chunk
	if len(chunks) == 0 && strings.TrimSpace(text) != "" {
		chunks = append(chunks, Chunk{
			Index:   0,
			Content: strings.TrimSpace(text),
			Words:   len(strings.Fields(text)),
		})
	}

	return chunks
}
