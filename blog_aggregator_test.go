package blogaggregatormodule_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/segmentio/ksuid"
	"github.com/wepala/blog-aggregator-module"
	"github.com/wepala/weos"
)

type TestBlog struct
{
	Title string
	URL string
	FeedLink string
}

type TestUser struct
{
	Name string
	Site string
	IsLoggedIn bool
	Blog *TestBlog
}

type FeedItem struct {
	Title string 
	Description string
	Link string
	Category string
	PublishDate string
}

var testUsers map[string]*TestUser
var testBlogs map[string]*TestBlog
var testBlog *TestBlog
var testFeed string
var testCommand *weos.Command
var app *weos.BaseApplication
var err error
var currentID string

func reset(*godog.Scenario) {
	testBlog = nil
	testFeed = ""
	testCommand = nil
	testUsers = make(map[string]*TestUser)
	testBlogs = make(map[string]*TestBlog)
	err = nil
	currentID = ksuid.New().String()

	blogaggregatormodule.GenerateID = func() string {
		return currentID
	}

	
	
}

func aPingbackUrlShouldBeGenerated() error {
	return godog.ErrPending
}

func aUserNamed(arg1 string) error {
	testUsers[arg1] = &TestUser{
		Name: arg1,
	}
	return err
}

func anAuthorShouldBeCreatedForEachAuthorInTheFeed() error {
	return godog.ErrPending
}

func anErrorScreenShouldBeShown(arg1 string) error {
	return godog.ErrPending
}

func followsTheBlog(arg1, arg2 string) error {
	return nil
}

func hasABlog(arg1, arg2 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.Blog = &TestBlog{
			Title: arg2,
		}
		testBlogs[arg2] = user.Blog
		testBlog = user.Blog
		return err
	}
	err = fmt.Errorf("user %s not defined",arg1)
	return err
}

func hitsTheSubmitButton(arg1 string) error {
	err = app.Dispatcher().Dispatch(context.Background(),testCommand)
	return err
}

func isLoggedIn(arg1 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.IsLoggedIn = true
		return err
	}
	
	err =  fmt.Errorf("user %s not defined",arg1)
	return err
}

func isLoggedInWithGoogle(arg1 string) error {
	return godog.ErrPending
}

func isNotLoggedIn(arg1 string) error {
	if user,ok := testUsers[arg1]; ok {
		user.IsLoggedIn = false
		return nil
	}
	
	return fmt.Errorf("user %s not defined",arg1)
}

func isOnTheBlogSubmitScreen(arg1 string) error {
	return nil
}

func postsShouldBeCreatedForEachPost() error {
	return godog.ErrPending
}

func profilesForTheBlogAuthorsShouldBeCreated() error {
	return godog.ErrPending
}

func shouldBeRedirectedToTheProfilePageForThatBlog(arg1 string) error {
	return godog.ErrPending
}

func successfullyCompletesTheCaptcha(arg1 string) error {
	return nil
}

func successfullySubmitsAFeed(arg1 string) error {
	return godog.ErrPending
}

func theAggregatorSupportsAtomFeedsAsWellAsRssFeeds() error {
	return nil
}

func theBlogDetailsStoredInTheAggregator() error {
	return godog.ErrPending
}

func theBlogHasALinkToAFeed(arg1 string) error {
	testBlog.FeedLink = arg1
	return nil
}

func theBlogPostsFromTheFeedShouldBeAddedToTheAggregator() error {
	return godog.ErrPending
}

func theBlogShouldBeAddedToTheAggregator() error {
	//check to see that there is an event in the database for adding the blog
	events, err := app.EventRepository().GetByAggregateAndType(currentID,"Blog")
	if len(events) == 0 {
		err = errors.New("There should be an event for adding a blog '%s'")
	}
	return err
}

func theFeedDetailsShouldBeExtracted() error {
	return godog.ErrPending
}

func theFeedHasPosts(arg1 *messages.PickleStepArgument_PickleTable) error {
	var err error
	testFeed = `
	<?xml version="1.0" encoding="windows-1252"?>
	<rss version="2.0">
	  <channel>
		<title>%s</title>
		<description>RSS is a fascinating technology. The uses for RSS are expanding daily. Take a closer look at how various industries are using the benefits of RSS in their businesses.</description>
		<link>http://www.feedforall.com/industry-solutions.htm</link>
		<category domain="www.dmoz.com">Computers/Software/Internet/Site Management/Content Management</category>
		<copyright>Copyright 2004 NotePage, Inc.</copyright>
		<docs>http://blogs.law.harvard.edu/tech/rss</docs>
		<language>en-us</language>
		<lastBuildDate>Tue, 19 Oct 2004 13:39:14 -0400</lastBuildDate>
		<managingEditor>marketing@feedforall.com</managingEditor>
		<pubDate>Tue, 19 Oct 2004 13:38:55 -0400</pubDate>
		<webMaster>webmaster@feedforall.com</webMaster>
		<generator>FeedForAll Beta1 (0.0.1.8)</generator>
		<image>
		  <url>http://www.feedforall.com/ffalogo48x48.gif</url>
		  <title>FeedForAll Sample Feed</title>
		  <link>http://www.feedforall.com/industry-solutions.htm</link>
		  <description>FeedForAll Sample Feed</description>
		  <width>48</width>
		  <height>48</height>
		</image>
		%s
	  </channel>
	</rss>`
	//TODO loop through the table and add feed item to the feed 
	items := ""
	itemColumns := make([]string,len(arg1.Rows[0].Cells))
	for i,_ := range arg1.Rows {
		if i == 0 {
			for _,column := range arg1.Rows[i].Cells {
				itemColumns = append(itemColumns,column.Value)
			}
		} else {
			feedItem := &FeedItem{}
			for j,column := range arg1.Rows[i].Cells {
				if itemColumns[j] == "title" {
					feedItem.Title = column.Value
				}

				if itemColumns[j] == "content" {
					feedItem.Description = column.Value
				}

				if itemColumns[j] == "publish date" {
					feedItem.PublishDate = column.Value
				}
			}
			
			items = items + fmt.Sprintf(`<item>
			<title>%s</title>
			<description>%s</description>
			<link>%s</link>
			<pubDate>%s</pubDate>
		  </item>`,feedItem.Title,feedItem.Link, feedItem.Description,feedItem.PublishDate)

		}
	}


	testFeed = fmt.Sprintf(testFeed,testBlog.Title,"",items)
	return err
}

func theUrlIsEntered(arg1 string) error {
	testCommand = blogaggregatormodule.AddBlogCommand(arg1)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	appConfig := &weos.ApplicationConfig{
		ModuleID: "123",
		Title:    "Test App",
		Database: &weos.DBConfig{
			Driver: "sqlite3",
			Database: "test",
		},
		Log: nil,
	}
	app, err = weos.NewApplicationFromConfig(appConfig,nil,nil,nil,nil)
	//run migrations to setup all the necessary tables
	err = app.Migrate(context.TODO())
	err = blogaggregatormodule.Initialize(app)

	ctx.BeforeScenario(reset)
	ctx.Step(`^a pingback url should be generated$`, aPingbackUrlShouldBeGenerated)
	ctx.Step(`^a user named "([^"]*)"$`, aUserNamed)
	ctx.Step(`^an author should be created for each author in the feed$`, anAuthorShouldBeCreatedForEachAuthorInTheFeed)
	ctx.Step(`^an error screen should be shown "([^"]*)"$`, anErrorScreenShouldBeShown)
	ctx.Step(`^"([^"]*)" follows the blog "([^"]*)"$`, followsTheBlog)
	ctx.Step(`^"([^"]*)" has a blog "([^"]*)"$`, hasABlog)
	ctx.Step(`^"([^"]*)" hits the submit button$`, hitsTheSubmitButton)
	ctx.Step(`^"([^"]*)" is logged in$`, isLoggedIn)
	ctx.Step(`^"([^"]*)" is logged in with google$`, isLoggedInWithGoogle)
	ctx.Step(`^"([^"]*)" is not logged in$`, isNotLoggedIn)
	ctx.Step(`^"([^"]*)" is on the blog submit screen$`, isOnTheBlogSubmitScreen)
	ctx.Step(`^posts should be created for each post$`, postsShouldBeCreatedForEachPost)
	ctx.Step(`^profiles for the blog authors should be created$`, profilesForTheBlogAuthorsShouldBeCreated)
	ctx.Step(`^"([^"]*)" should be redirected to the profile page for that blog$`, shouldBeRedirectedToTheProfilePageForThatBlog)
	ctx.Step(`^"([^"]*)" successfully completes the captcha$`, successfullyCompletesTheCaptcha)
	ctx.Step(`^"([^"]*)" successfully submits a feed$`, successfullySubmitsAFeed)
	ctx.Step(`^the aggregator supports atom feeds as well as rss feeds$`, theAggregatorSupportsAtomFeedsAsWellAsRssFeeds)
	ctx.Step(`^the blog details stored in the aggregator$`, theBlogDetailsStoredInTheAggregator)
	ctx.Step(`^the blog has a link to a feed "([^"]*)"$`, theBlogHasALinkToAFeed)
	ctx.Step(`^the blog posts from the feed should be added to the aggregator$`, theBlogPostsFromTheFeedShouldBeAddedToTheAggregator)
	ctx.Step(`^the blog should be added to the aggregator$`, theBlogShouldBeAddedToTheAggregator)
	ctx.Step(`^the feed details should be extracted$`, theFeedDetailsShouldBeExtracted)
	ctx.Step(`^the feed has posts$`, theFeedHasPosts)
	ctx.Step(`^the url "([^"]*)" is entered$`, theUrlIsEntered)
}