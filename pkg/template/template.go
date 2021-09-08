package template

import (
	"bytes"
	"html/template"
)

func Template(name string, text string, v interface{}) (string, error) {
	var str bytes.Buffer
	t, err := template.New(name).Parse(string(text))
	err = t.Execute(&str, v)
	return str.String(), err
}