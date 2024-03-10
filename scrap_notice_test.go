package uiuscraper

import (
	"testing"
)

func TestScrapNotice(t *testing.T) {
	testCases := []struct {
		testName    string
		department  Department
		allowDomain string
		noticeSite  string
	}{
		{"UIU Notices", DepartmentAll, AllowDomainUIU, Notice_Site_UIU},
		{"CSE Notices", DepartmentCSE, AllowDomainCSE, Notice_Site_CSE},
		{"EEE Notices", DepartmentEEE, AllowDomainEEE, Notice_Site_EEE},
		{"CE Notices", DepartmentCivil, AllowDomainCE, Notice_Site_CE},
		{"Pharmacy Notices", DepartmentPharmacy, AllowDomainPharmacy, Notice_Site_Pharmacy},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			config := setupConfig(tc.department, tc.allowDomain, tc.noticeSite)

			notices := ScrapNotice(config)
			logAndAssertNotEmpty(t, notices)

			config.LastNoticeId = &notices[0].ID
			notices = ScrapNotice(config)
			logAndAssertNotEmpty(t, notices)
		})
	}
}

func TestGenerateNoticeID(t *testing.T) {
	title := "Test Title"
	date := "2021-07-20 00:00:00 +0000 UTC"

	noticeID := GenerateNoticeID(title, date)
	if noticeID != GenerateNoticeID(title, date) {
		t.Errorf("GenerateNoticeID() returned wrong value")
	}
}

func setupConfig(department Department, allowDomain string, noticeSite string) *NoticeScrapConfig {
	config := NoticeScrapConfig{
		LastNoticeId: nil,
		Department:   department,
		AllowDomain:  allowDomain,
		NOTICE_SITE:  noticeSite,
	}
	return &config
}

func logAndAssertNotEmpty(t *testing.T, notices []Notice) {
	totalNotices := len(notices)
	t.Logf("Total Notices: %d", totalNotices)

	if totalNotices == 0 {
		t.Errorf("ScrapNotice() returned no notices")
	}
}
