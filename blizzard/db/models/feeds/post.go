package feeds

type Post struct {
	ID        string `bun:",type:uuid,unique,default:gen_random_uuid()" json:"id"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
	//Author        user.Author `json:"author"`
	Content string `json:"content"`
}
