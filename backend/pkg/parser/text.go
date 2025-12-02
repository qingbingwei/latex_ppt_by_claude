package parser

import (
	"os"
)

type TextParser struct{}

func (p *TextParser) Parse(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
