package goemail

import (
	"bytes"
	"html/template"
)

// MessageFromHtmlTemplate parses html template file
// and fill the given values in the variables accordingly
func MessageFromHtmlTemplate(templateFilePath string, values interface{}) (string, error) {
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	var message bytes.Buffer
	err = t.Execute(&message, values)
	if err != nil {
		return "", err
	}

	return message.String(), nil
}
