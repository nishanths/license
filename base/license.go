package base

type License struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Category string `json:"category"`
	Body     string `json:"body"`
}
