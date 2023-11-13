package discovery

type Registry interface {
	RegisterHTTP(
		serviceName string,
		instanceID string,
		healthCheckURL string,
		instanceHost string,
		instancePort int, meta map[string]string,
	) error
	RegisterRPC(
		serviceName string,
		instanceID string,
		healthCheckURL string,
		instanceHost string,
		instancePort int, meta map[string]string,
	) error
	Deregister(instanceID string) error
	ServiceAddresses(serviceName string) ([]string, error)
}
