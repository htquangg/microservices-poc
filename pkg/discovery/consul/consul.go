package consul

import (
	"fmt"
	"strconv"
	"sync"

	consul_sd_kit "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/htquangg/microservices-poc/pkg/discovery"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

var _ discovery.Registry = (*consulClient)(nil)

type consulClient struct {
	cfg         *Config
	mutex       sync.Mutex
	instanceMap sync.Map
	log         logger.Logger
	client      consul_sd_kit.Client
}

func New(cfg *Config, log logger.Logger) (*consulClient, error) {
	dc := &consulClient{
		cfg: cfg,
		log: log,
	}
	if err := dc.initClient(); err != nil {
		return nil, err
	}

	return dc, nil
}

func (dc *consulClient) initClient() (err error) {
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = dc.cfg.Address()

	var apiClient *consul.Client

	if apiClient, err = consul.NewClient(consulConfig); err != nil {
		return err
	}

	dc.client = consul_sd_kit.NewClient(apiClient)

	return nil
}

func (dc *consulClient) RegisterHTTP(
	serviceName,
	instanceID,
	healthCheckURL string,
	instanceHost string,
	instancePort int, meta map[string]string,
) error {
	serviceRegistration := &consul.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Check: &consul.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "120s",
			HTTP: fmt.Sprintf(
				"http://%s:%s%s",
				instanceHost,
				strconv.Itoa(instancePort),
				healthCheckURL,
			),
			Interval: "10s",
			Notes:    "Basic health checks",
		},
	}

	if err := dc.client.Register(serviceRegistration); err != nil {
		return err
	}

	dc.log.Infof(
		"register consul service success with id and name: %s -- %s",
		serviceRegistration.ID,
		serviceRegistration.Name,
	)

	return nil
}

func (dc *consulClient) RegisterRPC(
	serviceName,
	instanceID,
	healthCheckURL string,
	instanceHost string,
	instancePort int, meta map[string]string,
) error {
	serviceRegistration := &consul.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Check: &consul.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "120s",
			GRPC: fmt.Sprintf(
				"%s:%s/grpc.health.v1.%v",
				instanceHost,
				strconv.Itoa(instancePort),
				healthCheckURL,
			),
			Interval: "10s",
			Notes:    "Basic health checks",
		},
	}

	if err := dc.client.Register(serviceRegistration); err != nil {
		return err
	}

	dc.log.Infof(
		"register consul service success with id and name: %s -- %s",
		serviceRegistration.ID,
		serviceRegistration.Name,
	)

	return nil
}

func (dc *consulClient) Deregister(instanceID string) error {
	serviceRegistration := &consul.AgentServiceRegistration{
		ID: instanceID,
	}

	if err := dc.client.Deregister(serviceRegistration); err != nil {
		return err
	}

	dc.log.Infof("deregister consul service success with id: %s", instanceID)

	return nil
}

func (dc *consulClient) ServiceAddresses(
	serviceName string,
) ([]string, error) {
	instanceList, ok := dc.instanceMap.Load(serviceName)
	if ok {
		return instanceList.([]string), nil
	}

	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	if instanceList, ok = dc.instanceMap.Load(serviceName); ok {
		return instanceList.([]string), nil
	}

	go func() {
		params := make(map[string]interface{})
		params["type"] = "service"
		params["service"] = serviceName
		plan, _ := watch.Parse(params)
		plan.Handler = func(_ uint64, i interface{}) {
			if i == nil {
				return
			}

			v, ok := i.([]*consul.ServiceEntry)
			if !ok {
				return
			}

			if len(v) == 0 {
				dc.instanceMap.Store(serviceName, []interface{}{})
			}

			var healthServices []interface{}

			for _, service := range v {
				if service.Checks.AggregatedStatus() == consul.HealthPassing {
					healthServices = append(healthServices, service)
				}
			}

			dc.instanceMap.Store(serviceName, healthServices)
		}

		defer plan.Stop()
		plan.Run(dc.cfg.Address())
	}()

	entries, _, err := dc.client.Service(serviceName, "", false, nil)
	if err != nil {
		dc.instanceMap.Store(serviceName, []interface{}{})
		dc.log.Error("Discover Service Error")

		return nil, err
	}

	instances := make([]string, 0, len(entries))

	for _, instance := range entries {
		instances = append(
			instances,
			fmt.Sprintf(
				"%s:%d",
				instance.Service.Address,
				instance.Service.Port,
			),
		)
	}

	dc.instanceMap.Store(serviceName, instances)

	return instances, nil
}
