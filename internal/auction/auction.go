package auction

import (
	"sync"
	"time"
)


// LiveAuction manages a single auction while it is running.
type LiveAuction struct {
	mu   sync.Mutex
	data Auction
}

func NewLiveAuction(id string, attrs map[string]interface{}, start, deadline time.Time) *LiveAuction {
	return &LiveAuction{
		data: Auction{
			AuctionID:  id,
			Attributes: attrs,
			StartedAt:  start,
			Deadline:   deadline,
			Bids:       []Bid{},
			Status:     "open",
		},
	}
}

// AcceptBid records the bid if it arrived before deadline.
func (a *LiveAuction) AcceptBid(b Bid) {
	if b.Time.After(a.data.Deadline) {
		return // too late
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	a.data.Bids = append(a.data.Bids, b)
	if a.data.Winner == nil ||
		b.Amount > a.data.Winner.Amount ||
		(b.Amount == a.data.Winner.Amount && b.Time.Before(a.data.Winner.Time)) {
		a.data.Winner = &b
	}
}

// Close marks the auction closed and records ClosedAt.
func (a *LiveAuction) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data.Status = "closed"
	a.data.ClosedAt = time.Now()
}

// Snapshot returns a copy of the final AuctionResult.
func (a *LiveAuction) Snapshot() Auction {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.data
}
