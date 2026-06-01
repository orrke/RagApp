package extract

import (
	"log"
	"path/filepath"
)

var extractors = map[string]func(string) (string, error){
	".txt":  getTextFromPlainText,
	".csv":  getTextFromPlainText,
	".bat":  getTextFromPlainText,
	".doc":  getTextFromDocx,
	".docx": getTextFromDocx,
	".xlsx": getTextFromDocx,
	".pdf":  getTextFromPdf,
	".odt":  getTextFromOdt,
}

// Extract :
//
// Universal text extractor from a bunch of supported files.
// If the file isn't supported it will return an error.
// The map in this file is how we get the content of the plain text file
func Extract(filePath string) string {
	extension := filepath.Ext(filePath)

	if extension == "" {
		log.Println("Extract: no extension")
		return ""
	}

	extractor, ok := extractors[extension]
	if !ok {
		log.Println("Unsupported file extension:", extension)
		return ""
	}
	result, err := extractor(filePath)

	if err != nil {
		log.Println(err)
	}

	return result
}
