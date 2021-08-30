package template

import (
	"bytes"
	"html/template"
	"io/ioutil"
)

func Template(name string, file string, v interface{}) (string, error) {
	var str bytes.Buffer
	m, err := ioutil.ReadFile(file)
	t, err := template.New(name).Parse(string(m))
	err = t.Execute(&str, v)
	return str.String(), err
}