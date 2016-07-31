package license

// License represents a license from the GitHub API.
type License struct {
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	URL            string   `json:"url"`
	HtmlURL        string   `json:"html_url"`
	Featured       bool     `json:"featured"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Implementation string   `json:"implementation"`
	Required       []string `json:"required"`
	Permitted      []string `json:"permitted"`
	Forbidden      []string `json:"forbidden"`
	Body           string   `json:"body"`
}
