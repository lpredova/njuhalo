package model

// Query represents one url that needs to be checked if active
type Query struct {
	ID        string
	Name      string
	IsActive  int64
	URL       string
	Filters   string
	CreatedAt int64
}
