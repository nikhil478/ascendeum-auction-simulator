package bidder

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nikhil478/ascendeum-action-simulator/internal/auction"
)

type Bidder struct {
	ID                  string
	ResponseProbability float64
	MinDelay, MaxDelay  time.Duration
	MinBid, MaxBid      float64
	Rng                 *rand.Rand
}

func NewBidders(n int, seed int64) []*Bidder {
	g := rand.New(rand.NewSource(seed))
	out := make([]*Bidder, n)
	for i := 0; i < n; i++ {
		out[i] = &Bidder{
			ID:                  fmt.Sprintf("bidder-%03d", i+1),
			ResponseProbability: 0.85,
			MinDelay:            10 * time.Millisecond,
			MaxDelay:            800 * time.Millisecond,
			MinBid:              1.0,
			MaxBid:              100.0,
			Rng:                 rand.New(rand.NewSource(g.Int63())),
		}
	}
	return out
}

func (b *Bidder) SimulateBid() (float64, bool) {
	if b.Rng.Float64() > b.ResponseProbability {
		return 0, false
	}
	delay := b.MinDelay
	if b.MaxDelay > b.MinDelay {
		delay += time.Duration(b.Rng.Int63n(int64(b.MaxDelay-b.MinDelay)))
	}
	time.Sleep(delay)
	amount := b.MinBid + b.Rng.Float64()*(b.MaxBid-b.MinBid)
	return amount, true
}

// safeSend protects against send on closed channel.
func safeSend(ctx context.Context, ch chan<- auction.Bid, bid auction.Bid) {
	defer func() { _ = recover() }() // catch panic if channel closed
	select {
	case <-ctx.Done():
		return
	case ch <- bid:
	}
}

// Start launches bidder goroutines.
func Start(ctx context.Context, bidders []*Bidder, auctionIDs []string, eng *auction.Engine) {
	for _, bidder := range bidders {
		b := bidder
		go func() {
			for _, aid := range auctionIDs {
				select {
				case <-ctx.Done():
					return
				default:
					if amt, ok := b.SimulateBid(); ok {
						safeSend(ctx, eng.BidChannel(), auction.Bid{
							AuctionID: aid,
							BidderID:  b.ID,
							Amount:    amt,
							// Time will be stamped by Engine.SubmitBid
						})
					}
				}
			}
		}()
	}
}
