package event

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"github.com/raidnav/go-cqrs-microservices/schema"
)

/**
Nats event store structure.
*/

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

/**
Close connection to nats
*/

func (e *NatsEventStore) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
	if e.meowCreatedSubscription != nil {
		err := e.meowCreatedSubscription.Unsubscribe()
		if err != nil {
			panic("Cannot unsubscribe event store, " + err.Error())
		}
	}
	close(e.meowCreatedChan)
}

func (e *NatsEventStore) PublishMeowCreated(meow schema.Meow) error {
	m := MeowCreatedMessage{meow.ID, meow.Body, meow.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.nc.Publish(m.Key(), data)
}

/**
Writes message to nats
*/

func (mq *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

/**
Starts getting message
*/

func (e *NatsEventStore) OnMeowCreated(f func(MeowCreatedMessage)) (err error) {
	m := MeowCreatedMessage{}
	e.meowCreatedSubscription, err = e.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		err := e.readMessage(msg.Data, &m)
		if err != nil {
			panic("Unable to read message")
		}
		f(m)
	})
	return
}

func (mq *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

/**
Create an intermediate channel is created to transform messages into appropriate type.
*/

func (e *NatsEventStore) SubscribeMeowCreate() (<-chan MeowCreatedMessage, error) {
	m := MeowCreatedMessage{}
	e.meowCreatedChan = make(chan MeowCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	e.meowCreatedSubscription, err = e.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				e.readMessage(msg.Data, &m)
				e.meowCreatedChan <- m
			}
		}
	}()
	return (<-chan MeowCreatedMessage)(e.meowCreatedChan), nil
}
