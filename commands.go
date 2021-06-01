package blogaggregatormodule

import (
	"encoding/json"
	"github.com/wepala/weos"
)

func AddBlogCommand(url string) *weos.Command {
	payload := &AddBlogRequest{
		Url: url,
	}
	payloadJson,_ := json.Marshal(payload)
	return &weos.Command{
		Type: "blog.add",
		Payload: payloadJson,
		Metadata: weos.CommandMetadata{
			Version: 1,
		},
	}
}