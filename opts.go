package main

import (
	"log"
	"os"
	"path/filepath"
)

type options struct {
	Port     int    `short:"p" long:"port" description:"the port to use for the web server" default:"9000"`
	NoGen    bool   `short:"G" long:"no-generation" description:"when set, the site is not automatically generated" default:false`
	SiteName string `short:"n" long:"site-name" description:"the name of the site" default:"Site Name"`
	TagLine  string `short:"t" long:"tag-line" description:"the site's tag line"`
}

var (
	Options      options
	PublicDir    string
	PostsDir     string
	TemplatesDir string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL ", err)
	}
	PublicDir = filepath.Join(pwd, "public")
	PostsDir = filepath.Join(pwd, "posts")
	TemplatesDir = filepath.Join(pwd, "templates")
}
