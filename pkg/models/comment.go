package models

// Comment represents a ClickUp task comment
type Comment struct {
	ID          string    `json:"id"`
	CommentText string    `json:"comment_text"`
	User        User      `json:"user"`
	DateCreated Timestamp `json:"date"`
	Resolved    bool      `json:"resolved"`
}

// CommentsResponse is the API response for listing comments
type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}
