package auction

import "time"

type Bid struct {
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Time      time.Time `json:"time"`
}

type Auction struct {
	AuctionID  string                 `json:"auction_id"`
	Attributes map[string]interface{} `json:"attributes"`
	StartedAt  time.Time              `json:"start_time"`
	ClosedAt   time.Time              `json:"closed_at"`
	Deadline   time.Time              `json:"deadline_time"`
	Bids       []Bid                  `json:"bids"`
	Winner     *Bid                   `json:"winner"`
	Status     string                 `json:"status"`
}