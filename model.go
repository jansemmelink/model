package model

import (
	"github.com/pkg/errors"
)

type IModel interface {
	Name() string
	Type(name string) IType
	Types() []IType
}

func New(n string, f FileModel) (IModel, error) {
	if !typeNameRegex.MatchString(n) {
		return nil, errors.Errorf("invalid model name=\"%s\"", n)
	}
	m := model{
		name:  n,
		types: map[string]IType{},
	}
	for typeName, ft := range f.Types {
		t, err := NewType(typeName, ft)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid type(%s)", typeName)
		}
		m.types[typeName] = t
	}

	//resolve references between types
	for typeName, ft := range f.Types {
		t, ok := m.types[typeName]
		if !ok {
			return nil, errors.Errorf("model(%s).type(%s) not found for resolve", m.name, typeName)
		}
		if ti, ok := t.(ITypeItem); ok {
			if ft.Item == nil {
				return nil, errors.Errorf("unexpected not found %s after create", typeName)
			}
			if err := ti.Resolve(m, *ft.Item); err != nil {
				return nil, errors.Wrapf(err, "cannot resolve type(%s)", typeName)
			}
		}
	}

	//list children types for item types (do this after parent references were defined in Resolve
	for _, t := range m.types {
		if ti, ok := t.(ITypeItem); ok {
			if ti.Parent() != nil {
				ti.Parent().AddChild(ti)
			}
		}
	}
	return m, nil
}

type model struct {
	name  string
	types map[string]IType
}

func (m model) Name() string {
	return m.name
}

func (m model) Type(name string) IType {
	if t, ok := m.types[name]; ok {
		return t
	}
	//check base type (should implement tree and domains...)
	if t, ok := baseTypes[name]; ok {
		return t
	}
	return nil
}

func (m model) Types() []IType {
	list := []IType{}
	for _, t := range m.types {
		list = append(list, t)
	}
	return list
}

type FileModel struct {
	//IncludeFiles []string            `json:"includes" doc:"List of files to include"`
	Types map[string]FileType `json:"types" doc:"Type definitions"`
}
