package model

type ITypeValue interface {
	IType
}

type typeValueInt struct {
	typeBase
}

type typeValueStr struct {
	typeBase
}

var baseTypes = map[string]IType{}

func init() {
	baseTypes["int"] = typeValueInt{typeBase: typeBase{name: "int"}}
	baseTypes["string"] = typeValueStr{typeBase: typeBase{name: "string"}}
	baseTypes["bool"] = typeValueStr{typeBase: typeBase{name: "bool"}}
	baseTypes["float64"] = typeValueStr{typeBase: typeBase{name: "float64"}}
	baseTypes["amount"] = typeValueStr{typeBase: typeBase{name: "amount"}}
}
