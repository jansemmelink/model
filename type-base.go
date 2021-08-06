package model

//implements IType
type typeBase struct {
	name string
}

func (t typeBase) Name() string {
	return t.name
}
