package blogaggregatormodule

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

func GetFeedLink(baseUrl string,body io.Reader) string {
	tokenizer := html.NewTokenizer(body)
		for {
			tt := tokenizer.Next()
			switch tt {
			case html.ErrorToken:
				return ""
			case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
				token := tokenizer.Token()
				if token.Data == "link" {
					var href string
					var isFeedLink bool

					for _, attr := range token.Attr {
						if !isFeedLink {
							isFeedLink = attr.Key == "rel" && attr.Val == "alternate"
						}
					
						if attr.Key == "href" {
							href = attr.Val
							turl, err := url.Parse(href)
							if err != nil {
								return ""
							}
							//it's a relative link
							if turl.Host == "" {
								href = baseUrl + href
							}
						}
					}
					if isFeedLink {
						return href
					}
				}
			}
		}
}