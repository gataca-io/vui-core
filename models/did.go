package models

import (
	"encoding/json"
	"errors"
	"time"
)

//Request token key to save/get the request token from the context
const RequestToken = "REQUEST_TOKEN"

// Public Key types following https://w3c-ccg.github.io/ld-cryptosuite-registry/#signature-suites and https://w3c.github.io/did-spec-registries/#verification-method-types
const (
	TypeEd25519         = "Ed25519VerificationKey2018" // Uses LD-Context for Serialization
	TypeRSA             = "RsaVerificationKey2018"
	TypeEC              = "EcdsaSecp256k1VerificationKey2019"
	TypeEC2             = "EcdsaSecp256k1RecoveryMethod2020"
	TypeJWS             = "JwsVerificationKey2020"
	TypeGPG             = "GpgVerificationKey2020"
	TypeJCS_Ed          = "JcsEd25519Signature2020" // Thats the one we should be using and supporting up until know
	TypeBBS1            = "Bls12381G1Key2020"
	TypeBBS2            = "Bls12381G2Key2020"
	TypeX25519Agreement = "X25519KeyAgreementKey2019"
	TypeECSchnorr       = "SchnorrSecp256k1VerificationKey2019" //Listed but undefined
)

// DIDDocument represent the a client which is using Gataca
// Ledger is not a standard variable. We need to know where store the new DID when user create.
// See @https://www.w3.org/TR/did-core/#verification-methods
// See @https://w3c.github.io/did-spec-registries
// TODO: Switch PublicKey to Verification Method because of deprecation
type DIDDocument struct {
	Context            *SSIContext          `json:"@context,omitempty" description:"Context for JSON-LD"`
	Assertion          []VerificationMethod `json:"assertionMethod,omitempty" description:"Used to specify how the DID subject is expected to express claims: Ordered set of one or more verification methods. Each verification method MAY be embedded or referenced. "`
	Authentication     []VerificationMethod `json:"authentication,omitempty" description:"Used to specify how the DID subject is expected to be authenticated: Ordered set of one or more verification methods. Each verification method MAY be embedded or referenced."`
	Delegation         []VerificationMethod `json:"capabilityDelegation,omitempty" description:"Used to specify a mechanism that might be used by the DID subject to delegate a cryptographic capability to another party: Ordered set of one or more verification methods. Each verification method MAY be embedded or referenced."`
	EbsiToken          string               `json:"ebsiToken,omitempty" description:"Authorization JWT token used on write operations on the EBSI network provided by the user from the wallet"`
	Id                 string               `json:"id,omitempty" example:"did:example:xxxxxxxxxxxxx" description:"DID unique Identifier for resolution"`
	Invocation         []VerificationMethod `json:"capabilityInvocation,omitempty" description:"Used to specify a verification method that might be used by the DID subject to invoke a cryptographic capability: Ordered set of one or more verification methods. Each verification method MAY be embedded or referenced."`
	KeyAgreement       []VerificationMethod `json:"KeyAgreement,omitempty" description:"Used to specify how to encrypt information to the DID subject: Ordered set of one or more verification methods. Each verification method MAY be embedded or referenced."`
	Ledger             string               `json:"ledger,omitempty" example:"ETH" description:"DLT storing the DID Document"`
	Proof              *SSIProof            `json:"proof,omitempty" description:"Proof validating DID Document and identifying the source and control of the document."`
	PublicKey          []*PublicKey         `json:"publicKey,omitempty" description:"Description of the public keys associated to this DID - Deprecated"`
	Service            *SSIService          `json:"service,omitempty" swaggerignore:"true"`
	VerificationMethod []*PublicKey         `json:"verificationMethod,omitempty" description:"Description of the public keys associated to this DID. Replacing Public Keys"`
}

type DIDDocumentMetadata struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

// PublicKeys are yet usually only Ed25519 Keys -> publicKeyBase58
type PublicKey struct {
	Context      *SSIContext `json:"@context,omitempty" description:"Context for JSON-LD"`
	Controller   string      `json:"controller,omitempty" example:"did:example:xxxxxxxxxxxxx" description:"DID key owner"`
	Id           string      `json:"id,omitempty"  example:"did:example:xxxxxxxxxxxxx#keys-1" description:"Key Identifier as a URI"`
	EthAddress   string      `json:"ethereumAddress,omitempty" example:"0x89a932207c485f85226d86f7cd486a89a24fcc12"`
	KeyB58       string      `json:"publicKeyBase58,omitempty" example:"2pju8d2E3LWkDJaJpm6BBf73v5DzSRyVNVf3JQwgV7DW" description:"Key data codified in base58"`
	KeyGPG       string      `json:"publicKeyGpg,omitempty" example:"-----BEGIN PGP PUBLIC KEY BLOCK-----\r\nVersion: OpenPGP.js v4.9.0\r\nComment: https://openpgpjs.org\r\n\r\nxjMEXkm5LRYJKwYBBAHaRw8BAQdASmfrjYr7vrjwHNiBsdcImK397Vc3t4BL\r\nE8rnN......v6\r\nDw==\r\n=wSoi\r\n-----END PGP PUBLIC KEY BLOCK-----\r\n"`
	KeyHex       string      `json:"publicKeyHex,omitempty"` //Unsupported on DID-Core but widely used
	KeyJwk       *JWK        `json:"publicKeyJwk,omitempty"`
	KeyPem       string      `json:"publicKeyPem,omitempty"`
	KeyMultibase string      `json:"publicKeyMultibase,omitempty"`
	Type         string      `json:"type,omitempty" example:"Ed25519VerificationKey2018" description:"Cryptographic suite definition"` //Ed25519VerificationKey2018
	Usage        string      `json:"usage,omitempty" example:"signing,recovery"`
}

type JWK struct {
	Alg     string `json:"alg,omitempty" example:"jwa values: https://www.rfc-editor.org/rfc/rfc7518.html#page-6"`
	Curve   string `json:"crv,omitempty"`
	D       string `json:"d,omitempty" description:"value of d (ECC private key, RSA private key)"`
	DP      string `json:"dp,omitempty" description:"value of first factor crt exponent (RSA private)"`
	DQ      string `json:"dq,omitempty" description:"value of second factor crt exponent (RSA private)"`
	E       string `json:"e,omitempty" description:"value of exponent (RSA public key)"`
	K       string `json:"k,omitempty" description:"value of the symmetric key (Symmetric, OCT keys)"`
	KeyId   string `json:"kid,omitempty"`
	KeyType string `json:"kty,omitempty"`
	Q       string `json:"q,omitempty" description:"value of second prime number (RSA private key)"`
	QI      string `json:"qi,omitempty" description:"value of first CRT coefficient (RSA private key)"`
	N       string `json:"n,omitempty" description:"value of moduluss (RSA public key)"`
	P       string `json:"p,omitempty" description:"value of first prime number (RSA Private key)"`
	Use     string `json:"use,omitempty" example:"sig,enc"`
	X       string `json:"x,omitempty" description:"value of x coordinate (ECC, OKP public key)"`
	Y       string `json:"y,omitempty" description:"value of y coordinate (ECC public key)"`
}

func (p *PublicKey) GetContext() *SSIContext {
	return p.Context
}

func (p *PublicKey) SetContexts(context *SSIContext) {
	p.Context = context
}

func (p *PublicKey) GetKey() string {
	if p.KeyMultibase != "" {
		return p.KeyMultibase
	}
	switch p.Type {
	case TypeEd25519, TypeJCS_Ed, TypeX25519Agreement:
		return p.KeyB58
	case TypeRSA:
		return p.KeyPem
	case TypeEC, TypeEC2, TypeECSchnorr:
		if p.KeyHex != "" {
			return p.KeyHex
		}
		if p.KeyB58 != "" {
			return p.KeyB58
		}
		if p.KeyJwk != nil {
			return "" //UNSUPPORTED: Handle conversion to Hex
		}
		if p.EthAddress != "" {
			return "" //UNSUPPORTED: Handle obtention from address
		}
	case TypeJWS:
		if p.KeyJwk != nil {
			return "" //UNSUPPORTED: Handle conversion to Hex
		}
	case TypeGPG:
		return p.KeyGPG
	case TypeBBS1, TypeBBS2:
		return p.KeyB58
	}
	return ""
}

/** @Deprecated
// PublicKey represent the a client which is using Gataca
type PublicKeyEd25519 struct {
	Context    *SSIContext `json:"@context,omitempty" example:"https://w3id.org/security/v1" description:"Context for JSON-LD"`
	Id         string      `json:"id,omitempty" example:"did:example:xxxxxxxxxxxxx" description:"DID key owner"`
	Type       string      `json:"type,omitempty" example:"Ed25519VerificationKey2018" description:"Cryptographic suite definition"` //Ed25519VerificationKey2018
	Controller string      `json:"controller,omitempty" example:"did:example:xxxxxxxxxxxxx#keys-1" description:"Key Identifier as a URI"`
	Key        string      `json:"publicKeyBase58,omitempty" example:"2pju8d2E3LWkDJaJpm6BBf73v5DzSRyVNVf3JQwgV7DW" description:"Key data codified in base58"`
}

type PublicKeyRSA struct {
	Context    *SSIContext `json:"@context,omitempty"`
	Id         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty" example:"RsaVerificationKey2018"` //RsaVerificationKey2018
	Controller string      `json:"controller,omitempty"`
	Key        string      `json:"publicKeyPem,omitempty"`
}

type PublicKeySecp256k1 struct {
	Context    *SSIContext `json:"@context,omitempty"`
	Id         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty" example:"EcdsaSecp256k1VerificationKey2019"` //EcdsaSecp256k1VerificationKey2019
	Controller string      `json:"controller,omitempty"`
	Key        string      `json:"publicKeyHex,omitempty"`
}
*/

// Service represent the a client which is using Gataca
type Service struct {
	Id              string `json:"id,omitempty"`
	Type            string `json:"type,omitempty"`
	ServiceEndpoint string `json:"serviceEndpoint,omitempty"`
}

type SSIService struct {
	Service  *Service
	Services *[]Service
}

func (ss *SSIService) GetServices() *[]Service {
	if ss.Service != nil {
		arr := []Service{*ss.Service}
		return &arr
	}
	if ss.Services != nil {
		return ss.Services
	}
	return &[]Service{}
}

func (ss *SSIService) UnmarshalJSON(jsonData []byte) error {
	s := &Service{}
	if errO := json.Unmarshal([]byte(jsonData), s); errO != nil {
		var array []Service
		if errA := json.Unmarshal([]byte(jsonData), &array); errA != nil {
			return errors.New("unrecognized entity")
		}
		ss.Services = &array
	}
	ss.Service = s
	return nil
}

func (ss *SSIService) MarshalJSON() ([]byte, error) {
	if ss.Service != nil {
		return json.Marshal(ss.Service)
	}
	return json.Marshal(ss.Services)
}

type SSIContext struct {
	Context  string
	Contexts []string
}

func (sc *SSIContext) GetContext() []string {
	if sc.Context != "" {
		arr := []string{sc.Context}
		return arr
	}
	if sc.Contexts != nil {
		return sc.Contexts
	}
	return []string{}
}

func (sc *SSIContext) UnmarshalJSON(jsonData []byte) error {
	var array []string
	if errA := json.Unmarshal([]byte(jsonData), &array); errA != nil {
		// Could be a single string
		var s string
		if errS := json.Unmarshal([]byte(jsonData), &s); errS != nil {
			//Not a string, type unkonwn
			return errors.New("unrecognized entity")
		}
		sc.Context = s
		return nil
	}
	sc.Contexts = array
	return nil
}

func (sc SSIContext) MarshalJSON() ([]byte, error) {
	if sc.Context != "" {
		return json.Marshal(sc.Context)
	}
	return json.Marshal(sc.Contexts)
}

type VerificationMethod struct {
	Reference string
	Method    *PublicKey
}

func (vm *VerificationMethod) GetKeys(parentKeys []*PublicKey) *PublicKey {
	if vm.Method != nil {
		return vm.Method
	}
	for _, key := range parentKeys {
		if vm.Reference == key.Id {
			return key
		}
	}
	return nil
}

func (vm *VerificationMethod) GetId() string {
	if vm.Method != nil {
		return vm.Method.Id
	}
	return vm.Reference
}

func (vm *VerificationMethod) UnmarshalJSON(jsonData []byte) error {
	m := PublicKey{}
	if errO := json.Unmarshal([]byte(jsonData), &m); errO != nil {
		// Could be a single string as a reference
		var s string
		if errS := json.Unmarshal([]byte(jsonData), &s); errS != nil {
			//Not a string, type unkonwn
			return errors.New("unrecognized entity")
		}
		vm.Reference = s
		return nil
	}
	vm.Method = &m
	return nil
}

func (vm *VerificationMethod) MarshalJSON() ([]byte, error) {
	if vm.Reference != "" {
		return json.Marshal(vm.Reference)
	}
	return json.Marshal(vm.Method)
}

func (d *DIDDocument) GetProofs() *SSIProof {
	return d.Proof
}

func (d *DIDDocument) GetProofChain() *[]Proof {
	return nil
}

func (d *DIDDocument) SetProofs(proof *SSIProof) {
	d.Proof = proof
}

func (d *DIDDocument) SetProofChain(proof *[]Proof) {
	//No support
}

func (d *DIDDocument) GetContext() *SSIContext {
	return d.Context
}

func (d *DIDDocument) SetContexts(context *SSIContext) {
	d.Context = context
}

//For retrocompatibility with old dids
func (d *DIDDocument) GetVerificationMethods() []*PublicKey {
	keys := []*PublicKey{}
	if d.PublicKey != nil {
		keys = append(keys, d.PublicKey...)
	}
	if d.VerificationMethod != nil {
		keys = append(keys, d.VerificationMethod...)
	}
	return keys
}

func (d *DIDDocument) GetKeyFromId(id string) *PublicKey {
	mergeKeys := d.GetVerificationMethods()

	if len(mergeKeys) > 0 {
		for k, v := range mergeKeys {
			if v.Id == id {
				return mergeKeys[k]
			}
		}
	}
	return nil
}
