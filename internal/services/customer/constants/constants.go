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
	CommandPublisherKey    = "commandPublisher"
	ReplyPublisherKey      = "replyPublisher"
	InboxStoreKey          = "inboxStore"
	ApplicationKey         = "app"
	DomainEventHandlersKey = "domainEventHandlers"
	CommandHandlersKey     = "commandHandlers"
	ReplyHandlersKey       = "replyHandlers"

	CustomerRepoKey = "customerRepo"
)
