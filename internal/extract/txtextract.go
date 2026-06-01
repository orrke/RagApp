package extract

import (
	"RagApp/internal/logging"
	"bufio"
	"os"
)

//We don't need an external library for this one it's really easy in go

// getTextFromPlainText is a function to get the contents of any plain text file
// (which can be a txt, bat, csv... so on)
func getTextFromPlainText(filePath string) (string, error) {
	logging.Trace("getTextFromPlainText")

	//Open the reader
	r, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var text string

	//Create the file scanner
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		//add the text line to the content.
		text += scanner.Text()
	}

	logging.Trace("returning getTextFromPlainText")
	return text, nil
}
