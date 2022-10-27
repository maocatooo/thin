package broker

import "context"

/*
消息结构体接口,Publish 发布消息,Subscribe 接受消息的结构
*/
type Event interface {
	Publish(ctx context.Context, val interface{}) error
	Subscribe(h interface{}) error
}

type event struct {
	topic string
}

/*
创建一个Event,并且设置一个接受消息的结构
*/
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

/*
创建一个空Event
*/
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
