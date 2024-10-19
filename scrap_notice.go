package uiuscraper

import (
	"fmt"
	"strings"
	"time"
	"unicode"

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

	c := colly.NewCollector(
		colly.AllowedDomains(config.AllowDomain),
	)

	// Configure transport settings
	c.SetRequestTimeout(90 * time.Second)

	// Configure retry settings
	maxRetries := 3
	retryDelay := 5 * time.Second

	c.OnRequest(func(r *colly.Request) {
		retryCount, ok := r.Ctx.GetAny("retry_count").(int)
		if !ok {
			r.Ctx.Put("retry_count", 0)
		}
		log.Info().
			Str("url", r.URL.String()).
			Int("attempt", retryCount+1).
			Msg("Visiting")
	})

	c.OnError(func(r *colly.Response, err error) {
		retryCount, _ := r.Ctx.GetAny("retry_count").(int)

		if retryCount < maxRetries {
			log.Warn().
				Err(err).
				Str("url", r.Request.URL.String()).
				Int("retry", retryCount+1).
				Msg("Error while scraping, retrying...")

			time.Sleep(retryDelay)
			r.Ctx.Put("retry_count", retryCount+1)
			r.Request.Retry()
			return
		}

		log.Error().
			Err(err).
			Str("url", r.Request.URL.String()).
			Msg("Max retries reached, giving up")
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	var notices []Notice
	stopNextPage := false
	page := 1

	c.OnHTML("div[class=notice]", func(e *colly.HTMLElement) {
		title := e.ChildText("div[class=title] a")
		image := e.ChildAttr("div[class=image] img", "src")
		date := e.ChildText("div[class=date-container] span[class=date]")
		link := e.ChildAttr("div[class=title] a", "href")

		// log.Info().
		// 	Str("title", title).
		// 	Str("image", image).
		// 	Str("date", date).
		// 	Str("link", link).
		// 	Msg("Notice found")

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

		if link[len(link)-1] == '/' {
			link = link[:len(link)-1]
			link = fmt.Sprintf("%s?scrapper=%s", link, "github.com/monzim/uiu-notice-scraper")
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

		linkWithTrack := fmt.Sprintf("%s?track_id=%s", link, noticeID)
		err = c.Visit(linkWithTrack)
		if err != nil {
			if !strings.Contains(err.Error(), "Forbidden domain") {
				log.Error().
					Err(err).
					Str("url", linkWithTrack).
					Msg("Failed to visit notice detail page")
			}

		}
	})

	c.OnHTML("div[class=notice-details]", func(e *colly.HTMLElement) {
		summary := e.ChildText("p")
		noticeID := e.Request.URL.Query().Get("track_id")
		summary = removeExtraSpaces(summary)

		// log.Info().Str("notice_id", noticeID).Msg("Processing notice ")

		for i, notice := range notices {
			if notice.ID == noticeID {
				notices[i].Summary = summary
				break
			}
		}
	})

	c.OnHTML("div[class=nav-links]", func(e *colly.HTMLElement) {
		log.Info().
			Int("page", page).
			Str("department", string(config.Department)).
			Msg("Scraping page")

		page++
		if stopNextPage {
			log.Info().Msg("Stopping - Already have up to date notices")
			return
		}

		nextPage := e.ChildAttr("a.next.page-numbers", "href")
		if nextPage == "" {
			log.Info().Msg("No next page available")
			return
		}

		// Visit next page with built-in retry mechanism
		err := c.Visit(nextPage)
		if err != nil {
			log.Error().
				Err(err).
				Str("url", nextPage).
				Msg("Failed to visit next page")
		}
	})

	// Initial visit with retry mechanism
	err := c.Visit(config.NOTICE_SITE)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", config.NOTICE_SITE).
			Msg("Failed to start scraping")
		return nil
	}

	if len(notices) == 0 {
		log.Warn().
			Str("department", string(config.Department)).
			Msg("No notices found")
	}

	log.Info().
		Int("count", len(notices)).
		Str("department", string(config.Department)).
		Msg("Scraping completed")

	return notices
}

func removeExtraSpaces(paragraph string) string {
	paragraph = strings.Join(strings.Fields(paragraph), " ")
	var result strings.Builder
	var prev rune
	for _, char := range paragraph {
		if !unicode.IsSpace(prev) || !unicode.IsPunct(char) {
			result.WriteRune(char)
		}
		prev = char
	}
	return result.String()
}
