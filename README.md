# vui-spec

Open source library to support the usage of the VUI interfaces on any golang project

This is not a complete model, but it supports just the core validation logic.

It provides support to any verifier that wants to integrate the validations of a full VUI Presentation Exchange.

It holds the models for the:

- VUI Extendended Presentation Exchange:
  - Presentation Definition with extensions
  - Presentation submission with Extensions
- Data Agreements
- Trusted Issuers for Issuer Resolution
- Universal DID Resolver

It also provides some implementated services:

- To demonstrate the validations performed on a Presentation Exchange.
- To invoke a universal DID resolver.
- To operate with Data Agreements.

The library will be further completed with an implementation of a Governance Service making use of the Issuer Resolution interface, when the EBSI Trusted Issuer Registry is fully functional.

In order to use this library on a real-working project, it should be completed with additional repositories, services and controllers:

- Controllers supporting the DID-SIOP v2 specification
- A service to validate verifiable objects (credentials, data agreements, presentations) adhering to the interface stated.
- DB repositories to store the different documents in place on each flow.

All the required service and respository interfaces have been defined.
