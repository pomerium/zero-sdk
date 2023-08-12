package zerosdk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pomerium/zero-sdk/apierror"
	cluster_api "github.com/pomerium/zero-sdk/cluster"
	connect_api "github.com/pomerium/zero-sdk/connect"
	connect_mux "github.com/pomerium/zero-sdk/connect-mux"
	token_api "github.com/pomerium/zero-sdk/token"
)

// API is a Pomerium Zero Cluster API client
type API struct {
	cfg              *config
	cluster          cluster_api.ClientWithResponsesInterface
	downloadURLCache *cluster_api.URLCache
	tokenProvider    func(ctx context.Context, minTTL time.Duration) (string, error)
}

// NewAPI creates a new API client
func NewAPI(opts ...Option) (*API, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	fetcher, err := cluster_api.NewTokenFetcher(cfg.clusterAPIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating token fetcher: %w", err)
	}

	tokenCache := token_api.NewCache(fetcher, cfg.apiToken)

	clusterClient, err := cluster_api.NewAuthorizedClient(cfg.clusterAPIEndpoint, tokenCache.GetToken, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("error creating cluster client: %w", err)
	}

	return &API{
		cfg:              cfg,
		cluster:          clusterClient,
		downloadURLCache: cluster_api.NewURLCache(),
		tokenProvider:    tokenCache.GetToken,
	}, nil
}

// Connect creates a new connect mux client and starts it
func (api *API) Connect(ctx context.Context) (*connect_mux.Mux, error) {
	client, err := connect_api.NewAuthorizedConnectClient(ctx, api.cfg.connectAPIEndpoint, api.tokenProvider)
	if err != nil {
		return nil, fmt.Errorf("error creating connect client: %w", err)
	}

	return connect_mux.Start(ctx, client), nil
}

// GetClusterBootstrapConfig fetches the bootstrap configuration from the cluster API
func (api *API) GetClusterBootstrapConfig(ctx context.Context) (*cluster_api.BootstrapConfig, error) {
	return apierror.CheckResponse[cluster_api.BootstrapConfig](
		api.cluster.GetClusterBootstrapConfigWithResponse(ctx),
	)
}

// GetClusterResourceBundles fetches the resource bundles from the cluster API
func (api *API) GetClusterResourceBundles(ctx context.Context) (*cluster_api.GetBundlesResponse, error) {
	return apierror.CheckResponse[cluster_api.GetBundlesResponse](
		api.cluster.GetClusterResourceBundlesWithResponse(ctx),
	)
}

// DownloadClusterResourceBundle obtains a download URL for a resource bundle from the cluster API
func (api *API) DownloadClusterResourceBundle(ctx context.Context, id string, minTTL time.Duration) (*url.URL, error) {
	u, ok := api.downloadURLCache.Get(id, minTTL)
	if ok {
		return &u, nil
	}

	return api.updateBundleDownloadURL(ctx, id)
}

func (api *API) updateBundleDownloadURL(ctx context.Context, id string) (*url.URL, error) {
	now := time.Now()

	resp, err := apierror.CheckResponse[cluster_api.DownloadBundleResponse](
		api.cluster.DownloadClusterResourceBundleWithResponse(ctx, id),
	)
	if err != nil {
		return nil, fmt.Errorf("get bundle download URL: %w", err)
	}

	expiresSeconds, err := strconv.ParseInt(resp.ExpiresInSeconds, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse expiration: %w", err)
	}

	u, err := url.Parse(resp.Url)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	api.downloadURLCache.Set(id, *u, now.Add(time.Duration(expiresSeconds)*time.Second))
	return u, nil
}
