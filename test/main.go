package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"os"
	"time"

	uiuscraper "github.com/monzim/uiu-notice-scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cases := []struct {
		testName    string
		department  uiuscraper.Department
		allowDomain string
		noticeSite  string
	}{
		{"UIU Notices", uiuscraper.DepartmentAll, uiuscraper.AllowDomainUIU, uiuscraper.Notice_Site_UIU},
		{"CSE Notices", uiuscraper.DepartmentCSE, uiuscraper.AllowDomainCSE, uiuscraper.Notice_Site_CSE},
		{"EEE Notices", uiuscraper.DepartmentEEE, uiuscraper.AllowDomainEEE, uiuscraper.Notice_Site_EEE},
		{"CE Notices", uiuscraper.DepartmentCivil, uiuscraper.AllowDomainCE, uiuscraper.Notice_Site_CE},
		{"Pharmacy Notices", uiuscraper.DepartmentPharmacy, uiuscraper.AllowDomainPharmacy, uiuscraper.Notice_Site_Pharmacy},
	}

	for _, tc := range cases {
		config := setupConfig(tc.department, tc.allowDomain, tc.noticeSite)

		notices := uiuscraper.ScrapNotice(config)
		logAndAssertNotEmpty(tc.department, notices)
	}
}

func setupConfig(department uiuscraper.Department, allowDomain string, noticeSite string) *uiuscraper.NoticeScrapConfig {
	config := uiuscraper.NoticeScrapConfig{
		LastNoticeId: nil,
		Department:   department,
		AllowDomain:  allowDomain,
		NOTICE_SITE:  noticeSite,
	}
	return &config
}

func logAndAssertNotEmpty(department uiuscraper.Department, notices []uiuscraper.Notice) {
	totalNotices := len(notices)
	log.Info().Msgf("Total Notices: %d", totalNotices)

	// save the notices to a local file
	jsonNotices, err := json.MarshalIndent(notices, "", " ")
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling notices")
		return
	}

	_ = ioutil.WriteFile(fmt.Sprintf("%s_notices_%s.json", department, time.Now().Format("2006-01-02")), jsonNotices, 0644)

	if totalNotices == 0 {
		log.Error().Msg("ScrapNotice() returned no notices")
	}
}
