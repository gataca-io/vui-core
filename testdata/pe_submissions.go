package testdata

const SampleSubmission = `{
	"presentation_submission": {
	  "id": "a30e3b91-fb77-4d22-95fa-871689c322e2",
	  "definition_id": "32f54163-7166-48f1-93d8-ff217bdb0653",
	  "descriptor_map": [
		{
		  "id": "banking_input_2",
		  "format": "jwt_vc",
		  "path": "$.verifiableCredential[0]"
		},
		{
		  "id": "employment_input",
		  "format": "ldp_vc",
		  "path": "$.verifiableCredential[1]"
		},
		{
		  "id": "citizenship_input_1",
		  "format": "ldp_vc",
		  "path": "$.verifiableCredential[2]"
		}
	  ]
	}
  }`

const SampleVerifiablePresentation = `
{
	"@context": [
	  "https://www.w3.org/2018/credentials/v1",
	  "https://identity.foundation/presentation-exchange/submission/v1"
	],
	"type": [
	  "VerifiablePresentation",
	  "PresentationSubmission"
	],
	"presentation_submission": {
	  "id": "a30e3b91-fb77-4d22-95fa-871689c322e2",
	  "definition_id": "32f54163-7166-48f1-93d8-ff217bdb0653",
	  "descriptor_map": [
		{
		  "id": "banking_input_2",
		  "format": "jwt_vc",
		  "path": "$.verifiableCredential[0]"
		},
		{
		  "id": "employment_input",
		  "format": "ldp_vc",
		  "path": "$.verifiableCredential[1]"
		},
		{
		  "id": "citizenship_input_1",
		  "format": "ldp_vc",
		  "path": "$.verifiableCredential[2]"
		}
	  ]
	},
	"verifiableCredential": [
	  {
		"@context": "https://www.w3.org/2018/credentials/v1",
		"id": "cred:example:aeourhiuq4q38wq8q3",
		"type": [
		  "BankAccount"
		],
		"issuer": "did:example:123",
		"issuanceDate": "2010-01-01T19:43:24Z",
		"credentialSubject": {
		  "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
		  "account": [
			{
			  "id": "1234567890",
			  "route": "DE-9876543210"
			},
			{
			  "id": "2457913570",
			  "route": "DE-0753197542"
			}
		  ]
		},
		"credentialSchema": {
			"id":"https://bank-schemas.org/2.0.0/accounts.json",
			"type":"JSONCredentialSchema2020"
		},
		"proof": {
			"type": "Ed25519VerificationKey2018",
			"created": "2017-06-18T21:19:10Z",
			"proofPurpose": "assertionMethod",
			"verificationMethod": "did:example:123#keys-1",
			"jws": "..."
		  }
	  },
	  {
		"@context": "https://www.w3.org/2018/credentials/v1",
		"id": "cred:example:employee12345",
		"type": [
		  "EmployeeCredential"
		],
		"issuer": "did:example:ebfeb1f712ebc6f1c276e12ec21",
		"issuanceDate": "2010-01-01T19:43:24Z",
		"credentialSubject": {
		  "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
		  "jobs": [{
			"number": "34DGE352",
			"dob": "07/13/80"
		  }]
		},
		"credentialSchema": {
			"id":"https://business-standards.org/schemas/employment-history.json",
			"type":"JSONCredentialSchema2020"
		},
		"proof": {
		  "type": "EcdsaSecp256k1VerificationKey2019",
		  "created": "2017-06-18T21:19:10Z",
		  "proofPurpose": "assertionMethod",
		  "verificationMethod": "did:example:ebfeb1f712ebc6f1c276e12ec21#keys-1",
		  "jws": "..."
		}
	  },
	  {
		"@context": "https://www.w3.org/2018/credentials/v1",
		"id": "https://eu.com/claims/DriversLicense/idnag",
		"type": [
		  "EUDriversLicense"
		],
		"issuer": "did:example:gov2",
		"issuanceDate": "2010-01-01T19:43:24Z",
		"credentialSubject": {
		  "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
		  "license": {
			"number": "34DGE352",
			"dob": "1980-07-13T00:00:00Z"
		  }
		},
		"credentialSchema": {
			"id":"https://eu.com/claims/DriversLicense.json",
			"type":"JSONCredentialSchema2020"
		},
		"proof": {
		  "type": "RsaSignature2018",
		  "created": "2017-06-18T21:19:10Z",
		  "proofPurpose": "assertionMethod",
		  "verificationMethod": "did:example:gov2#keys-1",
		  "jws": "..."
		}
	  }
	],
	"proof": {
	  "type": "RsaSignature2018",
	  "created": "2018-09-14T21:19:10Z",
	  "proofPurpose": "authentication",
	  "verificationMethod": "did:example:ebfeb1f712ebc6f1c276e12ec21#keys-1",
	  "challenge": "1f44d55f-f161-4938-a659-f8026467f126",
	  "domain": "4jt78h47fh47",
	  "jws": "..."
	}
  }`
