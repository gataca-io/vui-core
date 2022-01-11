package service

import (
	"encoding/json"
	"testing"

	"github.com/gataca-io/vui-core/models"
	"github.com/gataca-io/vui-core/testdata"
	"github.com/stretchr/testify/assert"
)

var jVal JSONValidator

func init() {
	jVal = &jsonValidator{}
}

func TestPresentationDefinitionSchema(t *testing.T) {

	t.Run("Validate the schema is valid", func(t *testing.T) {
		var schemaMap models.JSONSchemaMap
		err := json.Unmarshal([]byte(testdata.PresentationDefinitionSchema), &schemaMap)
		assert.NoError(t, err)

		err = ValidateJSONSchema(schemaMap)
		assert.NoError(t, err)
	})

	t.Run("Validate a sample definition is valid against the schema", func(t *testing.T) {

		// Validate against presDefSchema
		assert.NoError(t, jVal.ValidateStrings(testdata.PresentationDefinitionSchema, testdata.BasicPresentationDefinition))
	})
}

func TestPresentationSubmissionSchema(t *testing.T) {

	t.Run("Validate the schema is valid", func(t *testing.T) {
		var schemaMap models.JSONSchemaMap
		err := json.Unmarshal([]byte(testdata.PresentationSubmissionSchema), &schemaMap)
		assert.NoError(t, err)

		err = ValidateJSONSchema(schemaMap)
		assert.NoError(t, err)
	})

	t.Run("Validate a sample submission is valid against the schema", func(t *testing.T) {
		// Validate against presSubSchema
		assert.NoError(t, jVal.ValidateStrings(testdata.PresentationSubmissionSchema, testdata.SampleSubmission))
	})
}
