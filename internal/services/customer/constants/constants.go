package constants

// SerivceName The name of this module/service
const ServiceName = "customer"

// Dependency Injection Keys
const (
	RegistryKey            = "registry"
	DomainDispatcherKey    = "domainDispatcher"
	MessagePublisherKey    = "messagePublisher"
	MessageSubscriberKey   = "messageSubscriber"
	EventPublisherKey      = "eventPublisher"
	ApplicationKey         = "app"
	DomainEventHandlersKey = "domainEventHandlers"

	CustomerRepoKey = "customerRepo"
)
