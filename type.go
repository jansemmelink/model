package model

import (
	"regexp"

	"github.com/pkg/errors"
)

type IType interface {
	Name() string
}

func NewType(name string, ft FileType) (IType, error) {
	switch {
	case ft.Item != nil:
		return newTypeItem(typeBase{name: name}, *ft.Item)
	default:
		return nil, errors.Errorf("missing \"item\":{...}") //or... if more types...
	}
}

type FileType struct {
	//one of the following must be defined
	Item *FileTypeItem `json:"item"`
}

const typeNamePattern = `[a-z]([a-z0-9_]*[a-z0-9])*`

var typeNameRegex = regexp.MustCompile(typeNamePattern)
