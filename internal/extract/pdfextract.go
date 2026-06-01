package extract

import (
	"RagApp/internal/logging"
	"os"

	"github.com/ledongthuc/pdf"
)

// getTextFromPdf is a function to extract text from a pdf file
func getTextFromPdf(filePath string) (string, error) {
	logging.Trace("getTextFromPdf")

	//open the file
	r, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	r.Close()

	//get the file size
	fileInfo, err := r.Stat()
	if err != nil {
		return "", err
	}

	//Create the pdf reader for the file
	readCloser, err := pdf.NewReader(r, fileInfo.Size())
	if err != nil {
		return "", err
	}

	var text string

	//iterate through each page of the pdf
	for i := 1; i <= readCloser.NumPage(); i++ {
		//read the text content of the pdf
		page := readCloser.Page(i)
		if page.V.IsNull() {
			continue
		}

		//fonts arg is nil, I didn't understand what it was for anyway
		b, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}

		text += b
	}

	logging.Trace("returning getTextFromPdf")
	return text, nil
}
