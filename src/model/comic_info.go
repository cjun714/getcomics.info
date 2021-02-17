package model

type ComicInfo struct {
	ID           int      `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Cover        string   `json:"cover,omitempty"`
	PageURL      string   `json:"page_url,omitempty"`
	Category     string   `json:"category,omitempty"`
	StartYear    int      `json:"start_year,omitempty"`
	EndYear      int      `json:"end_year,omitempty"`
	Size         int      `json:"size,omitempty"`
	DownloadURLs []string `json:"download_urls,omitempty"`
}

type URLTypeEnum int
type DownloadURL struct {
	ID      int
	URl     string
	URLType URLTypeEnum
}
