package blogaggregatormodule_test

import (
	"context"
	"net/http"
	"testing"

	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
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
		HTTPClientFunc: func() *http.Client {
			return &http.Client{}
		},
	}

	receiver := blogaggregatormodule.NewReceiver(mockedApplication)
	receiver.AddBlog(context.Background(),blogaggregatormodule.AddBlogCommand("https://ak33m.com"))

	//confirm that the event repository persist is called 
	if len(mockEventRepository.PersistCalls()) != 1 {
		t.Errorf("expected events to be persisted")
	}
}