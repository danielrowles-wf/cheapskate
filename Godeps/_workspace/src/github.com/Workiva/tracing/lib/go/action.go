package tracing

type Endpoint struct {
	Ipv4        string
	Port        int16
	ServiceName string
}

type Action struct {
	Timestamp int64
	Value     string
	Endpoint  Endpoint
}
