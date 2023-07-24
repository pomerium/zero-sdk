package connect

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	host string
	port string
	// requireTLS is whether TLS should be used or cleartext
	requireTLS bool
	// opts are additional options to pass to the gRPC client
	opts []grpc.DialOption
}

// NewConfig returns a new Config from an endpoint string, that has to be in a URL format.
// The endpoint can be either http:// or https:// that will be used to determine whether TLS should be used.
// if port is not specified, it will be inferred from the scheme (80 for http, 443 for https).
func NewConfig(endpoint string) (*Config, error) {
	c := new(Config)
	err := c.parseEndpoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}
	return c, nil
}

// GetConnectionURI returns connection string conforming to https://github.com/grpc/grpc/blob/master/doc/naming.md
func (c *Config) GetConnectionURI() string {
	return fmt.Sprintf("dns:%s:%s", c.host, c.port)
}

func (c *Config) RequireTLS() bool {
	return c.requireTLS
}

func (c *Config) parseEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("error parsing endpoint url: %w", err)
	}

	if u.Path != "" && u.Path != "/" {
		return fmt.Errorf("endpoint path is not supported: %s", u.Path)
	}

	host, port, err := splitHostPort(u.Host)
	if err != nil {
		return fmt.Errorf("error splitting host and port: %w", err)
	}

	requireTLS := true
	var opts []grpc.DialOption
	if u.Scheme == "http" {
		requireTLS = false
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if port == "" {
			port = "80"
		}
	} else if u.Scheme == "https" {
		if port == "" {
			port = "443"
		}
	} else {
		return fmt.Errorf("unsupported url scheme: %s", u.Scheme)
	}

	c.host = host
	c.port = port
	c.requireTLS = requireTLS
	c.opts = append(c.opts, opts...)

	return nil
}

var rePort = regexp.MustCompile(`:(\d+)$`)

func splitHostPort(hostport string) (host, port string, err error) {
	if hostport == "" {
		return "", "", fmt.Errorf("empty hostport")
	}
	if rePort.MatchString(hostport) {
		host, port, err = net.SplitHostPort(hostport)
		if host == "" {
			return "", "", fmt.Errorf("empty host")
		}
		if port == "" {
			return "", "", fmt.Errorf("empty port")
		}
		return host, port, err
	}
	return hostport, "", nil
}
