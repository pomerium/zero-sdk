package zerosdk

import "fmt"

type Option func(*config)

type config struct {
	clusterAPIEndpoint string
	connectAPIEndpoint string
	apiToken           string
}

// WithClusterAPIEndpoint sets the cluster API endpoint
func WithClusterAPIEndpoint(endpoint string) Option {
	return func(cfg *config) {
		cfg.clusterAPIEndpoint = endpoint
	}
}

// WithConnectAPIEndpoint sets the connect API endpoint
func WithConnectAPIEndpoint(endpoint string) Option {
	return func(cfg *config) {
		cfg.connectAPIEndpoint = endpoint
	}
}

// WithAPIToken sets the API token
func WithAPIToken(token string) Option {
	return func(cfg *config) {
		cfg.apiToken = token
	}
}

func newConfig(opts ...Option) (*config, error) {
	cfg := new(config)
	for _, opt := range opts {
		opt(cfg)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *config) validate() error {
	if c.clusterAPIEndpoint == "" {
		return fmt.Errorf("cluster API endpoint is required")
	}
	if c.connectAPIEndpoint == "" {
		return fmt.Errorf("connect API endpoint is required")
	}
	if c.apiToken == "" {
		return fmt.Errorf("API token is required")
	}
	return nil
}
