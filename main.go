package uiuscraper

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

const (
	AllowDomain = "www.uiu.ac.bd"
	WebsiteURL  = "https://www.uiu.ac.bd/notice"
	layout      = "January 2, 2006"
)

type Notice struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Date      time.Time `json:"date"`
	Link      string    `json:"link"`
	ScrapedAt time.Time `json:"scraped_at"`
}

func GenerateNoticeID(title, date string) string {
	hashInput := fmt.Sprintf("%s-%s", title, date)
	hash := sha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}

func ScrapUIU(lastNoticeId *string) []Notice {
	var lastId string

	if lastNoticeId != nil {
		lastId = *lastNoticeId
	}

	c := colly.NewCollector(colly.AllowedDomains(AllowDomain))
	var notices []Notice
	stopScraping := false

	c.OnHTML("div[class=notice]", func(e *colly.HTMLElement) {
		title := e.ChildText("div[class=title] a")
		image := e.ChildAttr("div[class=image] img", "src")
		date := e.ChildText("div[class=date-container] span[class=date]")
		link := e.ChildAttr("div[class=title] a", "href")

		parsedTime, err := time.Parse(layout, date)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing date")
			return
		}

		noticeID := GenerateNoticeID(title, parsedTime.String())
		if noticeID == lastId {
			stopScraping = true
			return
		}

		notice := Notice{
			ID:        noticeID,
			Title:     title,
			Image:     image,
			Date:      parsedTime,
			Link:      link,
			ScrapedAt: time.Now(),
		}

		notices = append(notices, notice)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Info().Str("url", r.URL.String()).Msg("Visiting")
	})

	c.OnHTML("div[class=nav-links]", func(e *colly.HTMLElement) {
		if len(notices) == 0 || stopScraping {
			return
		}

		nextPage := e.ChildAttr("a.next.page-numbers", "href")
		if nextPage == "" {
			log.Info().Msg("No Next Page Found")
			return
		}

		c.Visit(nextPage)
	})

	c.Visit(WebsiteURL)

	if len(notices) == 0 {
		log.Warn().Msg("No notices found")
	}

	if stopScraping {
		return nil
	}

	return notices
}
