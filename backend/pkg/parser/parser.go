package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

type DocumentParser interface {
	Parse(filePath string) (string, error)
}

func GetParser(filename string) (DocumentParser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".pdf":
		return &PDFParser{}, nil
	case ".docx":
		return &DOCXParser{}, nil
	case ".txt":
		return &TextParser{}, nil
	case ".md":
		return &TextParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func SplitIntoChunks(text string, chunkSize int, overlap int) []string {
	var chunks []string
	words := strings.Fields(text)
	
	if len(words) == 0 {
		return chunks
	}
	
	for i := 0; i < len(words); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)
		
		if end >= len(words) {
			break
		}
	}
	
	return chunks
}
