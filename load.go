package model

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func Load(filename string) (IModel, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file(%s)", filename)
	}
	defer f.Close()
	fm := FileModel{}
	if err := json.NewDecoder(f).Decode(&fm); err != nil {
		return nil, errors.Wrapf(err, "invalid JSON in file(%s)", filename)
	}
	name := strings.TrimSuffix(path.Base(filename), ".json")
	m, err := New(name, fm)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create model from file definition")
	}
	return m, nil
}
