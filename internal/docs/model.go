package docs

type Category string

const (
	CategoryForm        Category = "form"
	CategoryDataDisplay Category = "data-display"
	CategoryFeedback    Category = "feedback"
	CategoryLayout      Category = "layout"
	CategoryNavigation  Category = "navigation"
	CategoryUtility     Category = "utility"
	CategoryGuide       Category = "guide"
)

type DocEntry struct {
	Name        string   `json:"name"`
	Category    Category `json:"category"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
}

type DocEntrySummary struct {
	Name        string   `json:"name"`
	Category    Category `json:"category"`
	Description string   `json:"description"`
}
