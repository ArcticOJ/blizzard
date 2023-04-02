package problems

type Problem struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Contest string   `json:"contest"`
	Tags    []string `bun:"tags" json:"tags"`
	//Author  users.Author `json:"author"`
	Content string `json:"content"`
}
