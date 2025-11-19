package api

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/grpc/resolver"
)

// mcrunnerBuilder implements resolver.Builder for the "mcrunner" scheme.
type mcrunnerBuilder struct{}

func (b *mcrunnerBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// Parse the target in a forgiving way. The user may supply:
	//  - "mcrunner://host:port"
	//  - "mcrunner:///host:port"
	// Prefer an explicit endpoint, then URL.Host, then full Path.
	var endpoint string
	if target.Endpoint() != "" {
		endpoint = target.Endpoint()
	} else if target.URL.Host != "" {
		endpoint = target.URL.Host
	} else if target.URL.Path != "" {
		// Path may include a leading '/'
		endpoint = strings.TrimPrefix(target.URL.Path, "/")
	}
	if endpoint == "" {
		return nil, errors.New("mcrunner resolver needs an endpoint (host[:port])")
	}
	// If there's no explicit port, default to 50051
	if !strings.Contains(endpoint, ":") {
		endpoint = fmt.Sprintf("%s:%d", endpoint, 50051)
	}

	r := &mcrunnerResolver{
		target: target,
		cc:     cc,
		addr:   endpoint,
	}
	r.resolveNow() // Initial resolution.
	return r, nil
}

func (b *mcrunnerBuilder) Scheme() string {
	return "mcrunner"
}

// mcrunnerResolver implements resolver.Resolver.
type mcrunnerResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	addr   string
}

func (r *mcrunnerResolver) ResolveNow(resolver.ResolveNowOptions) {
	r.resolveNow()
}

func (r *mcrunnerResolver) Close() {
	// Cleanup logic if needed (e.g., stop watchers).
}

func (r *mcrunnerResolver) resolveNow() {
	// Validate a little by parsing as URL (allowing host:port without scheme)
	if _, err := url.Parse(fmt.Sprintf("//%s", r.addr)); err != nil {
		r.cc.ReportError(fmt.Errorf("invalid address %q: %w", r.addr, err))
		return
	}

	state := resolver.State{
		Addresses: []resolver.Address{{Addr: r.addr}},
	}
	if err := r.cc.UpdateState(state); err != nil {
		r.cc.ReportError(err)
	}
}
