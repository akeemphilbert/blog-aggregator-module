package blogaggregatormodule_test

import (
	"context"
	"testing"
	"github.com/wepala/weos"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
)

func TestAddBlog (t *testing.T) {
	mockEventRepository := &EventRepositoryMock{
		PersistFunc: func(entity weos.AggregateInterface) error {
			return nil
		},
		AddSubscriberFunc: func(handler weos.EventHandler) {

		},
	}
	mockedApplication := &ApplicationMock{
		EventRepositoryFunc: func () weos.EventRepository {
			return mockEventRepository
		},
	}

	receiver := blogaggregatormodule.NewReceiver(mockedApplication)
	receiver.AddBlog(context.Background(),blogaggregatormodule.AddBlogCommand("https://ak33m.com"))

	//confirm that the event repository persist is called 
	if len(mockEventRepository.PersistCalls()) != 1 {
		t.Errorf("expected events to be persisted")
	}
}