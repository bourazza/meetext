package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
)

// Extractor extracts plain text from a PDF file using the system's pdftotext tool.
type Extractor struct {
	log            zerolog.Logger
	maxUploadBytes int64
}

func NewExtractor(log zerolog.Logger) *Extractor {
	return &Extractor{
		log:            log.With().Str("component", "pdf_extractor").Logger(),
		maxUploadBytes: 250 * 1024 * 1024,
	}
}

// Extract reads a PDF from the provided reader and returns its text.
func (e *Extractor) Extract(ctx context.Context, r io.Reader, fileName string) (string, error) {
	if !strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
		return "", fmt.Errorf("invalid file extension: %s (only .pdf supported)", fileName)
	}

	// Invoke pdftotext to read from stdin (-) and write to stdout (-)
	cmd := exec.CommandContext(ctx, "pdftotext", "-", "-")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = r
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract pdf text via pdftotext: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.String(), nil
}
