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
	URL string `json:"url,omitempty"`
	FeedURL string `json:"feedUrl,omitempty"`
	Authors []*Author
	Posts []*Post
}

func (b *Blog) Init(blog *AddBlogRequest) (*Blog, error) {
	err := b.Validate(blog)
	if err != nil {
		return nil, err
	}
	createdPayload := &BlogCreatedPayload{
		Blog: Blog {
			URL: blog.Url,
		},
	}
	createdPayload.ID = GenerateID()
	event, _ := weos.NewBasicEvent(BLOG_ADDED,createdPayload.ID,"Blog",createdPayload)
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
	updatedPayload := &BlogCreatedPayload{
		Blog: Blog {
			Title: feed.Title,
		},
	}
	//set the link to be the blog url if link is not empty
	if feed.Link != "" && b.URL != feed.Link {
		updatedPayload.URL = feed.Link
		updatedPayload.FeedURL = b.URL
	}

	//update blog info based on what is returned in the feed
	event, _ := weos.NewBasicEvent(BLOG_UPDATED,b.ID,"Blog",updatedPayload)
	b.NewChange(event)
	b.ApplyChanges([]*weos.Event{event})

	for _,author := range feed.Authors {
		event, err := weos.NewBasicEvent(AUTHOR_CREATED,b.ID,"Blog",&AuthorCreatedPayload{
			BlogID: b.ID,
			Name: author.Name,
			Email: author.Email,
		})
		if err != nil {
			return err
		}
		b.NewChange(event)
		b.ApplyChanges([]*weos.Event{event})
	}

	for _, post := range feed.Items {
		event, err := weos.NewBasicEvent(POST_CREATED,b.ID,"Blog",&PostCreatedPayload{
			BlogID: b.ID,
			Item: *post,
		})
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
		b.SequenceNo = change.Meta.SequenceNo
		switch change.Type {
		case BLOG_ADDED:
			b.ID = change.Meta.EntityID
			err := json.Unmarshal(change.Payload, b)
			if err != nil {
				return err
			}
		case BLOG_UPDATED:
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