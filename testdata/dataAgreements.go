package testdata

const DataAgreementUnsignedExample = `{
    "@context": "https://schema.igrant.io/data-agreements/v1",
    "id": "d7216cb1-aedb-471e-96f7-7fef51dedb76",
    "version": "v1.0",
    "template_id": "91be609a-4acd-468f-b37a-0f379893b65c",
    "template_version": "v1.0",
    "data_receiver": {
        "id": "did:example:123",
        "name": "Happy Shopping AB",
        "url": "www.happyshopping.com",
        "service": "Example Service",
        "consent_duration": 365,
        "form_of_consent": "explicit | implicit"
    },
    "termination_timestamp": 1660992340,
    "purposes": [{
        "id": "Customized shopping experience",
        "purpose_description": "Collecting user data for offering custom tailored shopping experience",
        "legal_basis": "<consent/legal_obligation/contract/vital_interest/public_task/legitimate_interest>",
        "method_of_use": "<null/data-source/data-using-service>",
        "data_policy": {
            "policy_URL": "https://happyshoping.com/privacy-policy/",
            "jurisdictions": ["Sweden"],
            "industry_scope": "Retail",
            "data_retention_period": 30,
            "geographic_restriction": "Europe",
            "storage_location": "Europe"
        }
    }],

    "data_subject": "did:example:987",
    "data_holder": "did:example:8173",
    "personal_data": [{
            "attribute_id": "f216cb1-aedb-571e-46f7-2fef51dedb54",
            "attribute_name": "Name",
            "attribute_sensitive": true,
            "purposes": ["Customized shopping experience"]
        },
        {
            "attribute_id": "f216cb1-aedb-571e-46f7-2fef51dedb54",
            "attribute_name": "Age",
            "attribute_sensitive": true,
            "purposes": ["Customized shopping experience"]
        }
    ],
    "dpia": {
        "dpia_date": "2021-05-08T08:41:59+0000",
        "dpia_summary_url": "https://org.com/dpia_results.html"
    },

    "event": [{
            "timestamp": 1660992340,
            "principle-did": "did:mydata:1:<sender_did_value>",
            "state": "Definition/Preparation/Capture/Modification/Revocation",
            "proof": {
                "type": "RsaSignature2018",
                "created": "2017-06-18T21:19:10Z",
                "proofPurpose": "assertionMethod",
                "verificationMethod": "did:gatc:50e05cff9d34cc409f2095d4006698167821f439c263ff4f9693274df71ce7e2",
                "jws": "eyJhfahfrarariughair"
            }
        },
        {
            "timestamp": 1660992340,
            "principle-did": "did:mydata:1:<sender_did_value>",
            "state": "<Definition/Prepration/Capture>"
        }
    ]
}
`
