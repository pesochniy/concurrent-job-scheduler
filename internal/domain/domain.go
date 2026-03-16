package domain

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type FetchURL struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
}

type Report struct {
	UserID string `json:"user_id"`
	From   string `json:"from"`
	To     string `json:"to"`
}
