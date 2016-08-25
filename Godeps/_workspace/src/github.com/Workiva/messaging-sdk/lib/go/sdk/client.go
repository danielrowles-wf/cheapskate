package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Workiva/frugal/lib/go"
	"github.com/Workiva/go-vessel"
	"github.com/nats-io/nats"
	"golang.org/x/oauth2"
)

var (
	ErrIsOpen    = errors.New("frugal: client is already open")
	ErrNotOpen   = errors.New("frugal: client is not open")
	ErrNoServers = errors.New("frugal: no servers available for connection")
)

const (
	minHeartbeatInterval = 500 * time.Millisecond
	tlsPrefix            = "tls://"
	rpcQueueGroup        = "rpc"
)

// ConnectionHandler is a callback invoked for different connection events.
type ConnectionHandler func()

// ErrorHandler is a callback invoked when there is an asynchronous error
// processing inbound messages.
type ErrorHandler func(err error)

// MessagingClient is the top-level client API for the Messaging SDK.
type MessagingClient interface {

	// SetClosedHandler sets the handler invoked when the client has been
	// closed and will not attempt to reopen. The client must be manually
	// reopened at this point.
	SetClosedHandler(handler ConnectionHandler)

	// SetDisconnectHandler sets the handler invoked when the client has
	// disconnected but might be attempting to reconnect. The closed handler
	// will be invoked if reconnect fails. The reconnect handler will be
	// invoked if reconnect succeeds.
	SetDisconnectHandler(handler ConnectionHandler)

	// SetReconnectHandler sets the handler invoked when the client has
	// successfully reconnected.
	SetReconnectHandler(handler ConnectionHandler)

	// SetErrorHandler sets the handler invoked when there is an asynchronous
	// error processing inbound messages.
	SetErrorHandler(handler ErrorHandler)

	// Open prepares the Client for use. Returns ErrIsOpen if the Client is
	// already open and ErrNoServers if connecting to NATS times out.
	Open() error

	// IsOpen indicates if the Client has been opened for use.
	IsOpen() bool

	// Close disconnects the Client.
	Close() error

	// GetIamToken calls the IamTokenFetcher in the client options and returns a
	// JWT token.
	GetIamToken() (string, error)

	// NewClient returns the plumbing required for a Frugal client of the given
	// service. You must explicitly open the transport before using it with a
	// Frugal client.
	NewClient(service ServiceDescriptor) (frugal.FTransport, *frugal.FProtocolFactory, error)

	// NewServer returns a Frugal server for the given service.
	NewServer(service ServiceDescriptor, processor frugal.FProcessor) (frugal.FServer, error)

	// NewHttpHandlerFunc returns a http.HandlerFunc able to serve frugal
	// requests for the given processor.
	NewHttpHandlerFunc(service ServiceDescriptor, processor frugal.FProcessor) (http.HandlerFunc, error)

	// NewHttpClient returns the plumbing required for a Frugal client of the
	// given service running an HTTP handler. You must explicitly open the
	// transport before using it with a Frugal client.
	NewHttpClient(client *http.Client, service ServiceDescriptor) (frugal.FTransport, *frugal.FProtocolFactory, error)

	// NewPubSubProvider returns the plumbing required for Frugal pub/sub using
	// the given protocol.
	NewPubSubProvider(proto ThriftProtocol) (*frugal.FScopeProvider, error)

	// NewServiceDescriptor creates a ServiceDescriptor instance with the given
	// parameters. This will eventually be replaced by ResolveService which will
	// resolve this information via service discovery.
	NewServiceDescriptor(name, natsSubject, frugalURL string, proto ThriftProtocol) ServiceDescriptor

	// ResolveService resolves the given service and creates a ServiceDescriptor
	// for it, returning an error if it failed to resolve.
	ResolveService(name string) (ServiceDescriptor, error)
}

// Client provides access to the messaging platform. It implements the
// MessagingClient interface.
type Client struct {
	options Options
	conn    *nats.Conn
}

// New returns a new Client. Open must be called before the Client can be used.
// DEPRECATED. Will be removed in 2.0.0. Use NewMessagingClient instead.
func New(options Options) *Client {
	return &Client{options: options}
}

// NewMessagingClient creates a new MessagingClient which provides the API for
// interacting with the Messaging Platform.
func NewMessagingClient(options Options) MessagingClient {
	return &Client{options: options}
}

// SetClosedHandler sets the handler invoked when the client has been closed
// and will not attempt to reopen. The client must be manually reopened at this
// point.
func (c *Client) SetClosedHandler(handler ConnectionHandler) {
	c.conn.SetClosedHandler(func(_ *nats.Conn) {
		handler()
	})
}

// SetDisconnectHandler sets the handler invoked when the client has
// disconnected but might be attempting to reconnect. The closed handler will
// be invoked if reconnect fails. The reconnect handler will be invoked if
// reconnect succeeds.
func (c *Client) SetDisconnectHandler(handler ConnectionHandler) {
	c.conn.SetDisconnectHandler(func(_ *nats.Conn) {
		handler()
	})
}

// SetReconnectHandler sets the handler invoked when the client has
// successfully reconnected.
func (c *Client) SetReconnectHandler(handler ConnectionHandler) {
	c.conn.SetReconnectHandler(func(_ *nats.Conn) {
		handler()
	})
}

// SetErrorHandler sets the handler invoked when there is an asynchronous error
// processing inbound messages.
func (c *Client) SetErrorHandler(handler ErrorHandler) {
	c.conn.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		handler(err)
	})
}

// Open prepares the Client for use. Returns ErrIsOpen if the Client is
// already open and ErrNoServers if connecting to NATS times out.
func (c *Client) Open() error {
	if c.IsOpen() {
		return ErrIsOpen
	}

	// Chose a cluster representative to check for tls enabled
	clusterRep := c.options.NATSConfig.Url
	if len(c.options.NATSConfig.Servers) > 0 {
		clusterRep = c.options.NATSConfig.Servers[0]
	}

	// Set TLS credentials
	if strings.HasPrefix(clusterRep, tlsPrefix) {
		cert, err := tls.LoadX509KeyPair(c.options.ClientCertFilepath, c.options.ClientCertFilepath)
		if err != nil {
			return fmt.Errorf("messaging_sdk: Error loading client tls cert: %s", err)
		}

		cacert, err := ioutil.ReadFile(c.options.ClientCACertFilepath)
		if err != nil {
			return fmt.Errorf("messaging_sdk: Error loading client CA tls cert: %s", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(cacert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			MinVersion:   tls.VersionTLS12,
		}
		tlsConfig.BuildNameToCertificate()

		c.options.NATSConfig.TLSConfig = tlsConfig
	}
	conn, err := c.options.NATSConfig.Connect()
	if err != nil {
		if err == nats.ErrNoServers {
			return ErrNoServers
		}
		return err
	}
	c.conn = conn
	return nil
}

// IsOpen indicates if the Client has been opened for use.
func (c *Client) IsOpen() bool {
	return c.conn != nil && c.conn.Status() == nats.CONNECTED
}

// Close disconnects the Client.
func (c *Client) Close() error {
	if !c.IsOpen() {
		return ErrNotOpen
	}
	c.conn.Close()
	c.conn = nil
	return nil
}

// GetIamToken calls the IamTokenFetcher in the client options and returns a
// JWT token.
func (c *Client) GetIamToken() (string, error) {
	if c.options.IamTokenFetcher != nil {
		return c.options.IamTokenFetcher.GetIamToken()
	}
	return "", errors.New("messaging_sdk: No token fetcher present")
}

// Vessel returns a new Vessel client.
// DEPRECATED. Will be removed in 2.0.0. Use Frugal pub/sub instead.
func (c *Client) Vessel() (vessel.Vessel, error) {
	if c.options.VesselAuth != nil {
		c.options.VesselConfig.HTTPClient = c.options.VesselAuth.Client(oauth2.NoContext)
	}
	return vessel.New(c.options.VesselHosts, &c.options.VesselConfig)
}

// ProvideClient returns the plumbing required for a Frugal client of the
// given service. You must explicitly open the transport before using it with
// a Frugal client.
// DEPRECATED. Will be removed in 2.0.0. Use NewClient instead.
func (c *Client) ProvideClient(service Service) (frugal.FTransport, *frugal.FProtocolFactory, error) {
	if !c.IsOpen() {
		return nil, nil, ErrNotOpen
	}

	tProtocolFactory, err := c.options.newThriftProtocolFactory()
	if err != nil {
		return nil, nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	transportFactory := c.options.newThriftTransportFactory()
	transport := frugal.NewNatsServiceTTransport(c.conn, string(service),
		c.options.ConnectTimeout, c.options.MaxMissedHeartbeats)
	tr := transportFactory.GetTransport(transport)
	fTransport := frugal.NewFMuxTransport(tr, c.options.NumWorkers)
	return fTransport, fProtocolFactory, nil
}

// NewClient returns the plumbing required for a Frugal client of the given
// service. You must explicitly open the transport before using it with a
// Frugal client.
func (c *Client) NewClient(service ServiceDescriptor) (frugal.FTransport, *frugal.FProtocolFactory, error) {
	if !c.IsOpen() {
		return nil, nil, ErrNotOpen
	}

	tProtocolFactory, err := newTProtocolFactory(service.FrugalProtocol())
	if err != nil {
		return nil, nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	transport := frugal.NewFNatsTransport(c.conn, service.NatsSubject(), "")
	return transport, fProtocolFactory, nil
}

// ProvideServer returns a Frugal server for the given service.
// DEPRECATED. Will be removed in 2.0.0. Use NewServer instead.
func (c *Client) ProvideServer(service Service, processor frugal.FProcessor) (frugal.FServer, error) {
	if !c.IsOpen() {
		return nil, ErrNotOpen
	}

	tProtocolFactory, err := c.options.newThriftProtocolFactory()
	if err != nil {
		return nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	fTransportFactory := frugal.NewFMuxTransportFactory(c.options.NumWorkers)
	heartbeatInterval := c.options.HeartbeatInterval
	if heartbeatInterval < minHeartbeatInterval {
		heartbeatInterval = minHeartbeatInterval
	}
	return frugal.NewFNatsServerFactory(c.conn, string(service),
		heartbeatInterval, c.options.MaxMissedHeartbeats,
		frugal.NewFProcessorFactory(processor), fTransportFactory, fProtocolFactory), nil
}

// NewServer returns a Frugal server for the given service.
func (c *Client) NewServer(service ServiceDescriptor, processor frugal.FProcessor) (frugal.FServer, error) {
	if !c.IsOpen() {
		return nil, ErrNotOpen
	}

	tProtocolFactory, err := newTProtocolFactory(service.FrugalProtocol())
	if err != nil {
		return nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	return frugal.NewFNatsServerWithStatelessConfig(
		c.conn,
		[]string{service.NatsSubject()},
		c.options.NumWorkers,
		64,
		c.options.HeartbeatInterval,
		c.options.MaxMissedHeartbeats,
		frugal.NewFProcessorFactory(processor),
		frugal.NewAdapterTransportFactory(),
		fProtocolFactory,
	), nil
}

// ProvideHttpHandlerFunc returns a http.HandlerFunc able to serve frugal
// requests for the given processor.
// DEPRECATED. Will be removed in 2.0.0. Use NewHttpHandlerFunc instead.
func (c *Client) ProvideHttpHandlerFunc(processor frugal.FProcessor) (http.HandlerFunc, error) {
	tProtocolFactory, err := c.options.newThriftProtocolFactory()
	if err != nil {
		return nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	return frugal.NewFrugalHandlerFunc(processor, fProtocolFactory, fProtocolFactory), nil
}

// NewHttpHandlerFunc returns a http.HandlerFunc able to serve frugal requests
// for the given processor.
func (c *Client) NewHttpHandlerFunc(service ServiceDescriptor, processor frugal.FProcessor) (http.HandlerFunc, error) {
	tProtocolFactory, err := newTProtocolFactory(service.FrugalProtocol())
	if err != nil {
		return nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	return frugal.NewFrugalHandlerFunc(processor, fProtocolFactory, fProtocolFactory), nil
}

// NewHttpClient returns the plumbing required for a Frugal client of the given
// service running an HTTP handler. You must explicitly open the transport
// before using it with a Frugal client.
func (c *Client) NewHttpClient(client *http.Client, service ServiceDescriptor) (frugal.FTransport, *frugal.FProtocolFactory, error) {
	if client == nil {
		client = http.DefaultClient
	}
	tProtocolFactory, err := newTProtocolFactory(service.FrugalProtocol())
	if err != nil {
		return nil, nil, err
	}
	fProtocolFactory := frugal.NewFProtocolFactory(tProtocolFactory)
	transport := frugal.NewHttpFTransportBuilder(client, service.FrugalURL()).Build()
	return transport, fProtocolFactory, nil
}

// ProvidePubSub returns the plumbing required for Frugal pub/sub.
// DEPRECATED. Will be removed in 2.0.0. Use NewPubSubProvider instead.
func (c *Client) ProvidePubSub() (*frugal.FScopeProvider, error) {
	if !c.IsOpen() {
		return nil, ErrNotOpen
	}

	protocolFactory, err := c.options.newThriftProtocolFactory()
	if err != nil {
		return nil, err
	}

	natsFactory := frugal.NewFNatsScopeTransportFactory(c.conn)
	return frugal.NewFScopeProvider(natsFactory, frugal.NewFProtocolFactory(protocolFactory)), nil
}

// NewPubSubProvider returns the plumbing required for Frugal pub/sub using the
// given protocol.
func (c *Client) NewPubSubProvider(proto ThriftProtocol) (*frugal.FScopeProvider, error) {
	if !c.IsOpen() {
		return nil, ErrNotOpen
	}

	protocolFactory, err := newTProtocolFactory(proto)
	if err != nil {
		return nil, err
	}

	natsFactory := frugal.NewFNatsScopeTransportFactory(c.conn)
	return frugal.NewFScopeProvider(natsFactory, frugal.NewFProtocolFactory(protocolFactory)), nil
}

// NewServiceDescriptor creates a ServiceDescriptor instance with the given
// parameters. This will eventually be replaced by ResolveService which will
// resolve this information via service discovery.
func (c *Client) NewServiceDescriptor(name, natsSubject, frugalURL string, proto ThriftProtocol) ServiceDescriptor {
	return &serviceDescriptor{
		name:           name,
		natsSubject:    natsSubject,
		frugalURL:      frugalURL,
		frugalProtocol: proto,
	}
}

// ResolveService resolves the given service and creates a ServiceDescriptor
// for it, returning an error if it failed to resolve.
func (c *Client) ResolveService(name string) (ServiceDescriptor, error) {
	// TODO: Implement once service discovery proxy is available.
	return nil, errors.New("Not implemented: use NewServiceDescriptor instead")
}
