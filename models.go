package uiuscraper

import "time"

type Department string

const (
	DepartmentAll      Department = "ALL"
	DepartmentCSE      Department = "BSCSE"
	DepartmentEEE      Department = "EEE"
	DepartmentCivil    Department = "CE"
	DepartmentPharmacy Department = "Pharmacy"
	DepartmentEnglish  Department = "English"
	DepartmentEDS      Department = "EDS"
	DepartmentMSJ      Department = "MSJ"
	DepartmentSoBE     Department = "SoBE"
)

type Notice struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Image      string     `json:"image"`
	Date       time.Time  `json:"date"`
	Link       string     `json:"link"`
	ScrapedAt  time.Time  `json:"scraped_at"`
	Department Department `json:"department"`
}

type NoticeScrapConfig struct {
	LastNoticeId *string
	Department   Department
	AllowDomain  string
	NOTICE_SITE  string
}
