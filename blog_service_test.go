package blogaggregatormodule_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	blogaggregatormodule "github.com/wepala/blog-aggregator-module"
	"github.com/wepala/go-testhelpers"
)

func TestBlogServiceCreate_BasicBlog(t *testing.T) {
	blogDataFetched := 0
	//setup blog service with a mocked response for the http client
	service := blogaggregatormodule.NewBlogService(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
		blogDataFetched += 1
		//this is fetching the blog page 
		if blogDataFetched == 1 {
			content, err := ioutil.ReadFile("fixtures/html/ak33m_com.html")
			if err != nil {
				t.Fatalf("error setting up test fixtures '%s'",err)
			}
			resp := testhelpers.NewBytesResponse(200,content)
			resp.Header.Set("Content-Type", "text/html; charset=utf-8")
			return resp
		}

		content, err := ioutil.ReadFile("fixtures/feed/ak33m_com.xml")
		if err != nil {
			t.Fatalf("error setting up test fixtures '%s'",err)
		}
		resp := testhelpers.NewBytesResponse(200,content)
		resp.Header.Set("Content-Type", "application/rss+xml")
		return resp
		
	}),nil)
	blog, err := service.AddBlog(&blogaggregatormodule.AddBlogRequest{
		Url: "https://ak33m.com",
	})
	if err != nil {
		t.Fatalf("unexpected error creating blog '%s'",err)
	}

	if blog == nil {
		t.Fatal("expected blog to be returned")
	}

	if blogDataFetched != 2 {
		t.Errorf("expected the blog information to be fetched, called %d times",blogDataFetched)
	}
}

func TestBlogServiceCreate_InvalidBlog(t *testing.T) {
	t.Run("non blog page",func(t *testing.T) {
		//setup blog service with a mocked response for the http client
		service := blogaggregatormodule.NewBlogService(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
			content, err := ioutil.ReadFile("fixtures/html/google.html")
			if err != nil {
				t.Fatalf("error setting up test fixtures '%s'",err)
			}
			resp := testhelpers.NewBytesResponse(200,content)
			resp.Header.Set("Content-Type", "text/html")
			return resp
			
		}),nil)
		_, err := service.AddBlog(&blogaggregatormodule.AddBlogRequest{
			Url: "https://google.com",
		})
		if err == nil {
			t.Fatalf("expected an error from blog service")
		}
	})
}