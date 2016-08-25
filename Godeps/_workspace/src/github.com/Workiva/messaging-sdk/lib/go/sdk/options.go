package sdk

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/go-vessel"
	"github.com/nats-io/nats"
	"golang.org/x/oauth2/jwt"
)

// ThriftProtocol specifies a serialization protocol used by Frugal.
type ThriftProtocol string

const (
	ProtocolBinary  ThriftProtocol = "binary"
	ProtocolCompact ThriftProtocol = "compact"
	ProtocolJSON    ThriftProtocol = "json"

	defaultHeartbeatInterval   = 5 * time.Second
	defaultConnectTimeout      = 10 * time.Second
	defaultMaxMissedHeartbeats = 3
	defaultTransportBuffer     = 8192
)

// Options contains configuration settings for an sdk client.
type Options struct {
	NATSConfig           nats.Options
	ClientCertFilepath   string
	ClientCACertFilepath string
	ThriftProtocol       ThriftProtocol // DEPRECATED
	ThriftBufferSize     uint           // DEPRECATED
	HeartbeatInterval    time.Duration  // DEPRECATED
	MaxMissedHeartbeats  uint           // DEPRECATED
	ConnectTimeout       time.Duration  // DEPRECATED
	NumWorkers           uint
	VesselHosts          []string      // DEPRECATED
	VesselConfig         vessel.Config // DEPRECATED
	VesselAuth           *jwt.Config   // DEPRECATED
	IamTokenFetcher      IamTokenFetcher
}

// NewOptions creates a new Options for the given service client id.
// DEPRECATED. Will be removed in 2.0.0. Use NewClientOptions instead.
func NewOptions(clientID string) Options {
	cfg := nats.DefaultOptions
	servers := MsgURL()
	if len(servers) > 0 {
		// If MSG_URL set, override any URL set in the config
		cfg.Url = ""
		cfg.Servers = servers
	} else if cfg.Url == "" {
		// If nothing is set, take the default url
		cfg.Url = nats.DefaultURL
	}
	cfg.Name = fmt.Sprintf("%s-%d", clientID, os.Getpid())

	return Options{
		NATSConfig:           cfg,
		ClientCertFilepath:   MsgCert(),
		ClientCACertFilepath: MsgCACert(),
		ThriftProtocol:       ProtocolBinary,
		ThriftBufferSize:     defaultTransportBuffer,
		HeartbeatInterval:    defaultHeartbeatInterval,
		MaxMissedHeartbeats:  defaultMaxMissedHeartbeats,
		ConnectTimeout:       defaultConnectTimeout,
		NumWorkers:           uint(runtime.NumCPU()),
		VesselConfig:         *vessel.NewConfig(clientID),
	}
}

// NewClientOptions creates a new Options used to create a Messaging SDK
// Client.
func NewClientOptions() Options {
	cfg := nats.DefaultOptions
	servers := MsgURL()
	if len(servers) > 0 {
		// If MSG_URL set, override any URL set in the config
		cfg.Url = ""
		cfg.Servers = servers
	} else if cfg.Url == "" {
		// If nothing is set, take the default url
		cfg.Url = nats.DefaultURL
	}
	hostname, err := os.Hostname()
	if err != nil {
		// Not really sure what to do if this call fails.
		hostname = "hostname"
	}
	cfg.Name = fmt.Sprintf("%s-%d", hostname, os.Getpid())

	return Options{
		NATSConfig:           cfg,
		ClientCertFilepath:   MsgCert(),
		ClientCACertFilepath: MsgCACert(),
		ThriftProtocol:       ProtocolBinary,
		ThriftBufferSize:     defaultTransportBuffer,
		NumWorkers:           uint(runtime.NumCPU()),
	}
}

func (o Options) newThriftProtocolFactory() (thrift.TProtocolFactory, error) {
	return newTProtocolFactory(o.ThriftProtocol)
}

func (o Options) newThriftTransportFactory() thrift.TTransportFactory {
	var transportFactory thrift.TTransportFactory
	if int(o.ThriftBufferSize) > 0 {
		transportFactory = thrift.NewTBufferedTransportFactory(int(o.ThriftBufferSize))
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}
	return transportFactory
}

func newTProtocolFactory(protocol ThriftProtocol) (thrift.TProtocolFactory, error) {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case ProtocolBinary:
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	case ProtocolCompact:
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case ProtocolJSON:
		protocolFactory = thrift.NewTJSONProtocolFactory()
	default:
		return nil, fmt.Errorf("sdk: invalid protocol specified: %s", protocol)
	}
	return protocolFactory, nil
}
