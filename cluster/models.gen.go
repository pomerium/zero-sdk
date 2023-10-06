// Package cluster provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package cluster

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for BundleStatusFailureSource.
const (
	DatabrokerError BundleStatusFailureSource = "databroker_error"
	DownloadError   BundleStatusFailureSource = "download_error"
	InvalidBundle   BundleStatusFailureSource = "invalid_bundle"
	IoError         BundleStatusFailureSource = "io_error"
	UnknownError    BundleStatusFailureSource = "unknown_error"
)

// BootstrapConfig defines model for BootstrapConfig.
type BootstrapConfig struct {
	// DatabrokerStorageConnection databroker storage connection string
	DatabrokerStorageConnection *string `json:"databrokerStorageConnection,omitempty"`
}

// Bundle defines model for Bundle.
type Bundle struct {
	// Id bundle id
	Id string `json:"id"`
}

// BundleStatus defines model for BundleStatus.
type BundleStatus struct {
	Failure *BundleStatusFailure `json:"failure,omitempty"`
	Success *BundleStatusSuccess `json:"success,omitempty"`
}

// BundleStatusFailure defines model for BundleStatusFailure.
type BundleStatusFailure struct {
	Message string `json:"message"`

	// Source source of the failure
	Source BundleStatusFailureSource `json:"source"`
}

// BundleStatusFailureSource source of the failure
type BundleStatusFailureSource string

// BundleStatusSuccess defines model for BundleStatusSuccess.
type BundleStatusSuccess struct {
	// Metadata bundle metadata
	Metadata map[string]string `json:"metadata"`
}

// DownloadBundleResponse defines model for DownloadBundleResponse.
type DownloadBundleResponse struct {
	// CaptureMetadataHeaders bundle metadata that need be picked up by the client from the download URL
	CaptureMetadataHeaders []string `json:"captureMetadataHeaders"`
	ExpiresInSeconds       string   `json:"expiresInSeconds"`

	// Url download URL
	Url string `json:"url"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Error Error message
	Error string `json:"error"`
}

// ExchangeTokenRequest defines model for ExchangeTokenRequest.
type ExchangeTokenRequest struct {
	// RefreshToken cluster identity token
	RefreshToken string `json:"refreshToken"`
}

// ExchangeTokenResponse defines model for ExchangeTokenResponse.
type ExchangeTokenResponse struct {
	ExpiresInSeconds string `json:"expiresInSeconds"`

	// IdToken ID token
	IdToken string `json:"idToken"`
}

// GetBootstrapConfigResponse defines model for GetBootstrapConfigResponse.
type GetBootstrapConfigResponse = BootstrapConfig

// GetBundlesResponse defines model for GetBundlesResponse.
type GetBundlesResponse struct {
	Bundles []Bundle `json:"bundles"`
}

// BundleId defines model for bundleId.
type BundleId = string

// ReportClusterResourceBundleStatusJSONRequestBody defines body for ReportClusterResourceBundleStatus for application/json ContentType.
type ReportClusterResourceBundleStatusJSONRequestBody = BundleStatus

// ExchangeClusterIdentityTokenJSONRequestBody defines body for ExchangeClusterIdentityToken for application/json ContentType.
type ExchangeClusterIdentityTokenJSONRequestBody = ExchangeTokenRequest
