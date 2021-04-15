package goemail

import (
	"bytes"
	"html/template"
)

func MessageFromHtmlTemplate(templateFile string, templateValues map[string]interface{}) (string, error) {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	var message bytes.Buffer
	err = t.Execute(&message, templateValues)
	if err != nil {
		return "", err
	}

	return message.String(), nil
}
