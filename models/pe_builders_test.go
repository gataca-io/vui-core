package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://identity.foundation/presentation-exchange/#presentation-definition---basic-example
func TestPresentationDefinitionBuilder_BasicExample(t *testing.T) {
	b := NewPresentationDefinitionBuilder()
	b.SetID("32f54163-7166-48f1-93d8-ff217bdb0653")

	// shouldn't validate as empty
	_, err := b.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'InputDescriptors' failed on the 'required' tag")

	// create an input descriptor
	id1 := NewInputDescriptor("banking_input", "Bank Account Information", "We need your bank and account information.", "")
	err = id1.AddSchema(Schema{
		URI: "https://bank-standards.com/customer.json",
	})
	assert.NoError(t, err)

	field := NewConstraintsField([]string{"$.issuer", "$.vc.issuer", "$.iss"})
	field.SetPurpose("The claim must be from one of the specified issuers")
	err = field.SetFilter(Filter{
		Type:    "string",
		Pattern: "did:example:123|did:example:456",
	})
	assert.NoError(t, err)

	// now build
	err = id1.SetConstraints(*field)
	assert.NoError(t, err)

	id1.SetConstraintsLimitDisclosure(true)

	// validate
	err = Validate(id1)
	assert.NoError(t, err)

	// add the descriptor to the builder
	err = b.AddInputDescriptor(*id1)
	assert.NoError(t, err)

	// add a second input descriptor
	id2 := NewInputDescriptor("citizenship_input", "US Passport", "", "")
	err = id2.AddSchema(Schema{
		URI: "hub://did:foo:123/Collections/schema.us.gov/passport.json",
	})
	assert.NoError(t, err)

	field2 := NewConstraintsField([]string{"$.credentialSubject.birth_date", "$.vc.credentialSubject.birth_date", "$.birth_date"})
	err = field2.SetFilter(Filter{
		Type:    "string",
		Format:  "date",
		Minimum: "1999-5-16",
	})
	assert.NoError(t, err)

	err = id2.SetConstraints(*field2)
	assert.NoError(t, err)

	// validate
	err = Validate(id2)
	assert.NoError(t, err)

	// add the descriptor to the builder
	err = b.AddInputDescriptor(*id2)
	assert.NoError(t, err)

	presDef, err := b.Build()
	assert.NoError(t, err)

	assert.NoError(t, Validate(presDef))

	// presDefJSON, err := ToJSON(presDef)
	// assert.NoError(t, err)

	// // get sample json from packr
	// testPresDefJSON, err := testcases.GetJSONFile(testcases.BasicPresentationDefinition)
	// assert.NoError(t, err)

	// // Make sure our builder has the same result
	// same, err := CompareJSON(presDefJSON, testPresDefJSON)
	// assert.NoError(t, err)
	// assert.True(t, same)
}

// TODO as the spec is in flux the remaining tests will be implemented once v0.1.0 is finalized

// func TestPresentationDefinitionBuilder_SingleGroupExample(t *testing.T) {
// }
//
// func TestPresentationDefinitionBuilder_MultiGroupExample(t *testing.T) {
// }

func TestPresentationSubmissionBuilder(t *testing.T) {
	b := NewPresentationSubmissionBuilder("32f54163-7166-48f1-93d8-ff217bdb0653")
	b.SetID("a30e3b91-fb77-4d22-95fa-871689c322e2")

	// shouldn't validate as empty
	_, err := b.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'DescriptorMap' failed on the 'required' tag")

	// Add descriptors
	err = b.AddDescriptor(Descriptor{
		ID:     "banking_input_2",
		Format: CredentialFormat(JWTVC),
		Path:   "$.verifiableCredential.[0]",
	})
	assert.NoError(t, err)

	err = b.AddDescriptor(Descriptor{
		ID:     "employment_input",
		Format: CredentialFormat(LDPVC),
		Path:   "$.verifiableCredential.[1]",
	})
	assert.NoError(t, err)

	err = b.AddDescriptor(Descriptor{
		ID:     "citizenship_input_1",
		Format: CredentialFormat(LDPVC),
		Path:   "$.verifiableCredential.[2]",
	})
	assert.NoError(t, err)

	presSub, err := b.Build()
	assert.NoError(t, err)
	assert.NoError(t, Validate(presSub))

	// presSubJSON, err := util.ToJSON(presSub)
	// assert.NoError(t, err)

	// // get sample json from packr
	// testPresSubJSON, err := testcases.GetJSONFile(testcases.SampleSubmission)
	// assert.NoError(t, err)

	// // Make sure our builder has the same result
	// same, err := util.CompareJSON(presSubJSON, testPresSubJSON)
	// assert.NoError(t, err)
	// assert.True(t, same)
}
