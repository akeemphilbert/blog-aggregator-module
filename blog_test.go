package blogaggregatormodule_test

import (
	"math/rand"
	"testing"
	"time"

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