package parser

import (
	"os/exec"
	"strings"
)

type PDFParser struct{}

func (p *PDFParser) Parse(filePath string) (string, error) {
	// 使用 pdftotext 命令行工具提取 PDF 文本
	// -layout 保持布局，-enc UTF-8 使用 UTF-8 编码
	cmd := exec.Command("pdftotext", "-layout", "-enc", "UTF-8", filePath, "-")
	output, err := cmd.Output()
	if err != nil {
		// 如果 pdftotext 不可用，尝试使用 mutool
		cmd = exec.Command("mutool", "draw", "-F", "text", filePath)
		output, err = cmd.Output()
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(string(output)), nil
}
