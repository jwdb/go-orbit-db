package pubsub

import (
	"context"
	"fmt"
	"github.com/berty/go-orbit-db/ipfs"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

type pubSub struct {
	ipfs          ipfs.Services
	id            peer.ID
	subscriptions map[string]Subscription
}

func NewPubSub(is ipfs.Services, id peer.ID) (Interface, error) {
	if is == nil {
		return nil, errors.New("ipfs services are not defined")
	}

	ps := is.PubSub()

	if ps == nil {
		return nil, errors.New("pubsub service is not provided by the current ipfs instance")
	}

	return &pubSub{
		ipfs:          is,
		id:            id,
		subscriptions: map[string]Subscription{},
	}, nil
}

func (p *pubSub) Subscribe(ctx context.Context, topic string) (Subscription, error) {
	sub, ok := p.subscriptions[topic]
	if ok {
		return sub, nil
	}

	logger().Debug(fmt.Sprintf("starting pubsub listener for peer %s on topic %s", p.id, topic))

	s, err := NewSubscription(ctx, p.ipfs, topic)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create new pubsub subscription")
	}

	p.subscriptions[topic] = s

	return s, nil
}

func (p *pubSub) Publish(ctx context.Context, topic string, message []byte) error {
	if _, ok := p.subscriptions[topic]; !ok {
		return errors.New("to subscribed to this topic")
	}

	return p.ipfs.PubSub().Publish(ctx, topic, message)
}

func (p *pubSub) Close() error {
	for _, sub := range p.subscriptions {
		_ = sub.Close()
	}

	return nil
}

func (p *pubSub) Unsubscribe(topic string) error {
	s, ok := p.subscriptions[topic]
	if !ok {
		return errors.New("no subscription found")
	}

	_ = s.Close()

	return nil
}

var _ Interface = &pubSub{}
