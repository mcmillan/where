package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/dustin/go-humanize"
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
)

type APIResult struct {
	Response APIResponse
}

type APIResponse struct {
	Checkins Checkins
}

type Checkins struct {
	Items []Checkin
}

type Checkin struct {
	Venue     Venue
	CreatedAt int64
}

type Venue struct {
	Name       string
	Location   Location
	Categories []Category
}

type Category struct {
	Icon Icon
}

type Icon struct {
	Prefix string
	Suffix string
}

type Location struct {
	Lat     float64
	Lng     float64
	City    string
	Country string
}

func (c *Checkin) TimeAgo() string {
	time := time.Unix(c.CreatedAt, 0)
	return humanize.Time(time)
}

func (c *Checkin) IconURL() string {
	prefix := c.Venue.Categories[0].Icon.Prefix
	suffix := c.Venue.Categories[0].Icon.Suffix

	return fmt.Sprintf("%s88%s", prefix, suffix)
}

func main() {
	godotenv.Load()
	m := martini.Classic()

	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", getLastCheckin())
	})

	m.Run()
}

func getLastCheckin() *Checkin {
	oauthToken := os.Getenv("FOURSQUARE_ACCESS_TOKEN")
	url := fmt.Sprintf("https://api.foursquare.com/v2/users/self/checkins?limit=1&oauth_token=%s&v=20140806&m=swarm", oauthToken)
	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var result APIResult

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		panic(err)
	}

	return &result.Response.Checkins.Items[0]
}
