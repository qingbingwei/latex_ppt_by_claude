package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

type DOCXParser struct{}

func (p *DOCXParser) Parse(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var content string
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}

			content = extractTextFromXML(data)
			break
		}
	}

	return content, nil
}

func extractTextFromXML(data []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var text strings.Builder

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "t" {
				var content string
				if err := decoder.DecodeElement(&content, &element); err == nil {
					text.WriteString(content)
				}
			}
		}
	}

	return text.String()
}
