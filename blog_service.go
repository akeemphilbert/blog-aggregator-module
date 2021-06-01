package blogaggregatormodule

import (
	"fmt"
	"net/http"

	"github.com/mmcdole/gofeed"
	"github.com/wepala/weos"
)

type BlogService struct {
	eventRepository weos.EventRepository
	client *http.Client
	validFeedTypes []string
}
//Add a new blog to the aggregator. 
func (b *BlogService) AddBlog(blogRequest *AddBlogRequest) (*Blog, error) {
	//pull the information from the blog and ensure there is a feed associated. If there isn't a feed then return an error
	blog, err := new(Blog).Init(blogRequest)
	if err != nil {
		return nil, err
	}
	feedURL := blog.URL
	response, err := b.client.Get(blog.URL)
	if err != nil {
		return nil, weos.NewDomainError(fmt.Sprintf("Unable to fetch feed '%s'",blog.URL),"Blog",blog.ID,err)
	}
	//if it's html let's get the feed link and get the feed information
	if response.Header.Get("Content-Type") == "text/html" || response.Header.Get("Content-Type") == "application/xhtml+xml" {
		feedURL = GetFeedLink(response.Body)
		if feedURL != "" {
			response.Body.Close()
			response,err = b.client.Get(feedURL)
			if err != nil {
				return nil, weos.NewDomainError(fmt.Sprintf("Unable to fetch feed '%s'",feedURL),"Blog",blog.ID,err)
			}
		}
	} 

	defer response.Body.Close()

	//parse the feed
	fp := gofeed.NewParser()
	feed, err := fp.Parse(response.Body)
	if err != nil {
		return nil, weos.NewDomainError(fmt.Sprintf("Unable to parse feed '%s'",feedURL),"Blog",blog.ID,err)
	}
	err = blog.AddFeed(feed)

	return blog, err
}

//Create new repository service
func NewBlogService(client *http.Client, eventRepository weos.EventRepository) *BlogService {
	return &BlogService{
		client: client,
		eventRepository: eventRepository,
		validFeedTypes: []string{"application/rss+xml","application/xml"},
	}
}