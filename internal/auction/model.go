package auction

import "time"

type Bid struct {
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Time      time.Time `json:"time"`
}

type Auction struct {
	AuctionID string    `json:"auction_id"`
	StartedAt time.Time `json:"started_at"`
	Deadline  time.Time `json:"deadline"`
	Winner    *Bid      `json:"winner,omitempty"`
}