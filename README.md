# CardScanner

CardScanner is a robust Go-based tool for scanning local directories to detect sensitive card numbers in various file formats. It is designed for security audits, compliance checks, or data loss prevention tasks, and emphasizes reliability, efficiency, and clear reporting.

## Features

* Recursively scans a specified directory
* Supports multiple file types:

  * Plain text files (.txt, .csv)
  * Word documents (.docx)
  * Excel spreadsheets (.xlsx)
  * PDF files (via the external tool pdftotext)
* Extracts potential card numbers in common formats, including space or hyphen separators
* Validates detected card numbers using the Luhn algorithm to reduce false positives
* Identifies card type based on BIN or IIN prefixes
* Detects expiry date and CVV patterns when found in proximity
* Optionally hashes detected numbers using SHA-256 for added security
* Color-coded output for different risk levels
* Handles large files efficiently without excessive memory consumption
* Skips unreadable or inaccessible files gracefully

## How It Works

1. Recursively traverses directories starting from a user-specified path.
2. For each supported file, extracts readable text using format-appropriate libraries or tools.
3. Applies regular expressions to find candidate card numbers and related patterns.
4. Validates each candidate with the Luhn algorithm and checks known prefixes to classify card type.
5. Produces a clear report including file path, line number, detected type, and a hashed version of the number if configured.
6. Highlights findings using different colors to indicate confidence levels.

## Requirements

* Go version 1.20 or higher
* The command-line tool pdftotext must be installed and available in the system PATH (required for extracting text from PDF files)

## Installation

Clone the repository and download the required dependencies using the Go module system.

```bash
git clone <repository_url>
cd cardscanner
go mod tidy
```

## Build and Run

Build the executable:

```bash
go build -o cardscanner
```

Run the scanner on a target directory:

```bash
./cardscanner /path/to/target/directory
```

Example:

```bash
./cardscanner ~/Documents/Reports
```

## Example Matches

* Valid Visa card: 4111-1111-1111-1111
* Valid MasterCard: 5115 1051 0510 5100
* Invalid card (wrong checksum): 1234567890123456
* Invalid card (too short): 1234 5678 9012 345

## Output

The tool outputs each match with:

* Full file path
* Line number
* Matched text
* Card type if recognized (for example, Visa, MasterCard)
* Hash of the detected number to avoid storing plain card numbers
* Confidence level indicated by color: red for valid matches, yellow for possible but invalid patterns

## Best Practices

* Always run the scanner on directories you have permission to access.
* Ensure that pdftotext is installed to handle PDF files reliably.
* Review results carefully to filter out potential false positives due to loosely matching numeric patterns.

## Limitations

* Optical character recognition (OCR) for scanned images inside PDFs is not included. For scanned PDFs, ensure pdftotext can extract text correctly.
* Some uncommon card formats or region-specific cards may not be fully supported.
* False positives are still possible due to numeric patterns that resemble card numbers but are unrelated.

## License

This project is intended as a demonstration of secure local scanning practices. It does not store or transmit any detected card data. Use responsibly and ensure compliance with local regulations when scanning sensitive directories.
