package doc

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"archive/zip"
	"github.com/ledongthuc/pdf"
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

func TextToChunks(input string) []string {
	return make([]string, 0)
}

func DecodeBase64ToString(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func DecodeBase64ToFileAndRead(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	// Определение MIME-типа
	mimeType := http.DetectContentType(decodedBytes)

	// Определение расширения файла
	extension := strings.TrimPrefix(mime.TypeByExtension(mimeType), "/")
	if extension == "" {
		extension = "bin"
	}
	filename := "temp." + extension

	// Записываем в файл
	err = os.WriteFile(filename, decodedBytes, 0644)
	if err != nil {
		return "", err
	}

	// Читаем содержимое файла в зависимости от его типа
	switch mimeType {
	case "application/pdf":
		return readPDFFile(filename)
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return readWordFile(filename)
	case "text/plain":
		return readTextFile(filename)
	default:
		return "", errors.New("неподдерживаемый формат файла")
	}
}

func readTextFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func readPDFFile(filename string) (string, error) {
	f, r, err := pdf.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var text string
	for i := 0; i < r.NumPage(); i++ {
		page := r.Page(i)
		text += page.Content().Text[0].Font
	}
	return text, nil
}

func readWordFile(filename string) (string, error) {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Find and extract word/document.xml
	var xmlData []byte
	for _, file := range reader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			xmlData, err = io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if len(xmlData) == 0 {
		return "", fmt.Errorf("word/document.xml not found in DOCX")
	}

	// Parse XML content
	var wordDoc WordXML
	err = xml.Unmarshal(xmlData, &wordDoc)
	if err != nil {
		return "", err
	}

	// Extract text from paragraphs
	var textContent string
	for _, p := range wordDoc.Body.Paragraphs {
		for _, t := range p.Texts {
			textContent += t.Text + " "
		}
		textContent += "\n"
	}

	return textContent, nil
}
