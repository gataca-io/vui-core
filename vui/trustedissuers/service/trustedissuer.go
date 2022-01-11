package service

import (
	"time"

	"github.com/gataca-io/vui-core/models"
	"github.com/gataca-io/vui-core/vui/trustedissuers"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

const (
	id1     = "e20993d1-2430-462b-a9d0-2f2ead3345f8"
	id2     = "e30121a2-3023-875b-a1d4-2f2ead3312d2"
	issuer1 = `{  "@context": "https://gataca.io/schemas/tir/2020/v1",  "accreditations": [    {      "accreditor": "did:gatc:2abcd...ABC",      "createdAt": {},      "credentialSchema": "https://gataca.io/tsr/exampleSchema1.json",      "evidence": {        "evidenceDocuments": [          "Passport"        ],        "documentPresence": "Physical",        "id": "https://essif.europa.eu/tsr/evidence/f2aeec97-fc0d-42bf-8ca7-0548192d4231",        "subjectPresence": "Physical",        "type": [          "DocumentVerification"        ]      },      "expirationDate": {},      "levelOfTrust": 2,      "proof": [        {          "created": {},          "jws": "abc...123=",          "proofPurpose": "assertionMethod",          "type": "EidasSeal2021",          "verificationMethod": "did:gatc:2abcd...ABC#123456789"        }      ],      "validFrom": {}    }  ],  "eidasCertificates": [    {      "eidasCertificateIssuerNumber": "123456",      "eidasCertificatePem": "<PEM-ENC-CERT>",      "eidasCertificateSerialNumber": "123456"    }  ],  "dids": [    "did:gatc:2abcd...ABC#123456789",    "did:ebsi:2abcd...ABC#123456789"  ],  "domain": "Education",  "id": "e20993d1-2430-462b-a9d0-2f2ead3345f8",  "organizationInfo": {    "areaGroup": "Education",    "EORI": "AT123456789101",    "discoveryURL": "https://example.organization.com",    "domainName": "https://example.organization.com",    "identifierBag": "ddd1ebce-8305-4edf-b6b6-7588aa021311",    "legalAddress": "Example Street, 38, 3 Izq, Madrid, Spain",    "legalName": "Example Legal Name",    "legalPersonalIdentifier": "123456789",    "LEI": "12341212EXAMPLE34512",    "preferredDisplay": {      "background": {        "color": "#ABCDEF",        "uri": "https://example.org/background.jpg"      },      "logo": "https://example.org/logo.jpg",      "preferredName": "Brand Name",      "style": {        "color": "#ABCDEF",        "fontFamily": "arial"      }    },    "SEED": "AT12345678910",    "SIC": "1234",    "taxReference": "123456789",    "VATRegistration": "ATU12345678"  },  "proof": [    {      "created": {},      "jws": "abc...123=",      "proofPurpose": "assertionMethod",      "type": "EidasSeal2021",      "verificationMethod": "did:gatc:2abcd...ABC#123456789"    }  ],  "serviceEndpoints": [    {      "id": "did:gatc:2abcd...ABC#123456789#openid",      "serviceEndpoint": "https://openid.example.com/",      "type": "OpenIdConnectVersion1.0Service"    }  ]}`
	issuer2 = `{  "@context": "https://gataca.io/schemas/tir/2020/v1",  "accreditations": [    {      "accreditor": "did:gatc:2abcd...ABC",      "createdAt": {},      "credentialSchema": "https://gataca.io/tsr/exampleSchema1.json",      "evidence": {        "evidenceDocuments": [          "Passport"        ],        "documentPresence": "Physical",        "id": "https://essif.europa.eu/tsr/evidence/f2aeec97-fc0d-42bf-8ca7-0548192d4231",        "subjectPresence": "Physical",        "type": [          "DocumentVerification"        ]      },      "expirationDate": {},      "levelOfTrust": 2,      "proof": [        {          "created": {},          "jws": "abc...123=",          "proofPurpose": "assertionMethod",          "type": "EidasSeal2021",          "verificationMethod": "did:gatc:2abcd...ABC#123456789"        }      ],      "validFrom": {}    }  ],  "eidasCertificates": [    {      "eidasCertificateIssuerNumber": "123456",      "eidasCertificatePem": "<PEM-ENC-CERT>",      "eidasCertificateSerialNumber": "123456"    }  ],  "dids": [    "did:gatc:2abcd...ABC#123456789",    "did:ebsi:2abcd...ABC#123456789"  ],  "domain": "Education",  "id": "e30121a2-3023-875b-a1d4-2f2ead3312d2",  "organizationInfo": {    "areaGroup": "Education",    "EORI": "AT123456789101",    "discoveryURL": "https://example.organization.com",    "domainName": "https://example.organization.com",    "identifierBag": "ddd1ebce-8305-4edf-b6b6-7588aa021311",    "legalAddress": "Example Street, 38, 3 Izq, Madrid, Spain",    "legalName": "Example Legal Name",    "legalPersonalIdentifier": "123456789",    "LEI": "12341212EXAMPLE34512",    "preferredDisplay": {      "background": {        "color": "#ABCDEF",        "uri": "https://example.org/background.jpg"      },      "logo": "https://example.org/logo.jpg",      "preferredName": "Brand Name",      "style": {        "color": "#ABCDEF",        "fontFamily": "arial"      }    },    "SEED": "AT12345678910",    "SIC": "1234",    "taxReference": "123456789",    "VATRegistration": "ATU12345678"  },  "proof": [    {      "created": {},      "jws": "abc...123=",      "proofPurpose": "assertionMethod",      "type": "EidasSeal2021",      "verificationMethod": "did:gatc:2abcd...ABC#123456789"    }  ],  "serviceEndpoints": [    {      "id": "did:gatc:2abcd...ABC#123456789#openid",      "serviceEndpoint": "https://openid.example.com/",      "type": "OpenIdConnectVersion1.0Service"    }  ]}`
)

type trustedIssuerService struct {
	mock map[string]*models.TrustedIssuer
}

func NewTrustedIssuerService() trustedissuers.TrustedIssuers {
	var i1 *models.TrustedIssuer
	var i2 *models.TrustedIssuer

	err := json.Unmarshal([]byte(issuer1), i1)
	if err != nil {
		panic("could not initialize mock issuer 1")
	}

	err = json.Unmarshal([]byte(issuer2), i2)
	if err != nil {
		panic("could not initialize mock issuer 2")
	}

	m := map[string]*models.TrustedIssuer{
		id1: i1,
		id2: i2,
	}

	return &trustedIssuerService{
		mock: m,
	}
}

func (ts *trustedIssuerService) GetTrustedIssuer(ctx echo.Context, id string) (*models.TrustedIssuer, error) {
	var i *models.TrustedIssuer

	if i = ts.mock[id]; i == nil {
		return nil, models.ErrNotFound
	}

	return i, nil
}

func (ts *trustedIssuerService) GetAllTrustedIssuers(ctx echo.Context) (*models.TrustedIssuerList, error) {
	var iL *models.TrustedIssuerList

	iL.TrustedIssuers = append(iL.TrustedIssuers, *ts.mock[id1])
	iL.TrustedIssuers = append(iL.TrustedIssuers, *ts.mock[id2])

	p := &models.Proof{
		Created:            &models.TimeWithFormat{Time: time.Now()},
		Jws:                "abc...123",
		ProofPurpose:       "assertionMethod",
		Type:               "EidasSeal2021",
		VerificationMethod: "did:gatc:2abcd...ABC#123456789",
	}

	iL.Proof = p

	return iL, nil
}
