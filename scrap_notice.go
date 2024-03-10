package uiuscraper

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

func ScrapNotice(config *NoticeScrapConfig) []Notice {
	if config == nil {
		log.Error().Msg("Config is nil")
		return nil

	}

	var lastId string

	if config.LastNoticeId != nil {
		lastId = *config.LastNoticeId
	}

	log.Info().Msgf("Scraping notices for department: %s", config.Department)

	c := colly.NewCollector(colly.AllowedDomains(config.AllowDomain))
	var notices []Notice
	stopNextPage := false
	page := 1

	c.OnHTML("div[class=notice]", func(e *colly.HTMLElement) {
		title := e.ChildText("div[class=title] a")
		image := e.ChildAttr("div[class=image] img", "src")
		date := e.ChildText("div[class=date-container] span[class=date]")
		link := e.ChildAttr("div[class=title] a", "href")

		parsedTime, err := time.Parse(LayoutTime, date)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing date")
			return
		}

		noticeID := GenerateNoticeID(title, parsedTime.String())
		if noticeID == lastId {
			stopNextPage = true
			return
		}

		notice := Notice{
			ID:         noticeID,
			Title:      title,
			Image:      image,
			Date:       parsedTime,
			Link:       link,
			ScrapedAt:  time.Now(),
			Department: config.Department,
		}

		notices = append(notices, notice)
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	log.Info().Str("url", r.URL.String()).Msg("Visiting")
	// })

	c.OnHTML("div[class=nav-links]", func(e *colly.HTMLElement) {
		log.Info().
			Msgf("Scraping Department: %s, Page: %d", config.Department, page)
		page++

		if stopNextPage {
			log.Info().Msg("Stopping Already have up to date notices")
			return
		}

		nextPage := e.ChildAttr("a.next.page-numbers", "href")
		if nextPage == "" {
			log.Info().Msg("Don't have next page")
			return
		}

		c.Visit(nextPage)
	})

	c.Visit(config.NOTICE_SITE)

	if len(notices) == 0 {
		log.Warn().Msgf("No notices found for department: %s", config.Department)
	}

	log.Info().Msgf("Scraped %d notices for department: %s", len(notices), config.Department)

	return notices
}
