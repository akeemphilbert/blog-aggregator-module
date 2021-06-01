package blogaggregatormodule

import (
	"os"
	"testing"
)

func TestGetFeedLink (t *testing.T) {
	content, err := os.Open("fixtures/html/ak33m_com.html")
	if err != nil {
		t.Fatalf("error reading fixtures '%s'",err)
	}
	link := GetFeedLink(content)
	if link != "https://ak33m.com/index.xml" {
		t.Errorf("expected the link to be '%s', got '%s'","https://ak33m.com/index.xml",link)
	}
}