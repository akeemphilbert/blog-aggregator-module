package blogaggregatormodule

import (
	"encoding/json"

	"github.com/mmcdole/gofeed"
	"github.com/wepala/weos"
)

type Author struct {
	Name string
	Email string
}

type Post struct {
	Title string
	Description string
	Content string
}
type Blog struct {
	weos.AggregateRoot
	Title string
	Description string
	URL string `json:"url"`
	Authors []*Author
	Posts []*Post
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

func (b *Blog) AddFeed(feed *gofeed.Feed) error {
	b.Title = feed.Title
	for _,author := range feed.Authors {
		event, err := weos.NewBasicEvent(AUTHOR_CREATED,GenerateID(),"Blog",author)
		if err != nil {
			return err
		}
		b.NewChange(event)
		b.ApplyChanges([]*weos.Event{event})
	}

	for _, post := range feed.Items {
		event, err := weos.NewBasicEvent(POST_CREATED,GenerateID(),"Blog",post)
		if err != nil {
			return err
		}
		b.NewChange(event)
		b.ApplyChanges([]*weos.Event{event})
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
		case AUTHOR_CREATED:
			var author *Author
			err := json.Unmarshal(change.Payload, &author)
			if err != nil {
				return err
			}
			b.Authors = append(b.Authors,author)
		case POST_CREATED:
			var post *Post
			err := json.Unmarshal(change.Payload, &post)
			if err != nil {
				return err
			}
			b.Posts = append(b.Posts,post)
		}
	}

	return nil
}