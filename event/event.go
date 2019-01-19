package event

import "github.com/raidnav/go-cqrs-microservices/schema"

// Event store message API.
type EventStore interface {
	Close()
	PublishMeowCreated(meow schema.Meow) error
	SubscribeMeowCreated() (<-chan MeowCreatedMessage, error)
	OnMeowCreated(f func(message MeowCreatedMessage)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

// Close connection to event store
func Close() {
	impl.Close()
}

// Publish message
func PublishMeowCreated(meow schema.Meow) error {
	return impl.PublishMeowCreated(meow)
}

//Create message when object creation.
func OnMeowCreated(f func(message MeowCreatedMessage)) error {
	return impl.OnMeowCreated(f)
}
