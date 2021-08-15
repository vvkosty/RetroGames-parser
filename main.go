package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const BASE_URL = "https://www.retrogames.cc"

type Game struct {
	rate  int
	title string
	url   string
}

var items []*Game
var games []*Game

func Parse() {
	i := 1
	for {
		res, err := http.Get(BASE_URL + "/mastersystem-games/page/" + strconv.Itoa(i) + ".html")
		checkErr(err)
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		checkErr(err)

		doc.Find(".item").Each(func(i int, s *goquery.Selection) {
			game := Game{}
			game.title = strings.TrimSpace(s.Find("a").Text())
			url, _ := s.Find("a").Attr("href")
			game.url = BASE_URL + url
			game.rate, err = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(s.Find(".pull-left").First().Text(), ",", "")))
			checkErr(err)
			games = append(games, &game)
		})

		if len(games) == 0 {
			break
		}

		items = append(items, games...)
		games = nil

		fmt.Printf("%d - %d - %d\r", i, len(games), len(items))
		i++
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].rate < items[j].rate
	})
	fmt.Printf("%v", items)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("rate: %d\ntitle: %s\nurl: %s\n", g.rate, g.title, g.url)
}

func main() {
	Parse()
}
