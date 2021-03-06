package testdata

const PresentationDefinitionSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"definitions": {
	  "schema": {
		"type": "object",
		"properties": {
		  "uri": {
			"type": "string"
		  },
		  "required": {
			"type": "boolean"
		  }
		},
		"required": [
		  "uri"
		],
		"additionalProperties": false
	  },
	  "filter": {
		"type": "object",
		"properties": {
		  "type": {
			"type": "string"
		  },
		  "format": {
			"type": "string"
		  },
		  "pattern": {
			"type": "string"
		  },
		  "minimum": {
			"type": [
			  "number",
			  "string"
			]
		  },
		  "minLength": {
			"type": "integer"
		  },
		  "maxLength": {
			"type": "integer"
		  },
		  "exclusiveMinimum": {
			"type": [
			  "number",
			  "string"
			]
		  },
		  "exclusiveMaximum": {
			"type": [
			  "number",
			  "string"
			]
		  },
		  "maximum": {
			"type": [
			  "number",
			  "string"
			]
		  },
		  "const": {
			"type": [
			  "number",
			  "string"
			]
		  },
		  "enum": {
			"type": "array",
			"items": {
			  "type": [
				"number",
				"string"
			  ]
			}
		  },
		  "not": {
			"type": "object",
			"minProperties": 1
		  }
		},
		"required": [
		  "type"
		],
		"additionalProperties": false
	  },
	  "format": {
		"type": "object",
		"patternProperties": {
		  "^jwt$|^jwt_vc$|^jwt_vp$": {
			"type": "object",
			"properties": {
			  "alg": {
				"type": "array",
				"minItems": 1,
				"items": {
				  "type": "string"
				}
			  }
			},
			"required": [
			  "alg"
			],
			"additionalProperties": false
		  },
		  "^ldp_vc$|^ldp_vp$|^ldp$": {
			"type": "object",
			"properties": {
			  "proof_type": {
				"type": "array",
				"minItems": 1,
				"items": {
				  "type": "string"
				}
			  }
			},
			"required": [
			  "proof_type"
			],
			"additionalProperties": false
		  },
		  "additionalProperties": false
		},
		"additionalProperties": false
	  },
	  "submission_requirements": {
		"type": "object",
		"oneOf": [
		  {
			"properties": {
			  "name": {
				"type": "string"
			  },
			  "purpose": {
				"type": "string"
			  },
			  "rule": {
				"type": "string"
			  },
			  "count": {
				"type": "integer",
				"minimum": 1
			  },
			  "min": {
				"type": "integer",
				"minimum": 0
			  },
			  "max": {
				"type": "integer",
				"minimum": 0
			  },
			  "from": {
				"type": "string"
			  }
			},
			"required": [
			  "rule",
			  "from"
			],
			"additionalProperties": false
		  },
		  {
			"properties": {
			  "name": {
				"type": "string"
			  },
			  "purpose": {
				"type": "string"
			  },
			  "rule": {
				"type": "string"
			  },
			  "count": {
				"type": "integer",
				"minimum": 1
			  },
			  "min": {
				"type": "integer",
				"minimum": 0
			  },
			  "max": {
				"type": "integer",
				"minimum": 0
			  },
			  "from_nested": {
				"type": "array",
				"minItems": 1,
				"items": {
				  "$ref": "#/definitions/submission_requirements"
				}
			  }
			},
			"required": [
			  "rule",
			  "from_nested"
			],
			"additionalProperties": false
		  }
		]
	  },
	  "input_descriptors": {
		"type": "object",
		"properties": {
		  "id": {
			"type": "string"
		  },
		  "name": {
			"type": "string"
		  },
		  "purpose": {
			"type": "string"
		  },
		  "metadata": {
			"type": "string"
		  },
		  "group": {
			"type": "array",
			"items": {
			  "type": "string"
			}
		  },
		  "schema": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/schema"
			}
		  },
		  "constraints": {
			"type": "object",
			"properties": {
			  "limit_disclosure": {
				"type": "boolean"
			  },
			  "fields": {
				"type": "array",
				"items": {
				  "$ref": "#/definitions/field"
				}
			  },
			  "subject_is_issuer": {
				"type": "string",
				"enum": [
				  "required",
				  "preferred"
				]
			  },
			  "subject_is_holder": {
				"type": "string",
				"enum": [
				  "required",
				  "preferred"
				]
			  }
			},
			"additionalProperties": false
		  }
		},
		"required": [
		  "id",
		  "schema"
		],
		"additionalProperties": false
	  },
	  "field": {
		"type": "object",
		"oneOf": [
		  {
			"properties": {
			  "path": {
				"type": "array",
				"items": {
				  "type": "string"
				}
			  },
			  "purpose": {
				"type": "string"
			  },
			  "filter": {
				"$ref": "#/definitions/filter"
			  }
			},
			"required": [
			  "path"
			],
			"additionalProperties": false
		  },
		  {
			"properties": {
			  "path": {
				"type": "array",
				"items": {
				  "type": "string"
				}
			  },
			  "purpose": {
				"type": "string"
			  },
			  "filter": {
				"$ref": "#/definitions/filter"
			  },
			  "predicate": {
				"type": "string",
				"enum": [
				  "required",
				  "preferred"
				]
			  }
			},
			"required": [
			  "path",
			  "filter",
			  "predicate"
			],
			"additionalProperties": false
		  }
		]
	  }
	},
	"type": "object",
	"properties": {
	  "presentation_definition": {
		"type": "object",
		"properties": {
		  "id": {
			"type": "string"
		  },
		  "name": {
			"type": "string"
		  },
		  "purpose": {
			"type": "string"
		  },
		  "locale": {
			"type": "string"
		  },
		  "format": {
			"$ref": "#/definitions/format"
		  },
		  "submission_requirements": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/submission_requirements"
			}
		  },
		  "input_descriptors": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/input_descriptors"
			}
		  }
		},
		"required": [
		  "id",
		  "input_descriptors"
		],
		"additionalProperties": false
	  }
	}
  }`

const PresentationSubmissionSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"title": "Presentation Submission",
	"type": "object",
	"properties": {
	  "presentation_submission": {
		"type": "object",
		"properties": {
		  "id": {
			"type": "string"
		  },
		  "definition_id": {
			"type": "string"
		  },
		  "locale": {
			"type": "string"
		  },
		  "descriptor_map": {
			"type": "array",
			"items": {
			  "$ref": "#/definitions/descriptor"
			}
		  }
		},
		"required": [
		  "id",
		  "definition_id",
		  "descriptor_map"
		],
		"additionalProperties": false
	  }
	},
	"definitions": {
	  "descriptor": {
		"type": "object",
		"properties": {
		  "id": {
			"type": "string"
		  },
		  "path": {
			"type": "string"
		  },
		  "path_nested": {
			"type": "object",
			"$ref": "#/definitions/descriptor"
		  },
		  "format": {
			"type": "string",
			"enum": [
			  "jwt",
			  "jwt_vc",
			  "jwt_vp",
			  "ldp",
			  "ldp_vc",
			  "ldp_vp"
			]
		  }
		},
		"required": [
		  "id",
		  "path",
		  "format"
		],
		"additionalProperties": false
	  }
	},
	"required": [
	  "presentation_submission"
	],
	"additionalProperties": false
  }`
