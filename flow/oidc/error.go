package oidc

const (
	ErrOIDCProviderNotFound         = "oidc_provider_not_found"
	ErrOIDCProviderUnableToCreate   = "oidc_provider_unable_to_create"
	ErrOIDCVerifierUnableToGenerate = "oidc_verifier_unable_to_generate"

	ErrOIDCFlowUnableToCreate = "oidc_flow_unable_to_create"
	ErrOIDCFlowNotFound       = "oidc_flow_not_found"

	ErrOIDCUnableToExchangeCode = "oidc_unable_to_exchange_code"
)
