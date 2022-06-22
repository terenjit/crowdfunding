package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var notShowErrorListType = map[string]bool{
	"condition_else": true, "condition_then": true,
}

var jsonSchemaList = map[string]*gojsonschema.Schema{}

// LoadValidatorSchemas all schema from given path
func LoadValidatorSchemas(path string) {

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		if strings.HasSuffix(fileName, ".json") {
			s, err := ioutil.ReadFile(p)
			if err != nil {
				panic(fmt.Errorf("%s: %v", fileName, err))
			}

			var data map[string]interface{}
			if err := json.Unmarshal(s, &data); err != nil {
				panic(fmt.Errorf("%s: %v", fileName, err))
			}
			id, ok := data["$id"].(string)
			if !ok {
				id = strings.Trim(strings.TrimSuffix(strings.TrimPrefix(p, path), ".json"), "/") // take filename without extension
			}
			jsonSchemaList[id], err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(s))
			if err != nil {
				panic(fmt.Errorf("%s: %v", fileName, err))
			}
		}
		return nil
	})
}

// GetSchema json schema by ID
func GetSchema(schemaID string) (*gojsonschema.Schema, error) {
	schema, ok := jsonSchemaList[schemaID]
	if !ok {
		return nil, fmt.Errorf("schema '%s' not found", schemaID)
	}

	return schema, nil
}

// ValidateSchema from Go data type
func ValidateSchema(schemaID string, input interface{}) error {
	multiError := NewMultiError()

	schema, err := GetSchema(schemaID)
	if err != nil {
		multiError.Append("getSchema", err)
		return multiError
	}

	document := gojsonschema.NewGoLoader(input)
	return validate(schema, document)
}

// ValidateDocument document
func ValidateDocument(schemaID string, jsonByte []byte) error {
	multiError := NewMultiError()

	schema, err := GetSchema(schemaID)
	if err != nil {
		multiError.Append("getSchema", err)
		return multiError
	}

	document := gojsonschema.NewBytesLoader(jsonByte)
	return validate(schema, document)
}

func validate(schema *gojsonschema.Schema, document gojsonschema.JSONLoader) error {
	multiError := NewMultiError()

	result, err := schema.Validate(document)
	if err != nil {
		multiError.Append("validateInput", errors.New("gagal memuat input data"))
		return multiError
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			if notShowErrorListType[desc.Type()] {
				continue
			}
			var field = desc.Field()
			if desc.Type() == "required" || desc.Type() == "additional_property_not_allowed" {
				field = fmt.Sprintf("%s.%s", field, desc.Details()["property"])
				field = strings.TrimPrefix(field, "(root).")
			}
			multiError.Append(field, errors.New(desc.Description()))
		}
	}

	if multiError.HasError() {
		return multiError
	}

	return nil
}
