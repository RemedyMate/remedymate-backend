package mail

import (
	"bytes"
	"html/template"
)

// RenderTemplate parses and executes an HTML template file with the given data.
func RenderTemplate(path string, data any) (string, error) {
	t, err := template.New("email").ParseFiles(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "activation_email.html", data); err != nil {
		// If named template execution fails (e.g., single-file), try default Execute
		buf.Reset()
		if err2 := t.Execute(&buf, data); err2 != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
