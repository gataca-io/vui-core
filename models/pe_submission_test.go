package models

import (
	"encoding/json"
	"testing"

	"github.com/gataca-io/vui-core/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPresentationSubmission(t *testing.T) {
	var presSub PresentationSubmissionHolder
	presSubBytes := []byte(testdata.SampleSubmission)
	err := json.Unmarshal(presSubBytes, &presSub)
	assert.NoError(t, err)

	assert.NoError(t, Validate(presSub))

	// Roundtrip and compare
	roundTripBytes, err := json.Marshal(presSub)
	assert.NoError(t, err)

	true, err := JSONBytesEqual(presSubBytes, roundTripBytes)
	assert.NoError(t, err)
	assert.True(t, true)
}
