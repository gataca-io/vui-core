package testdata

const BasicPresentationDefinition = `{
	"presentation_definition": {
	  "id": "32f54163-7166-48f1-93d8-ff217bdb0653",
	  "input_descriptors": [
		{
		  "id": "banking_input",
		  "name": "Bank Account Information",
		  "purpose": "We need your bank and account information.",
		  "schema": [
			{
			  "uri": "https://bank-standards.com/customer.json"
			}
		  ],
		  "constraints": {
			"limit_disclosure": true,
			"fields": [
			  {
				"path": [
				  "$.issuer",
				  "$.vc.issuer",
				  "$.iss"
				],
				"purpose": "The claim must be from one of the specified issuers",
				"filter": {
				  "type": "string",
				  "pattern": "did:example:123|did:example:456"
				}
			  }
			]
		  }
		},
		{
		  "id": "citizenship_input",
		  "name": "US Passport",
		  "schema": [
			{
			  "uri": "hub://did:foo:123/Collections/schema.us.gov/passport.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.credentialSubject.birth_date",
				  "$.vc.credentialSubject.birth_date",
				  "$.birth_date"
				],
				"filter": {
				  "type": "string",
				  "format": "date",
				  "minimum": "1999-5-16"
				}
			  }
			]
		  }
		}
	  ]
	}
  }`

const GroupPresentationDefinition = `{
	"presentation_definition": {
	  "id": "32f54163-7166-48f1-93d8-ff217bdb0653",
	  "submission_requirements": [
		{
		  "name": "Citizenship Information",
		  "rule": "pick",
		  "count": 1,
		  "from": "A"
		}
	  ],
	  "input_descriptors": [
		{
		  "id": "citizenship_input_1",
		  "name": "EU Driver's License",
		  "group": [
			"A"
		  ],
		  "schema": [
			{
			  "uri": "https://eu.com/claims/DriversLicense.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.issuer",
				  "$.vc.issuer",
				  "$.iss"
				],
				"purpose": "The claim must be from one of the specified issuers",
				"filter": {
				  "type": "string",
				  "pattern": "did:example:gov1|did:example:gov2"
				}
			  },
			  {
				"path": [
				  "$.credentialSubject.dob",
				  "$.vc.credentialSubject.dob",
				  "$.dob"
				],
				"filter": {
				  "type": "string",
				  "format": "date",
				  "maximum": "1999-6-15"
				}
			  }
			]
		  }
		},
		{
		  "id": "citizenship_input_2",
		  "name": "US Passport",
		  "group": [
			"A"
		  ],
		  "schema": [
			{
			  "uri": "hub://did:foo:123/Collections/schema.us.gov/passport.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.credentialSubject.birth_date",
				  "$.vc.credentialSubject.birth_date",
				  "$.birth_date"
				],
				"filter": {
				  "type": "string",
				  "format": "date",
				  "maximum": "1999-5-16"
				}
			  }
			]
		  }
		}
	  ]
	}
  }`

const MultiGroupPresentationDefinition = `{
	  "id": "32f54163-7166-48f1-93d8-ff217bdb0653",
	  "submission_requirements": [
		{
		  "name": "Banking Information",
		  "purpose": "We need to know if you have an established banking history.",
		  "rule": "pick",
		  "count": 1,
		  "from": "A"
		},
		{
		  "name": "Employment Information",
		  "purpose": "We need to know that you are currently employed.",
		  "rule": "all",
		  "from": "B"
		},
		{
		  "name": "Citizenship Information",
		  "rule": "pick",
		  "count": 1,
		  "from": "C"
		}
	  ],
	  "input_descriptors": [
		{
		  "id": "banking_input_1",
		  "name": "Bank Account Information",
		  "purpose": "We need your bank and account information.",
		  "group": [
			"A"
		  ],
		  "schema": [
			{
			  "uri": "https://bank-standards.com/customer.json"
			}
		  ],
		  "constraints": {
			"limit_disclosure": true,
			"fields": [
			  {
				"path": [
				  "$.issuer",
				  "$.vc.issuer",
				  "$.iss"
				],
				"purpose": "The claim must be from one of the specified issuers",
				"filter": {
				  "type": "string",
				  "pattern": "did:example:123|did:example:456"
				}
			  },
			  {
				"path": [
				  "$.credentialSubject.account[*].account_number",
				  "$.vc.credentialSubject.account[*].account_number",
				  "$.account[*].account_number"
				],
				"purpose": "We need your bank account number for processing purposes",
				"filter": {
				  "type": "string",
				  "minLength": 10,
				  "maxLength": 12
				}
			  },
			  {
				"path": [
				  "$.credentialSubject.account[*].routing_number",
				  "$.vc.credentialSubject.account[*].routing_number",
				  "$.account[*].routing_number"
				],
				"purpose": "You must have an account with a German, US, or Japanese bank account",
				"filter": {
				  "type": "string",
				  "pattern": "^DE|^US|^JP"
				}
			  }
			]
		  }
		},
		{
		  "id": "banking_input_2",
		  "name": "Bank Account Information",
		  "purpose": "We need your bank and account information.",
		  "group": [
			"A"
		  ],
		  "schema": [
			{
			  "uri": "https://bank-schemas.org/1.0.0/accounts.json"
			},
			{
			  "uri": "https://bank-schemas.org/2.0.0/accounts.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.issuer",
				  "$.vc.issuer",
				  "$.iss"
				],
				"purpose": "The claim must be from one of the specified issuers",
				"filter": {
				  "type": "string",
				  "pattern": "did:example:123|did:example:456"
				}
			  },
			  {
				"path": [
				  "$.credentialSubject.account[*].id",
				  "$.vc.credentialSubject.account[*].id",
				  "$.account[*].id"
				],
				"purpose": "We need your bank account number for processing purposes",
				"filter": {
				  "type": "string",
				  "minLength": 10,
				  "maxLength": 12
				}
			  },
			  {
				"path": [
				  "$.credentialSubject.account[*].route",
				  "$.vc.credentialSubject.account[*].route",
				  "$.account[*].route"
				],
				"purpose": "You must have an account with a German, US, or Japanese bank account",
				"filter": {
				  "type": "string",
				  "pattern": "^DE|^US|^JP"
				},
				"predicate":"required"
			  }
			]
		  }
		},
		{
		  "id": "employment_input",
		  "name": "Employment History",
		  "purpose": "We need to know your work history.",
		  "group": [
			"B"
		  ],
		  "schema": [
			{
			  "uri": "https://business-standards.org/schemas/employment-history.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.credentialSubject.jobs[*].active"
				],
				"filter": {
					"type": "boolean",
					"pattern": "true"
				},
				"predicate":"preferred"
			  }
			]
		  }
		},
		{
		  "id": "citizenship_input_1",
		  "name": "EU Driver's License",
		  "group": [
			"C"
		  ],
		  "schema": [
			{
			  "uri": "https://eu.com/claims/DriversLicense.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.issuer",
				  "$.vc.issuer",
				  "$.iss"
				],
				"purpose": "The claim must be from one of the specified issuers",
				"filter": {
				  "type": "string",
				  "pattern": "did:example:gov1|did:example:gov2"
				}
			  },
			  {
				"path": [
				  "$.credentialSubject.dob",
				  "$.credentialSubject.license.dob",
				  "$.vc.credentialSubject.dob",
				  "$.dob"
				],
				"filter": {
				  "type": "string",
				  "format": "date",
				  "maximum": "1996-05-16T00:00:00Z"
				}
			  }
			]
		  }
		},
		{
		  "id": "citizenship_input_2",
		  "name": "US Passport",
		  "group": [
			"C"
		  ],
		  "schema": [
			{
			  "uri": "hub://did:foo:123/Collections/schema.us.gov/passport.json"
			}
		  ],
		  "constraints": {
			"fields": [
			  {
				"path": [
				  "$.credentialSubject.birth_date",
				  "$.vc.credentialSubject.birth_date",
				  "$.birth_date"
				],
				"filter": {
				  "type": "string",
				  "format": "date",
				  "minimum": "1999-5-16"
				}
			  }
			]
		  }
		}
	  ]
  }`
