package broker

import (
	"context"
	"testing"
)

func Test_event_Publish(t *testing.T) {
	e, err := NewEvent("Test_event_Publish", F)
	if err != nil {
		t.Error(err)
	}
	err = e.Publish(context.Background(), 1)
	if err == nil || err.Error() != "1" {
		t.Error(err)
		return
	}
}

func Test_event_Subscribe(t *testing.T) {
	e, err := NewEmptyEvent("Test_event_Publish")
	if err != nil {
		t.Error(err)
	}
	err = e.Subscribe(F)
	if err != nil {
		t.Error(err)
		return
	}
	err = e.Publish(context.Background(), 1)
	if err == nil || err.Error() != "1" {
		t.Error(err)
		return
	}
}
