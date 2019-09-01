package main

import (
	"os"
	"strings"
	"text/template"
)

type modelTemplateData struct {
	Package  string
	Model    string
	ModelVar string
	Owner    string
	OwnerVar string
}

func Generate(path string, model string) error {
	var d modelTemplateData
	substrings := strings.Split(path, "/")
	d.Package = strings.ToLower(substrings[len(substrings)-1])
	d.ModelVar = strings.ToLower(model)
	d.Model = strings.Title(d.ModelVar)

	t, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(path + "/" + d.ModelVar + ".go")
	if err != nil {
		return err
	}

	err = t.Execute(f, d)
	return err
}

func GenerateWithOwner(path string, model string, owner string) error {
	var d modelTemplateData
	substrings := strings.Split(path, "/")
	d.Package = strings.ToLower(substrings[len(substrings)-1])
	d.ModelVar = strings.ToLower(model)
	d.Model = strings.Title(d.ModelVar)
	d.OwnerVar = strings.ToLower(owner)
	d.Owner = strings.Title(d.OwnerVar)

	t, err := template.New("model with owner").Parse(modelWithOwnerTemplate)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(path + "/" + d.ModelVar + ".go")
	if err != nil {
		return err
	}

	err = t.Execute(f, d)
	return err
}

var modelTemplate = `package {{.Package}}

import "github.com/lmuench/gobdb/gobdb"

type {{.Model}} struct {
	ID int
	// Add attributes here
}

func (self *{{.Model}}) GetID() int {
	return self.ID
}

func (self *{{.Model}}) SetID(ID int) {
	self.ID = ID
}

func GetAll{{.Model}}s(db *gobdb.DB) ([]{{.Model}}, error) {
	{{.ModelVar}}s := []{{.Model}}{}

	err := db.GetAll(&{{.Model}}{}, func(resource interface{}) {
		{{.ModelVar}} := *resource.(*{{.Model}})
		{{.ModelVar}}s = append({{.ModelVar}}s, {{.ModelVar}})
	})
	return {{.ModelVar}}s, err
}

func GetAll{{.Model}}sMap(db *gobdb.DB) (map[int]{{.Model}}, error) {
	{{.ModelVar}}s := make(map[int]{{.Model}})

	err := db.GetAll(&{{.Model}}{}, func(resource interface{}) {
		{{.ModelVar}} := *resource.(*{{.Model}})
		{{.ModelVar}}s[{{.ModelVar}}.GetID()] = {{.ModelVar}}
	})
	return {{.ModelVar}}s, err
}
`

var modelWithOwnerTemplate = `package {{.Package}}

import "github.com/lmuench/gobdb/gobdb"

type {{.Model}} struct {
	ID int
	{{.Owner}}ID int
	// Add attributes here
}

func (self *{{.Model}}) GetID() int {
	return self.ID
}

func (self *{{.Model}}) SetID(ID int) {
	self.ID = ID
}

func (self {{.Model}}) Get{{.Owner}}(db *gobdb.DB) ({{.Owner}}, error) {
	var {{.OwnerVar}} {{.Owner}}
	err := db.Get(&{{.OwnerVar}}, self.{{.Owner}}ID)
	return {{.OwnerVar}}, err
}

func GetAll{{.Model}}s(db *gobdb.DB) ([]{{.Model}}, error) {
	{{.ModelVar}}s := []{{.Model}}{}

	err := db.GetAll(&{{.Model}}{}, func(resource interface{}) {
		{{.ModelVar}} := *resource.(*{{.Model}})
		{{.ModelVar}}s = append({{.ModelVar}}s, {{.ModelVar}})
	})
	return {{.ModelVar}}s, err
}

func GetAll{{.Model}}sMap(db *gobdb.DB) (map[int]{{.Model}}, error) {
	{{.ModelVar}}s := make(map[int]{{.Model}})

	err := db.GetAll(&{{.Model}}{}, func(resource interface{}) {
		{{.ModelVar}} := *resource.(*{{.Model}})
		{{.ModelVar}}s[{{.ModelVar}}.GetID()] = {{.ModelVar}}
	})
	return {{.ModelVar}}s, err
}
`