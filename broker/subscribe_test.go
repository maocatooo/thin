package broker

import (
	"context"
	"testing"
)

func TestSubscribe(t *testing.T) {
	topic := "TestSubscribe"
	err := Subscribe(topic, F)
	if err != nil {
		t.Error(err)
	}
	subs := subscribes(topic)
	if len(subs) != 1 {
		t.Error("Subscribe err, the length must 1")
	}
	err = subs[0].call(context.Background(), 1)
	if err == nil || err.Error() != `1` {
		t.Error(err)
	}
}
