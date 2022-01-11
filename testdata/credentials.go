package testdata

const W3CCredential = `{
	"@context": [
	"https://www.w3.org/2018/credentials/v1",
	"https://www.w3.org/2018/credentials/examples/v1"
	],
	"id": "http://example.edu/credentials/3732",
	"type": ["VerifiableCredential", "UniversityDegreeCredential"],
	"credentialSubject": {
		"id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
		"degree": {
			"type": "BachelorDegree",
			"name": "Bachelor of Science and Arts"
		}
	}
}`

const GatacaCredential = `{
	"@context": [
		"https://www.w3.org/2018/credentials/v1",
		"https://s3.eu-west-1.amazonaws.com/gataca.io/contexts/v1.json"
	],
	"credentialStatus": {
		"id": "https://certify.gataca.io:9090/api/v1/group/diplomaUC3M/status",
		"type": "CredentialStatusList2017"
	},
	"credentialSubject": {
		"academicYear": "2012",
		"id": "did:gatc:NGQxZGE5ZTQ3YTk1ODI3MzhkNTA0ZGNi"
	},	
	"id": "cred:gatc:ovz8ekqoblbcmapjc31ku5prwd22s7ksj4ho",
	"type": ["VerifiableCredential", "academicYearCredential"]
}`

const GatacaCredentialEvidence = `{
	"@context": [
		"https://www.w3.org/2018/credentials/v1",
		"https://s3.eu-west-1.amazonaws.com/gataca.io/contexts/v1.json"
	],
	"credentialStatus": {
		"id": "https://certify.gataca.io:9090/api/v1/group/diplomaUC3M/status",
		"type": "CredentialStatusList2017"
	},
	"credentialSubject": {
		"academicYear": "2012",
		"id": "did:gatc:NGQxZGE5ZTQ3YTk1ODI3MzhkNTA0ZGNi"
	},
	"evidence": {
		"documentPresence": "Physical",
		"subjectPresence": "Physical",
		"evidenceDocument": {
			"type": "Passport",
			"documentCode": "P",
			"documentNumber": "SPECI2014",
			"documentIssuingState": "NLD",
			"documentExpirationDate": "2031-06-25T15:05:00Z"
		},
		"id": "https://example.edu/evidence/f2aeec97-fc0d-42bf-8ca7-0548192d4231",
		"type": [
  			"DocumentVerification"
		],
		"verifier": "did:ebsi:2962ASk1ODI3MzhkDAGASD4a7a"
	},	
	"id": "cred:gatc:ovz8ekqoblbcmapjc31ku5prwd22s7ksj4ho",
	"type": ["VerifiableCredential", "academicYearCredential"]
}`

/*


 */
