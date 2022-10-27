package broker

import (
	"context"
	"sync"
)

/*
处理 topic 对应的消息
*/
func Subscribe(topic string, h interface{}) error {
	return subscribe(topic, h)
}

func subscribe(topic string, h interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	subs, err := methodObjToHandler(h)
	if err != nil {
		return err
	}
	subscribers[topic] = append(subscribers[topic], subs...)
	return nil
}

var (
	subscribers map[string][]Subscriber
	lock        sync.RWMutex
)

func init() {
	subscribers = map[string][]Subscriber{}
}

func subscribes(topic string) []Subscriber {
	lock.RLock()
	defer lock.RUnlock()
	return subscribers[topic]
}

func subscribersCall(subs []Subscriber, ctx context.Context, val interface{}) error {
	v := valueCopy(val)
	for _, sub := range subs {
		err := sub.call(ctx, v)
		if err != nil {
			return err
		}
	}
	return nil
}
