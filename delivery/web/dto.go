package web

type Todo struct {
	ID        int    `json:"id"` // json tag for easy JSON encoding (like @JsonProperty in Java)
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}
