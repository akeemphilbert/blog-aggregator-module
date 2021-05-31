package blogaggregatormodule

import (
	"encoding/json"

	"github.com/wepala/weos"
)

type Blog struct {
	weos.AggregateRoot
	Title string
	Description string
	URL string `json:"url"`
}

func (b *Blog) Init(blog *AddBlogRequest) (*Blog, error) {
	err := b.Validate(blog)
	if err != nil {
		return nil, err
	}
	event, err := weos.NewBasicEvent(BLOG_ADDED,GenerateID(),"Blog",blog)
	b.NewChange(event)
	b.ApplyChanges([]*weos.Event{event})
	return b, nil
}

func (b *Blog) Validate(blog *AddBlogRequest) error {
	if blog.Url == "" {
		return weos.NewDomainError("a blog url must be specified","Blog","",nil)
	}
	return nil
}

func (b *Blog) ApplyChanges (changes []*weos.Event) error {
	for _, change := range changes {
		switch change.Type {
		case BLOG_ADDED:
			err := json.Unmarshal(change.Payload, b)
			if err != nil {
				return err
			}
		}
	}

	return nil
}