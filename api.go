package zerosdk

import (
	"context"
	"fmt"
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

	fetcher, err := cluster_api.NewTokenFetcher(cfg.clusterAPIEndpoint,
		cluster_api.WithHTTPClient(cfg.httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating token fetcher: %w", err)
	}

	tokenCache := token_api.NewCache(fetcher, cfg.apiToken)

	clusterClient, err := cluster_api.NewAuthorizedClient(cfg.clusterAPIEndpoint, tokenCache.GetToken, cfg.httpClient)
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
