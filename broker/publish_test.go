package broker

import (
	"context"
	"testing"
)

func TestPublish(t *testing.T) {
	err := Subscribe("TestPublish", F)
	if err != nil {
		t.Error(err)
	}
	err = Publish("TestPublish", context.Background(), 1)
	if err == nil || err.Error() != `1` {
		t.Error(err)
	}
}
