package main

type fieldType string

const (
	TypeUnknown fieldType = "interface{}"
	TypeUInt    fieldType = "uint"
	TypeInt     fieldType = "int"
	TypeFloat   fieldType = "float64"
	TypeInt32   fieldType = "int32"
	TypeUInt32  fieldType = "uint32"
	TypeString  fieldType = "string"
	TypeInt8    fieldType = "int8"
	TypeUInt8   fieldType = "uint8"
)
