package service

import (
	"encoding/json"
	"strings"

	"github.com/gataca-io/vui-core/log"
	"github.com/gataca-io/vui-core/models"
	"github.com/xeipuuv/gojsonschema"
)

type jsonValidator struct{}

// Validate any document with a reference to its json schema or having a string of its json schema.
func (j *jsonValidator) Validate(document models.JSONSchema) error {
	if document.IsRef() {
		return validateWithJSONLoader(gojsonschema.NewReferenceLoader(document.GetSchemaRef()), gojsonschema.NewGoLoader(document))
	}
	return validateWithJSONLoader(gojsonschema.NewStringLoader(document.GetSchema()), gojsonschema.NewGoLoader(document))
}

func (j *jsonValidator) ValidateWithRef(document models.JSONSchema, ref string) error {
	if document.IsRef() {
		return validateWithJSONLoader(gojsonschema.NewReferenceLoader(document.GetSchemaRef()), gojsonschema.NewGoLoader(document))
	}
	return validateWithJSONLoader(gojsonschema.NewStringLoader(ref), gojsonschema.NewGoLoader(document))
}

// Validate exists to hide gojsonschema logic within this file
// it is the entry-point to validation logic, requiring the caller pass in valid json strings for each argument
func (j *jsonValidator) ValidateStrings(schema, document string) error {
	if !isJSON(schema) {
		log.Errorf("schema is not valid json: %s", schema)
		return models.ErrInvalidFormat
	} else if !isJSON(document) {
		log.Errorf("document is not valid json: %s", document)
		return models.ErrInvalidFormat
	}
	return validateWithJSONLoader(gojsonschema.NewStringLoader(schema), gojsonschema.NewStringLoader(document))
}

// validateWithJSONLoader takes schema and document loaders; the document from the loader is validated against
// the schema from the loader. Nil if good, error if bad
func validateWithJSONLoader(schemaLoader, documentLoader gojsonschema.JSONLoader) error {
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Error("Cannot validate document with json schema")
		return err
	}

	if !result.Valid() {
		// Accumulate errs
		var errs []string
		for _, err := range result.Errors() {
			errs = append(errs, err.String())
		}
		log.Error("Cannot validate document with json schema", strings.Join(errs, ","))
		return models.ErrInvalidFormat
	}
	return nil
}

func ValidateJSONSchema(maybeSchema models.JSONSchemaMap) error {
	schemaLoader := gojsonschema.NewSchemaLoader()
	schemaLoader.Validate = true
	return schemaLoader.AddSchemas(gojsonschema.NewStringLoader(maybeSchema.ToJSON()))
}

// True if string is valid JSON, false otherwise
func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
