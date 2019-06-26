package model

// Offer struct represents one item parsed from njuskalo
type Offer struct {
	ID          int64
	QueryID     int64
	ItemID      string
	URL         string
	Name        string
	Image       string
	Price       string
	Description string
	Location    string
	Year        string
	Mileage     string
	Published   string
	CreatedAt   int64
}
