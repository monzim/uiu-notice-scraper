package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"os"
	"time"

	uiuscraper "github.com/monzim/uiu-notice-scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Ping(url string) {
	log.Info().Msgf("Pinging %s", url)
	_, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Error pinging %s", url)
	}
}

func main() {

	// pingUrl := "https://camo.githubusercontent.com/c23d0e21f14f15b7f1f631e3b73e49375694e90b5cc7330afeefe190023d0183/68747470733a2f2f6b6f6d617265762e636f6d2f67687076632f3f757365726e616d653d6d6f6e7a696d266c6162656c3d50726f66696c65253230766965777326636f6c6f723d306537356236267374796c653d666c6174"

	// // do 10000 requests to the ping url to check the performance use go routine to do parallel requests
	// for i := 0; i < 10000; i++ {
	// 	defer ping(pingUrl)
	// }

	// return

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
