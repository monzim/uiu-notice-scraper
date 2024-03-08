package uiuscraper

import "testing"

func TestGetNotice(t *testing.T) {

	notice := ScrapUIU(nil)
	if len(notice) == 0 {
		t.Errorf("ScrapUIU() returned no notices")
	}

	notice = ScrapUIU(&notice[0].ID)
	if len(notice) == 0 {
		return
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
