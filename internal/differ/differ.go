// Package differ содержит функцию, которая создает дифнутые png и pdf для двух pdf
package differ

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	exampleDir string = "examples/"
	pathDir    string = ""
	pdfOut     string = filepath.Join("diff", "output-diff.pdf")
	pngPrefix  string = filepath.Join("diff", "output-diff")
)

func DiffPdf(file1 string, file2 string) error {
	os.MkdirAll("diff", 0o755)

	file1 = filepath.Join(exampleDir, file1)
	file2 = filepath.Join(pathDir, file2)

	tmp1 := filepath.Join("diff", "tmp1.pdf")
	tmp2 := filepath.Join("diff", "tmp2.pdf")

	cmd1 := exec.Command("qpdf", file1, "--pages", ".", "1", "--", tmp1)
	cmd2 := exec.Command("qpdf", file2, "--pages", ".", "1", "--", tmp2)

	err := cmd1.Run()
	if err != nil {
		return fmt.Errorf("ошибка выполнения qpdf: %w", err)
	}
	defer os.Remove(tmp1)

	err = cmd2.Run()
	if err != nil {
		return fmt.Errorf("ошибка выполнения qpdf: %w", err)
	}
	defer os.Remove(tmp2)

	cmd3 := exec.Command("diff-pdf", "--output-diff="+pdfOut, tmp1, tmp2)
	err = cmd3.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				err = nil
			}
		}
	}
	if err != nil {
		return fmt.Errorf("ошибка выполнения diff-pdf: %w", err)
	}

	cmd4 := exec.Command("pdftoppm", "-png", "-f", "1", "-l", "1", "-r", "150", pdfOut, pngPrefix)
	err = cmd4.Run()
	if err != nil {
		return fmt.Errorf("ошибка выполнения pdftoppm: %w", err)
	}
	return nil
}
