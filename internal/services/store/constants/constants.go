package constants

// SerivceName The name of this module/service
const ServiceName = "store"

// Dependency Injection Keys
const (
	RegistryKey            = "registry"
	DomainDispatcherKey    = "domainDispatcher"
	MessagePublisherKey    = "messagePublisher"
	MessageSubscriberKey   = "messageSubscriber"
	EventPublisherKey      = "eventPublisher"
	AggregateStoreKey      = "aggregateStore"
	ApplicationKey         = "app"
	DomainEventHandlersKey = "domainEventHandlers"

	ProductRepoKey = "productRepo"
)
