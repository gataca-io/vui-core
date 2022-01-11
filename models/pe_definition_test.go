package models

import (
	"encoding/json"
	"testing"

	"github.com/gataca-io/vui-core/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPresentationDefinition(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {

		var presDef PresentationDefinitionHolder
		presDefBytes := []byte(testdata.BasicPresentationDefinition)
		err := json.Unmarshal(presDefBytes, &presDef)
		assert.NoError(t, err)

		assert.NoError(t, Validate(presDef))

		// Roundtrip and compare
		roundTripBytes, err := json.Marshal(presDef)
		assert.NoError(t, err)

		true, err := JSONBytesEqual(presDefBytes, roundTripBytes)
		assert.NoError(t, err)
		assert.True(t, true)
	})

	t.Run("Single Group", func(t *testing.T) {
		var presDef PresentationDefinitionHolder
		presDefBytes := []byte(testdata.GroupPresentationDefinition)
		err := json.Unmarshal(presDefBytes, &presDef)
		assert.NoError(t, err)

		assert.NoError(t, Validate(presDef))

		// Roundtrip and compare
		roundTripBytes, err := json.Marshal(presDef)
		assert.NoError(t, err)

		true, err := JSONBytesEqual(presDefBytes, roundTripBytes)
		assert.NoError(t, err)
		assert.True(t, true)
	})

	t.Run("Multi Group", func(t *testing.T) {

		var presDef PresentationDefinition
		presDefBytes := []byte(testdata.MultiGroupPresentationDefinition)
		err := json.Unmarshal(presDefBytes, &presDef)
		assert.NoError(t, err)

		assert.NoError(t, Validate(presDef))

		// Roundtrip and compare
		roundTripBytes, err := json.Marshal(presDef)
		assert.NoError(t, err)

		true, err := JSONBytesEqual(presDefBytes, roundTripBytes)
		assert.NoError(t, err)
		assert.True(t, true)
	})
}

func TestPresentationDefinitionBuilder(t *testing.T) {
	t.Run("one ldp", func(t *testing.T) {
		// minimally complaint pres def
		b := NewPresentationDefinitionBuilder()
		err := b.AddInputDescriptor(InputDescriptor{
			ID:   "test",
			Name: "test",
			Schema: []Schema{
				{
					URI: "test",
				},
			},
		})
		assert.NoError(t, err)

		err = b.SetLDPFormat(LDP, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)
		pres, err := b.Build()
		assert.NoError(t, err)
		assert.NotNil(t, pres)
	})

	t.Run("one jwt", func(t *testing.T) {
		// minimally complaint pres def
		b := NewPresentationDefinitionBuilder()
		err := b.AddInputDescriptor(InputDescriptor{
			ID:   "test",
			Name: "test",
			Schema: []Schema{
				{
					URI: "test",
				},
			},
		})
		assert.NoError(t, err)

		err = b.SetJWTFormat(JWT, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)
		pres, err := b.Build()
		assert.NoError(t, err)
		assert.NotNil(t, pres)
	})

	t.Run("mixed ldp and jwt", func(t *testing.T) {
		// minimally complaint pres def
		b := NewPresentationDefinitionBuilder()
		err := b.AddInputDescriptor(InputDescriptor{
			ID:   "test",
			Name: "test",
			Schema: []Schema{
				{
					URI: "test",
				},
			},
		})
		assert.NoError(t, err)

		err = b.SetJWTFormat(JWT, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)

		err = b.SetLDPFormat(LDP, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)

		pres, err := b.Build()
		assert.NoError(t, err)
		assert.NotNil(t, pres)
	})

	t.Run("mixed ldps and jwts", func(t *testing.T) {
		// minimally complaint pres def
		b := NewPresentationDefinitionBuilder()
		err := b.AddInputDescriptor(InputDescriptor{
			ID:   "test",
			Name: "test",
			Schema: []Schema{
				{
					URI: "test",
				},
			},
		})
		assert.NoError(t, err)

		err = b.SetJWTFormat(JWT, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)
		err = b.SetJWTFormat(JWTVC, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)
		err = b.SetJWTFormat(JWTVP, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)

		err = b.SetLDPFormat(LDP, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)
		err = b.SetLDPFormat(LDPVC, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)
		err = b.SetLDPFormat(LDPVP, []string{"Ed25519Signature2018"})
		assert.NoError(t, err)

		pres, err := b.Build()
		assert.NoError(t, err)
		assert.NotNil(t, pres)
	})
}
