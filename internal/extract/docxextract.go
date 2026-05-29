package extract

import (
	"os"

	"github.com/fumiama/go-docx"
)

// getTextFromDocx fetches the content of a docx or similar file.
func getTextFromDocx(filePath string) (string, error) {
	//Open the file
	r, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	//Get the file size for the docx package
	fileInfo, err := r.Stat()
	if err != nil {
		return "", err
	}
	size := fileInfo.Size()

	//Parse the docx file (which is essentially a zip file with XML)
	doc, err := docx.Parse(r, size)
	if err != nil {
		return "", err
	}

	var text string

	//iterate through the XML items
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		//Retrieve the paragraphs from the docx file
		case *docx.Paragraph:
			{
				text += it.(*docx.Paragraph).String() + "\n"
			}
		}
	}

	return text, nil
}
