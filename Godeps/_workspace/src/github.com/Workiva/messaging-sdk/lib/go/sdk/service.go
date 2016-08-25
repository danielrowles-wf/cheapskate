package sdk

// ServiceDescriptor exposes the necessary information for providing and
// connecting to a Harbour service. Eventually, there will be a ResolveService
// function which resolves this information and constructs a Service instance
// using service discovery. Until that service discovery mechanism is in place,
// this must be constructed manually with NewServiceDescriptor.
type ServiceDescriptor interface {
	// Name returns the service name, e.g. "linking".
	Name() string

	// NatsSubject returns the NATS subject this service is listening on for
	// Frugal requests.
	NatsSubject() string

	// FrugalURL returns the full URL hosting a Frugal HTTP server.
	FrugalURL() string

	// FrugalProtocol returns the protocol used for the service, e.g. binary.
	FrugalProtocol() ThriftProtocol

	// TODO: add other fields once those are determined.
}

// serviceDescriptor implements the ServiceDescriptor interface.
type serviceDescriptor struct {
	name           string
	natsSubject    string
	frugalURL      string
	frugalProtocol ThriftProtocol
}

// Name returns the service name, e.g. "linking".
func (s *serviceDescriptor) Name() string {
	return s.name
}

// NatsSubject returns the NATS subject this service is listening on for
// Frugal requests.
func (s *serviceDescriptor) NatsSubject() string {
	return s.natsSubject
}

// FrugalURL returns the full URL hosting a Frugal HTTP server.
func (s *serviceDescriptor) FrugalURL() string {
	return s.frugalURL
}

// FrugalProtocol returns the protocol used for the service, e.g. binary.
func (s *serviceDescriptor) FrugalProtocol() ThriftProtocol {
	return s.frugalProtocol
}
