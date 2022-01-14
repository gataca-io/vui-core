package models

type ExchangeRequest interface {
	IsRequest() bool
	ToPresentationDefinition() *PresentationDefinition
}

type ExchangeResponse interface {
	IsResponse() bool
	ToPresentation() *VerifiablePresentation
}
type Verifiable interface {
	GetProofs() *SSIProof
	GetProofChain() *[]Proof
	SetProofs(proof *SSIProof)
	SetProofChain(proof *[]Proof)
}

type LdContext interface {
	GetContext() *SSIContext
	SetContexts(context *SSIContext)
}

type JSONSchema interface {
	GetSchemaRef() string
	GetSchema() string
	IsRef() bool
}

type LdVerifiable interface {
	LdContext
	Verifiable
}

func GetProofForVerificationMethod(v Verifiable, vmethod string) *Proof {
	proofs := v.GetProofs().GetProof()
	for _, p := range *proofs {
		if p.GetCreator() == vmethod {
			return &p
		}
	}
	return nil
}
