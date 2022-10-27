package broker

import "context"

/*
向 topic 发送消息
*/
func Publish(topic string, ctx context.Context, val interface{}) error {
	return publish(topic, ctx, val)
}

func publish(topic string, ctx context.Context, val interface{}) error {
	subs := subscribes(topic)
	return subscribersCall(subs, ctx, val)
}
