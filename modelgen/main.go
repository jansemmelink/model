package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/jansemmelink/model"
	"github.com/pkg/errors"
)

func main() {
	modelFile := flag.String("f", "./example.json", "Model description as a JSON file")
	genPath := flag.String("o", "./example", "Output directory")
	flag.Parse()

	//load model file
	m, err := model.Load(*modelFile)
	if err != nil {
		panic(err)
	}
	if err := generate(m, *genPath); err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Successfully generated code in %s\n", *genPath)
} //main()

func generate(m model.IModel, dir string) error {
	mkdir(dir)

	//generate system package
	mkdir(dir + "/system")
	for _, t := range m.Types() {
		if it, ok := t.(model.ITypeItem); ok {
			if err := generateItemType(dir+"/system", it); err != nil {
				return errors.Wrapf(err, "failed to generate %s", it.Name())
			}
		} //if item type
	} //for each type

	//generate api package
	mkdir(dir + "/api")
	if err := generateModelApi(dir+"/api", m); err != nil {
		return errors.Wrapf(err, "failed to generate api")
	}
	return nil
}

func generateModelApi(dir string, m model.IModel) error {
	filename := dir + "/" + m.Name() + "_gen.go"
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to create file(%s)", filename)
	}
	defer f.Close()

	type apiRouteData struct {
		Method  string
		Path    string
		Handler string
	}
	type apiData struct {
		Package string
		Routes  []apiRouteData
		Items   []itemData
	}

	//template data
	modelData := apiData{
		Package: "api",
		Routes:  []apiRouteData{},
		Items:   []itemData{},
	}

	//using gorilla pat is not the best router, and routes has to be in sequence to match long before short patterns

	//start with direct routes, regardless of parent, when you have the id of a type, there is a direct get route
	//because you do not have to query with parent references, and security is provided by middle layer, so cannot be abused
	for _, t := range m.Types() {
		if ti, ok := t.(model.ITypeItem); ok {
			tnames := names(t.Name())

			if ti.Parent() == nil {
				//direct root item lists (no parent), e.g. GET /clients
				modelData.Routes = append(modelData.Routes,
					apiRouteData{
						Method:  "Get",
						Path:    fmt.Sprintf("/%ss", tnames.Ext),
						Handler: "get" + tnames.Pub + "s",
					},
				)
			}

			//child list operations on item, e.g. GET /product/{product_id}/variants
			for _, c := range m.Types() {
				if ti, ok := c.(model.ITypeItem); ok {
					if ti.Parent() == t {
						cnames := names(c.Name())
						modelData.Routes = append(modelData.Routes,
							apiRouteData{
								Method:  "Get",
								Path:    fmt.Sprintf("/%s/{%s_id}/%ss", tnames.Ext, tnames.Ext, cnames.Ext),
								Handler: "get" + tnames.Pub + cnames.Pub + "s",
							},
						)
					}
				}
			}

			//todo: put, post, del handlers

			//todo: add admin routes where you e.g. get children for multiple parents
			//e.g. Get /users where client in 21,22,23...

			//todo: add more api to search with more fields
			//e.g. GET '/users?client_id=21,22,23&name=*abc*'

			//direct operation on item, e.g. GET /variant/{variant_id}
			modelData.Routes = append(modelData.Routes,
				apiRouteData{
					Method:  "Get",
					Path:    fmt.Sprintf("/%s/{%s_id}", tnames.Ext, tnames.Ext),
					Handler: "get" + tnames.Pub,
				},
			)

			modelData.Items = append(modelData.Items, getItemData(ti))
		}
	}

	//generat template into file
	t, err := template.ParseFiles("./templates/api.go.tmpl")
	if err != nil {
		return errors.Wrapf(err, "failed to parse template")
	}
	if err := t.Execute(f, modelData); err != nil {
		return errors.Wrapf(err, "failed to execute template")
	}
	return nil
}

func mkdir(dir string) {
	if err := os.Mkdir(dir, 0755); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}
}

func generateItemType(dir string, ti model.ITypeItem) error {
	filename := dir + "/" + ti.Name() + "_gen.go"
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to create file(%s)", filename)
	}
	defer f.Close()

	//generat template into file
	t, err := template.ParseFiles("./templates/item-type.go.tmpl")
	if err != nil {
		return errors.Wrapf(err, "failed to parse template")
	}
	if err := t.Execute(f, getItemData(ti)); err != nil {
		return errors.Wrapf(err, "failed to execute template")
	}
	return nil
}

func getItemData(ti model.ITypeItem) itemData {
	//template data
	itemData := itemData{
		Package:  "system",
		Type:     itemTypeData{Name: names(ti.Name())},
		Fields:   []itemFieldData{},
		Children: []itemData{},
	}
	if ti.Parent() != nil {
		itemData.Parent = &itemParentData{
			Name: names(ti.Parent().Name()),
		}
	}
	itemData.FieldsExtCSV = ""
	for _, f := range ti.Fields() {
		fieldData := itemFieldData{
			Name:  names(f.Name()),
			Type:  itemTypeData{names(f.Type().Name())},
			IsRef: f.IsRef(),
		}
		itemData.Fields = append(itemData.Fields, fieldData)
		if !f.IsRef() {
			itemData.FieldsExtCSV += ",`" + fieldData.Name.Ext + "`"
		} else {
			itemData.FieldsExtCSV += ",`" + fieldData.Name.Ext + "_id`"
		}
	}
	if len(itemData.FieldsExtCSV) > 0 {
		itemData.FieldsExtCSV = itemData.FieldsExtCSV[1:] //skip leading comma
	}
	itemData.TitleFieldsExtCSV = itemData.FieldsExtCSV

	for _, child := range ti.Children() {
		childItemData := getItemData(child)
		itemData.Children = append(itemData.Children, childItemData)
	}
	return itemData
}

func prvName(s string) string {
	cc := strings.ToLower(s[0:1])
	var last rune = rune(s[0])
	for _, c := range s[1:] {
		if last == '_' && unicode.IsLetter(c) {
			cc = cc[0:len(cc)-1] + strings.ToUpper(string(c))
		} else {
			cc += string(c)
		}
		last = c
	}
	return cc
}

func pubName(s string) string {
	cc := strings.ToUpper(s[0:1])
	var last rune = rune(s[0])
	for _, c := range s[1:] {
		if last == '_' && unicode.IsLetter(c) {
			cc = cc[0:len(cc)-1] + strings.ToUpper(string(c))
		} else {
			cc += string(c)
		}
		last = c
	}
	return cc
}

func names(ext_name string) Names {
	return Names{
		Ext: ext_name,          //e.g. "user_role"
		Prv: prvName(ext_name), //e.g. "userRole"
		Pub: pubName(ext_name), //e.g. "UserRole"
	}
}

type Names struct {
	Ext string
	Prv string
	Pub string
}

type itemData struct {
	Package           string
	Type              itemTypeData
	Parent            *itemParentData
	Fields            []itemFieldData
	FieldsExtCSV      string
	TitleFieldsExtCSV string
	Children          []itemData
}

type itemTypeData struct {
	Name Names
}

type itemParentData struct {
	Name Names
}

type itemFieldData struct {
	Name  Names
	Type  itemTypeData
	IsRef bool
}
