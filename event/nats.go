package event

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"github.com/raidnav/go-cqrs-microservices/schema"
)

// Nats event store structure.
type NatsEventStore struct {
	nc                      *nats.Conn
	meowCreatedSubscription *nats.Subscription
	meowCreatedChan         chan MeowCreatedMessage
}

func newNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

// Close connection to nats
func (natsEventStore *NatsEventStore) Close() {
	if natsEventStore.nc != nil {
		natsEventStore.nc.Close()
	}
	if natsEventStore.meowCreatedSubscription != nil {
		err := natsEventStore.meowCreatedSubscription.Unsubscribe()
		if err != nil {
			panic("Cannot unsubscribe event store, " + err.Error())
		}
	}
	close(natsEventStore.meowCreatedChan)
}

func (natsEventStore *NatsEventStore) PublishMeowCreated(meow schema.Meow) error {
	m := MeowCreatedMessage{meow.ID, meow.Body, meow.CreatedAt}
	data, err := natsEventStore.writeMessage(&m)
	if err != nil {
		return err
	}
	return natsEventStore.nc.Publish(m.Key(), data)
}

// Writes message to nats
func (natsEventStore *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Starts getting message
func (natsEventStore *NatsEventStore) OnMeowCreated(f func(MeowCreatedMessage)) (err error) {
	m := MeowCreatedMessage{}
	natsEventStore.meowCreatedSubscription, err = natsEventStore.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		err := natsEventStore.readMessage(msg.Data, &m)
		if err != nil {
			panic("Unable to read message")
		}
		f(m)
	})
	return
}

func (natsEventStore *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

// Create an intermediate channel is created to transform messages into appropriate type.
func (natsEventStore *NatsEventStore) SubscribeMeowCreate() (<-chan MeowCreatedMessage, error) {
	m := MeowCreatedMessage{}
	natsEventStore.meowCreatedChan = make(chan MeowCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	natsEventStore.meowCreatedSubscription, err = natsEventStore.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				err := natsEventStore.readMessage(msg.Data, &m)

				if err != nil {
					panic("Unable to read message")
				}

				natsEventStore.meowCreatedChan <- m
			}
		}
	}()
	return (<-chan MeowCreatedMessage)(natsEventStore.meowCreatedChan), nil
}
