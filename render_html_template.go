package main

import (
	"bytes"
	"text/template"
)

func RenderHTMLTemplate(plan Plan) (string, error) {
	templ, err := template.ParseFiles("template.html")
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = templ.Execute(&buffer, plan)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}