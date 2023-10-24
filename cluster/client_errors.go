package cluster

import (
	"net/http"

	"github.com/pomerium/zero-sdk/apierror"
)

type EmptyResponse struct{}

var (
	_ apierror.APIResponse[ExchangeTokenResponse]  = (*ExchangeClusterIdentityTokenResp)(nil)
	_ apierror.APIResponse[GetBundlesResponse]     = (*GetClusterResourceBundlesResp)(nil)
	_ apierror.APIResponse[DownloadBundleResponse] = (*DownloadClusterResourceBundleResp)(nil)
	_ apierror.APIResponse[EmptyResponse]          = (*ReportClusterResourceBundleStatusResp)(nil)
)

func (r *ExchangeClusterIdentityTokenResp) GetBadRequestError() (string, bool) {
	if r.JSON400 == nil {
		return "", false
	}
	return r.JSON400.Error, true
}

func (r *ExchangeClusterIdentityTokenResp) GetInternalServerError() (string, bool) {
	if r.JSON500 == nil {
		return "", false
	}
	return r.JSON500.Error, true
}

func (r *ExchangeClusterIdentityTokenResp) GetValue() *ExchangeTokenResponse {
	return r.JSON200
}

func (r *ExchangeClusterIdentityTokenResp) GetHTTPResponse() *http.Response {
	return r.HTTPResponse
}

func (r *GetClusterResourceBundlesResp) GetBadRequestError() (string, bool) {
	if r.JSON400 == nil {
		return "", false
	}
	return r.JSON400.Error, true
}

func (r *GetClusterResourceBundlesResp) GetInternalServerError() (string, bool) {
	if r.JSON500 == nil {
		return "", false
	}
	return r.JSON500.Error, true
}

func (r *GetClusterResourceBundlesResp) GetValue() *GetBundlesResponse {
	return r.JSON200
}

func (r *GetClusterResourceBundlesResp) GetHTTPResponse() *http.Response {
	return r.HTTPResponse
}

func (r *DownloadClusterResourceBundleResp) GetBadRequestError() (string, bool) {
	if r.JSON400 == nil {
		return "", false
	}
	return r.JSON400.Error, true
}

func (r *DownloadClusterResourceBundleResp) GetInternalServerError() (string, bool) {
	if r.JSON500 == nil {
		return "", false
	}
	return r.JSON500.Error, true
}

func (r *DownloadClusterResourceBundleResp) GetValue() *DownloadBundleResponse {
	return r.JSON200
}

func (r *DownloadClusterResourceBundleResp) GetHTTPResponse() *http.Response {
	return r.HTTPResponse
}

func (r *ReportClusterResourceBundleStatusResp) GetBadRequestError() (string, bool) {
	if r.JSON400 == nil {
		return "", false
	}
	return r.JSON400.Error, true
}

func (r *ReportClusterResourceBundleStatusResp) GetInternalServerError() (string, bool) {
	if r.JSON500 == nil {
		return "", false
	}
	return r.JSON500.Error, true
}

func (r *ReportClusterResourceBundleStatusResp) GetValue() *EmptyResponse {
	return &EmptyResponse{}
}

func (r *ReportClusterResourceBundleStatusResp) GetHTTPResponse() *http.Response {
	return r.HTTPResponse
}
