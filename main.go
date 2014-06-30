package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

// This structure holds the command-line options.
type options struct {
	Port             int    `short:"p" long:"port" description:"the port to use for the web server" default:"9000"`
	GenOnly          bool   `short:"g" long:"generate-only" description:"generate the static site and exit"`
	NoGen            bool   `short:"G" long:"no-generation" description:"when set, the site is not automatically generated"`
	SiteName         string `short:"n" long:"site-name" description:"the name of the site" default:"Site Name"`
	TagLine          string `short:"t" long:"tag-line" description:"the site's tag line"`
	RecentPostsCount int    `short:"r" long:"recent-posts" description:"the number of recent posts to send to the templates" default:"5"`
	BaseURL          string `short:"b" long:"base-url" description:"the base URL of the web site" default:"http://localhost"`
}

type siteMeta struct {
	meta        TemplateData
	recentPosts int
}

var (
	// The one and only Options parsed from the command-line
	Options options

	PublicDir    string   // Public directory path
	PostsDir     string   // Posts directory path
	TemplatesDir string   // Templates directory path
	RssURL       string   // The RSS feed URL, parsed only once and stored for convenience
	SiteMeta     siteMeta // The site meta data can be used by posts
)

func init() {
	// Initialize directories
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	PublicDir = filepath.Join(pwd, "public")
	PostsDir = filepath.Join(pwd, "posts")
	TemplatesDir = filepath.Join(pwd, "templates")
}

func storeRssURL() {
	b, err := url.Parse(Options.BaseURL)
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	r, err := b.Parse("/rss")
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	RssURL = r.String()
}

func copyMeta() {
	SiteMeta.recentPosts = Options.RecentPostsCount
	SiteMeta.meta = make(TemplateData)
	SiteMeta.meta["BaseURL"] = Options.BaseURL
	SiteMeta.meta["SiteName"] = Options.SiteName
	SiteMeta.meta["TagLine"] = Options.TagLine
	SiteMeta.meta["RssURL"] = RssURL
}

func main() {
	// Parse the flags
	_, err := flags.Parse(&Options)
	if err == nil { // err != nil prints the usage automatically
		storeRssURL()
		copyMeta()
		if !Options.NoGen {
			// Generate the site
			if err := generateSite(); err != nil {
				log.Println("Error ", err)
			}
			// Terminate if set to generate only
			if Options.GenOnly {
				return
			}
			// Start the watcher
			defer startWatcher().Close()
		}
		// Start the web server
		run()
	}
}
