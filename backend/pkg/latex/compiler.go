package latex

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Compiler struct {
	outputDir string
}

func NewCompiler(outputDir string) *Compiler {
	return &Compiler{outputDir: outputDir}
}

func (c *Compiler) Compile(latexContent string, filename string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return "", err
	}

	// Create temporary directory for compilation
	tempDir, err := os.MkdirTemp(c.outputDir, "latex-*")
	if err != nil {
		return "", err
	}

	// Write LaTeX content to file
	texFile := filepath.Join(tempDir, "main.tex")
	if err := os.WriteFile(texFile, []byte(latexContent), 0644); err != nil {
		return "", err
	}

	// Run xelatex twice for proper compilation
	for i := 0; i < 2; i++ {
		cmd := exec.Command("xelatex", "-interaction=nonstopmode", "-output-directory="+tempDir, texFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("xelatex compilation failed: %s, output: %s", err, string(output))
		}
	}

	// Check if PDF was generated
	pdfFile := filepath.Join(tempDir, "main.pdf")
	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF file was not generated")
	}

	// Move PDF to final location
	finalPDF := filepath.Join(c.outputDir, filename)
	if err := os.Rename(pdfFile, finalPDF); err != nil {
		return "", err
	}

	// Clean up temp directory (optional, can keep for debugging)
	// os.RemoveAll(tempDir)

	return finalPDF, nil
}
