package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/unidoc/unioffice/document"
	"github.com/xuri/excelize/v2"
)

// Regex for card-like patterns
var cardPattern = regexp.MustCompile(`(?:\b(?:\d[ -]*?){13,19}\b)`)

// Regex for expiration date
var expDatePattern = regexp.MustCompile(`\b(0[1-9]|1[0-2])[\/\-](\d{2}|\d{4})\b`)

// CVV pattern
var cvvPattern = regexp.MustCompile(`\b\d{3,4}\b`)

// Card type prefixes
var visaPrefix = regexp.MustCompile(`^4`)
var mcPrefix = regexp.MustCompile(`^5[1-5]`)

// Luhn validation
func isValidLuhn(number string) bool {
	sum := 0
	alt := false
	for i := len(number) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(number[i]))
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

func detectType(number string) string {
	switch {
	case visaPrefix.MatchString(number):
		return "Visa"
	case mcPrefix.MatchString(number):
		return "MasterCard"
	default:
		return "Unknown"
	}
}

func scanText(path string, text string) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		matches := cardPattern.FindAllString(line, -1)
		for _, match := range matches {
			clean := strings.ReplaceAll(strings.ReplaceAll(match, " ", ""), "-", "")
			if isValidLuhn(clean) {
				cardType := detectType(clean)
				fmt.Printf("\033[31m[HIGH]\033[0m %s:%d | Card: %s | Type: %s | Hash: %x\n",
					path, i+1, match, cardType, sha256.Sum256([]byte(clean)))
			} else {
				fmt.Printf("\033[33m[LOW]\033[0m %s:%d | Possibly Invalid: %s\n", path, i+1, match)
			}
		}
	}
}

func processTextFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening %s: %v", path, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	scanText(path, strings.Join(lines, "\n"))
}

func processDocx(path string) {
	doc, err := document.Open(path)
	if err != nil {
		log.Printf("Error opening DOCX: %v", err)
		return
	}
	var text bytes.Buffer
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text.WriteString(run.Text())
		}
		text.WriteString("\n")
	}
	scanText(path, text.String())
}

func processXlsx(path string) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Printf("Error opening XLSX: %v", err)
		return
	}
	var text bytes.Buffer
	for _, name := range f.GetSheetList() {
		rows, err := f.GetRows(name)
		if err != nil {
			continue
		}
		for _, row := range rows {
			text.WriteString(strings.Join(row, " ") + "\n")
		}
	}
	scanText(path, text.String())
}

func processPdf(path string) {
	cmd := exec.Command("pdftotext", path, "-")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error processing PDF: %v", err)
		return
	}
	scanText(path, string(out))
}

func scanDir(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Cannot access %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".txt", ".csv":
			processTextFile(path)
		case ".docx":
			processDocx(path)
		case ".xlsx":
			processXlsx(path)
		case ".pdf":
			processPdf(path)
		default:
			// skip others
		}

		return nil
	})
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cardscanner <directory>")
		os.Exit(1)
	}
	dir := os.Args[1]
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		fmt.Println("Invalid directory:", dir)
		os.Exit(1)
	}
	scanDir(dir)
}
