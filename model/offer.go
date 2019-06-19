package model

// Offer struct represents one item parsed from njuskalo
type Offer struct {
	ID          string
	QueryID     int64
	ItemID      int64
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
