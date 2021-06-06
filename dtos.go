package blogaggregatormodule

import "github.com/mmcdole/gofeed"

type AddBlogRequest struct {
	Url string `json:"url"`
}

type AuthorCreatedPayload struct {
	BlogID string `json:"blogId"`
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
type PostCreatedPayload struct {
	gofeed.Item
	BlogID string `json:"blogId"`
}