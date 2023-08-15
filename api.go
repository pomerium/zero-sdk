package zerosdk

import (
	"context"
	"fmt"

	"github.com/pomerium/zero-sdk/apierror"
	cluster_api "github.com/pomerium/zero-sdk/cluster"
	connect_api "github.com/pomerium/zero-sdk/connect"
	connect_mux "github.com/pomerium/zero-sdk/connect-mux"
	"github.com/pomerium/zero-sdk/fanout"
	token_api "github.com/pomerium/zero-sdk/token"
)

// API is a Pomerium Zero Cluster API client
type API struct {
	cfg              *config
	cluster          cluster_api.ClientWithResponsesInterface
	mux              *connect_mux.Mux
	downloadURLCache *cluster_api.URLCache
}

// WatchOption defines which events to watch for
type WatchOption = connect_mux.WatchOption

// NewAPI creates a new API client
func NewAPI(ctx context.Context, opts ...Option) (*API, error) {
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

	connectClient, err := connect_api.NewAuthorizedConnectClient(ctx, cfg.connectAPIEndpoint, tokenCache.GetToken)
	if err != nil {
		return nil, fmt.Errorf("error creating connect client: %w", err)
	}

	return &API{
		cfg:              cfg,
		cluster:          clusterClient,
		mux:              connect_mux.New(connectClient),
		downloadURLCache: cluster_api.NewURLCache(),
	}, nil
}

// Connect connects to the connect API and allows watching for changes
func (api *API) Connect(ctx context.Context, opts ...fanout.Option) error {
	return api.mux.Run(ctx, opts...)
}

// Watch dispatches API updates
func (api *API) Watch(ctx context.Context, opts ...WatchOption) error {
	return api.mux.Watch(ctx, opts...)
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
