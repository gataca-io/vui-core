# vui-spec

Open source library to support the usage of the VUI interfaces on any golang project

This is not a complete model, but it supports just the core validation logic.

It provides support to any verifier that wants to integrate the validations of a full VUI Presentation Exchange.

This library is the result of the efforts have been pushed under the NGI EssifLab.

![Essif Lab Logo](docs/essiflabLogo.png)

## Features

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

## Pending work

The library will be further completed with an implementation of a Governance Service making use of the Issuer Resolution interface, when the EBSI Trusted Issuer Registry is fully functional.

## Usage

The applications can be located under the "vui" folder.

In order to use this library on a real-working project, it should be completed by implementing additional repositories and services :

- A SSIService service to validate verifiable objects (credentials, data agreements, presentations) adhering to the interface stated.
- DB repositories to store the different documents in place on each flow:
  - Presentation Exchanges
  - Data Agreements
  - Tenants

All the required service and respository interfaces have been defined.

## Test integration

- The process to test presentation exchange is described in the [VUI Specification](https://gataca-io.github.io/vui/):
- The application APIs are described in the [Swagger documentation](https://gataca-io.github.io/vui-core/index.html)

Having developed a complete SSI application using this library, the process to test should look like:

1. Use the API Controller to create a new presentation definition with an embedded data agreement template
2. Process the presentation definition:
   1. Select the credentials matching the input descriptors and submission requirements.
   2. Fill the data agreement with the selected credentials
   3. Submit the new data agreement
   4. Build the presentation submission referencing the data agreement
3. Submit the presentation to the corresponding endpoint
4. Check the validation results and compare with the expected behavior

Additionally, there is a working implemented application with this library that can be used for testing, available at:

```
https://connect.gataca.io:9090
```
