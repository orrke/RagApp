package extract

import (
	"RagApp/internal/logging"

	"sbinet.org/x/odf"
)

// getTextFromOdt is a function to retrieve the contents of an odt file
func getTextFromOdt(filePath string) (string, error) {
	logging.Trace("getTextFromOdt")

	//open the file directly with the odf library
	doc, err := odf.Open(filePath)

	if err != nil {
		return "", err
	}
	defer doc.Close()

	var text string

	//Iterate through the elements of the odt file
	err = odf.Walk(doc.Node(), func(n odf.Node, _ bool) (odf.WalkStatus, error) {
		switch n := n.(type) {
		//Retrieve the Text from the odt file
		case *odf.Text:
			if n.Value != "" {
				text += n.Value + "\n"
			}
		}

		return odf.WalkContinue, nil
	})

	if err != nil {
		return "", err
	}

	logging.Trace("returning getTextFromOdt")
	return text, nil
}
