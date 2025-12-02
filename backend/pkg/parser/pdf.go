package parser

import (
	"bytes"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type PDFParser struct{}

func (p *PDFParser) Parse(filePath string) (string, error) {
	f, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	numPages, err := f.GetNumPages()
	if err != nil {
		return "", err
	}

	for i := 1; i <= numPages; i++ {
		page, err := f.GetPage(i)
		if err != nil {
			continue
		}

		ex, err := extractor.New(page)
		if err != nil {
			continue
		}

		text, err := ex.ExtractText()
		if err != nil {
			continue
		}

		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return buf.String(), nil
}
