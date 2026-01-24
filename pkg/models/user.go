package models

// User represents a ClickUp user
type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
	Initials       string `json:"initials"`
}

// Team represents a ClickUp workspace/team
type Team struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Avatar  string   `json:"avatar"`
	Members []Member `json:"members"`
}

// Member represents a team member
type Member struct {
	User User `json:"user"`
}

// UserResponse is the API response for GET /user
type UserResponse struct {
	User User `json:"user"`
}

// TeamsResponse is the API response for GET /team
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}
