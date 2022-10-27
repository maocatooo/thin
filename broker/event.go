package broker

import "context"

type Event interface {
	Publish(ctx context.Context, val interface{}) error
	Subscribe(h interface{}) error
}

type event struct {
	topic string
}

func NewEvent(topic string, h interface{}) (Event, error) {
	err := subscribe(topic, h)
	if err != nil {
		return nil, err
	}
	var e Event
	e = &event{
		topic: topic,
	}
	return e, nil
}

func NewEmptyEvent(topic string) (Event, error) {
	var e Event
	e = &event{
		topic: topic,
	}
	return e, nil
}

func (p *event) Publish(ctx context.Context, val interface{}) error {
	return publish(p.topic, ctx, val)
}

func (p *event) Subscribe(h interface{}) error {
	err := subscribe(p.topic, h)
	return err
}
