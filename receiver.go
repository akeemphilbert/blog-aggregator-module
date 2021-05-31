package blogaggregatormodule

import (
	"context"
	"encoding/json"
	"github.com/segmentio/ksuid"
	"github.com/wepala/weos"
)

var GenerateID = func () string {
	return ksuid.New().String()
}

type Receiver struct {
	application weos.Application
}

func (r *Receiver) AddBlog(ctx context.Context, command *weos.Command) error {
	var request *AddBlogRequest
	err := json.Unmarshal(command.Payload,&request)
	if err != nil {
		return err
	}
	blog, err := new(Blog).Init(request)
	if err != nil {
		return err
	}
	err = r.application.EventRepository().Persist(blog)
	return err
}

func NewReceiver(application weos.Application) *Receiver {
	return &Receiver{application: application}
}

func Initialize(application weos.Application) error {
	receiver := NewReceiver(application)
	application.Dispatcher().AddSubscriber(AddBlogCommand(""),receiver.AddBlog)
	return nil
}