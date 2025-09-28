package auction

import (
	"sync"
	"time"
)

type LiveAuction struct {
	Auction
	mu sync.Mutex
}

func NewLiveAuction(id string, start, deadline time.Time) *LiveAuction {
	return &LiveAuction{
		Auction: Auction{
			AuctionID: id,
			StartedAt: start,
			Deadline:  deadline,
		},
	}
}

// AcceptBid checks deadline and updates winner safely.
func (a *LiveAuction) AcceptBid(b Bid) {
	if time.Now().After(a.Deadline) {
		return // auction closed
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Winner == nil || b.Amount > a.Winner.Amount ||
		(b.Amount == a.Winner.Amount && b.Time.Before(a.Winner.Time)) {
		b.AuctionID = a.AuctionID
		a.Winner = &b
	}
}