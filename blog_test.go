package blogaggregatormodule_test

import (
	"encoding/json"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
)


func TestBlogCreate(t *testing.T) {
	blogUrls := []string{"https://ak33m.com","https://wepala.com","https://example.org"}
	rand.Seed(time.Now().UnixNano())
	selectedURL := blogUrls[rand.Intn(2)]

	t.Run("basic blog create",func(t *testing.T) {
		blog, err := new(blogaggregatormodule.Blog).Init(&blogaggregatormodule.AddBlogRequest{
			Url: selectedURL,
		})
		if err !=nil {
			t.Fatalf("unexpected error creating blog %s",err)
		}
	
		if blog.URL != selectedURL {
			t.Errorf("expected the blog url to be '%s', got '%s'",selectedURL,blog.URL)
		}

		if blog.ID == "" {
			t.Errorf("expected the blog id to be set")
		}

		//check for the blog id in the event payload 
		blogCreatedEvent := blog.GetNewChanges()[0].(*weos.Event)
		var blogCreatePayload struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(blogCreatedEvent.Payload,&blogCreatePayload)
		if err != nil {
			t.Errorf("unexpected error un marshalling blog create event %s",err)
		}
	
		if blogCreatePayload.ID != blog.ID {
			t.Errorf("expected the blog id to be %s, got %s",blog.ID,blogCreatePayload.ID)
		}
	})

	t.Run("don't create blog with missing url",func(t *testing.T) {
		_, err := new(blogaggregatormodule.Blog).Init(&blogaggregatormodule.AddBlogRequest{})
		if err ==nil {
			t.Error("should not be able to create invalid blog",err)
		}
	})
}

func TestBlogValidate(t *testing.T) {
	err := new(blogaggregatormodule.Blog).Validate(&blogaggregatormodule.AddBlogRequest{})
	if err == nil {
		t.Error("expected error for misisng url")
	}
}

func TestBlogAddFeed(t *testing.T) {
	blog,err := new(blogaggregatormodule.Blog).Init(&blogaggregatormodule.AddBlogRequest{
		Url: "https://ak33m.com/index.xml",
	})
	if err != nil {
		t.Fatalf("unexpected error initializing blog '%s'",err)
	}

	file, _ := os.Open("fixtures/feed/ak33m_com.xml")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	blog.AddFeed(feed)

	if blog.Title != "Akeem Philbert's Blog" {
		t.Errorf("expected the blog title to be '%s', got '%s'","Akeem Philbert's Blog",blog.Title)
	}

	if len(blog.GetNewChanges()) != 12{
		t.Fatalf("expected %d events, got %d",12,len(blog.GetNewChanges()))
	}

	//confirm that the author created event payload has blog id in it
	authorCreatedEvent := blog.GetNewChanges()[2].(*weos.Event)
	if authorCreatedEvent.Type != blogaggregatormodule.AUTHOR_CREATED {
		t.Errorf("expected the second event to be %s, got %s",blogaggregatormodule.AUTHOR_CREATED,authorCreatedEvent.Type)
	}
	var authorCreatePayload struct {
		BlogID string `json:"blogId"`
	}
	err = json.Unmarshal(authorCreatedEvent.Payload,&authorCreatePayload)
	if err != nil {
		t.Errorf("unexpected error un marshalling author create event %s",err)
	}

	if authorCreatePayload.BlogID != blog.ID {
		t.Errorf("expected the blog id to be %s, got %s",blog.ID,authorCreatePayload.BlogID)
	}


	if len(blog.Authors) != 1 {
		t.Errorf("expected %d author, got %d",1,len(blog.Authors))
	}

	//confirm that the post created event payload has blog id in it
	postCreatedEvent := blog.GetNewChanges()[3].(*weos.Event)
	if postCreatedEvent.Type != blogaggregatormodule.POST_CREATED {
		t.Errorf("expected the third event to be %s, got %s",blogaggregatormodule.POST_CREATED,postCreatedEvent.Type)
	}
	var postCreatePayload struct {
		BlogID string `json:"blogId"`
	}
	err = json.Unmarshal(postCreatedEvent.Payload,&postCreatePayload)
	if err != nil {
		t.Errorf("unexpected error un marshalling post create event %s",err)
	}

	if postCreatePayload.BlogID != blog.ID {
		t.Errorf("expected the blog id to be %s, got %s",blog.ID,postCreatePayload.BlogID)
	}


	if len(blog.Posts) != 9 {
		t.Errorf("expected %d posts, got %d",9,len(blog.Posts))
	}

	if blog.URL != "https://ak33m.com" {
		t.Errorf("expected the blog url to be '%s', got '%s'","https://ak33m.com",blog.URL)
	}
}