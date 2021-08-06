package model

import (
	"strings"

	"github.com/pkg/errors"
)

type ITypeItem interface {
	IType
	Parent() ITypeItem

	Resolve(m IModel, fti FileTypeItem) error
	Fields() []IField
}

type IField interface {
	Name() string
	IsRef() bool
	Type() IType
}

func newTypeItem(b typeBase, fti FileTypeItem) (IType, error) {
	t := &typeItem{
		typeBase: b,
		parent:   nil, //set in Resolve()
		fields:   map[string]*typeItemField{},
	}

	for fieldName := range fti.Fields {
		if !typeNameRegex.MatchString(fieldName) {
			return nil, errors.Errorf("invalid field name=\"%s\"", fieldName)
		}
		f := &typeItemField{
			name:      fieldName,
			valueType: nil, //set in Resolve()
		}
		t.fields[fieldName] = f
	}
	return t, nil
}

//implements IType
type typeItem struct {
	typeBase
	parent ITypeItem
	fields map[string]*typeItemField
}

func (t typeItem) Parent() ITypeItem {
	return t.parent
}

func (t typeItem) Fields() []IField {
	list := []IField{}
	for _, f := range t.fields {
		list = append(list, f)
	}
	return list
}

//called after all types were added to the model, to link between types
//fti is this typeItem's own description in the file where we get rest of field definition from
func (t *typeItem) Resolve(m IModel, fti FileTypeItem) error {
	if fti.Parent != "" {
		pt := m.Type(fti.Parent)
		if pt == nil {
			return errors.Errorf("unknown parent(%s)", fti.Parent)
		}
		pit, ok := pt.(ITypeItem)
		if !ok {
			return errors.Errorf("parent(%s) is not an item type", fti.Parent)
		}
		//loop through parent's parents to detect circular references
		temp := pit
		for temp != nil {
			temp = temp.Parent()
			if temp == t {
				return errors.Errorf("parent(%s) cause circular reference", fti.Parent)
			}
		}
		t.parent = pit
	}
	for fieldName, ftif := range fti.Fields {
		//get this field created in NewTypeItem()
		f, ok := t.fields[fieldName]
		if !ok {
			return errors.Errorf("did not get field[%s] in typeItem(%s)", fieldName, t.name)
		}

		fieldTypeName := ftif.TypeRef
		f.isRef = false
		if strings.HasPrefix(fieldTypeName, "->") {
			f.isRef = true
			fieldTypeName = fieldTypeName[2:]
		}
		f.valueType = m.Type(fieldTypeName)
		if f.valueType == nil {
			return errors.Errorf("field(%s) has unknown type(%s)", fieldName, fieldTypeName)
		}
	}
	return nil
}

type typeItemField struct {
	name      string
	valueType IType
	isRef     bool
}

func (f typeItemField) Name() string {
	return f.name
}

func (f typeItemField) IsRef() bool {
	return f.isRef
}

func (f typeItemField) Type() IType {
	return f.valueType
}

type FileTypeItem struct {
	Parent string                       `json:"parent,omitempty"` //optional name of parent item type, may not be circular
	Fields map[string]FileTypeItemField `json:"fields"`           //fields describing this type
}

type FileTypeItemField struct {
	TypeRef string `json:"type" doc:"Name the type of value in this field, or ->type to refer to an item of that type"` //type of field, using a *FieldType so same type can be used in multiple places
}
