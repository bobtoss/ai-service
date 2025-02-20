package doc

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dslipak/pdf"
)

type WordXML struct {
	XMLName xml.Name `xml:"document"`
	Body    struct {
		Paragraphs []struct {
			Texts []struct {
				Text string `xml:",chardata"`
			} `xml:"r>t"`
		} `xml:"body>p"`
	} `xml:"body"`
}

func DecodeBase64ToString(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func DecodeBase64ToFileAndRead(encoded string) ([]string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// Определение MIME-типа
	mimeType := http.DetectContentType(decodedBytes)

	// Определение расширения файла
	extension := strings.TrimPrefix(mimeType, "application/")
	if extension == "" {
		extension = "bin"
	}
	switch extension {
	case "zip":
		extension = "docx"
	}
	filename := "temp." + extension

	// Записываем в файл
	err = os.WriteFile(filename, decodedBytes, 0644)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filename)

	// Читаем содержимое файла в зависимости от его типа
	switch mimeType {
	case "application/pdf":
		return readPDFFile(filename)
	case "application/zip":
		return readWordFile(filename)
	case "text/plain":
		return readTextFile(filename)
	default:
		return nil, errors.New("неподдерживаемый формат файла")
	}
}

func readTextFile(filename string) ([]string, error) {
	_, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func readPDFFile(filename string) ([]string, error) {
	r, err := pdf.Open(filename)
	if err != nil {
		return nil, err
	}
	var chunks []string
	for i := 0; i < r.NumPage(); i++ {
		page := r.Page(i + 1)
		if page.V.IsNull() {
			continue
		}
		s, err := page.GetPlainText(nil)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, s)
	}

	return chunks, nil
}

func readWordFile(filename string) ([]string, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var textContent string

	for _, file := range r.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			// Read XML content
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, rc)
			if err != nil {
				return nil, err
			}

			textContent = extractTextFromXML(buf.String())
			break
		}
	}

	if textContent == "" {
		return nil, fmt.Errorf("could not find 'word/document.xml' in DOCX file")
	}

	pages := strings.Split(textContent, "[PAGE_BREAK]")

	return pages, nil
}

func extractTextFromXML(xmlContent string) string {
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	var extractedText strings.Builder

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			// Detect page breaks
			if t.Name.Local == "br" {
				for _, attr := range t.Attr {
					if attr.Name.Local == "type" && attr.Value == "page" {
						extractedText.WriteString("[PAGE_BREAK]") // Mark page break
					}
				}
			}
		case xml.CharData:
			extractedText.WriteString(string(t))
		}
	}

	return extractedText.String()
}
