package blogaggregatormodule_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
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
		Url: "https://ak33m.com",
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

	if len(blog.GetNewChanges()) != 11 {
		t.Errorf("expected %d events, got %d",11,len(blog.GetNewChanges()))
	}

	if len(blog.Authors) != 1 {
		t.Errorf("expected %d author, got %d",1,len(blog.Authors))
	}

	if len(blog.Posts) != 9 {
		t.Errorf("expected %d posts, got %d",9,len(blog.Posts))
	}
}