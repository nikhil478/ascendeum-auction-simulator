package bidder

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nikhil478/ascendeum-action-simulator/internal/auction"
)

// Bidder represents a simulated bidder with configurable behaviour.
type Bidder struct {
	ID                  string
	ResponseProbability float64
	MinDelay, MaxDelay  time.Duration
	MinBid, MaxBid      float64
	Rng                 *rand.Rand
}

// NewBidders creates n bidders with deterministic random seeds.
func NewBidders(n int, seed int64) []*Bidder {
	globalRand := rand.New(rand.NewSource(seed))
	bidders := make([]*Bidder, n)

	for i := 0; i < n; i++ {
		bidders[i] = &Bidder{
			ID:                  fmt.Sprintf("bidder-%03d", i+1),
			ResponseProbability: 0.85,
			MinDelay:            10 * time.Millisecond,
			MaxDelay:            800 * time.Millisecond,
			MinBid:              1.0,
			MaxBid:              100.0,
			Rng:                 rand.New(rand.NewSource(globalRand.Int63())),
		}
	}
	return bidders
}

// SimulateBid waits a random delay and returns a bid amount
// if this bidder chooses to respond.  ok==false means no bid.
func (b *Bidder) SimulateBid() (amount float64, ok bool) {
	if b.Rng.Float64() > b.ResponseProbability {
		return 0, false
	}
	delay := b.MinDelay
	if b.MaxDelay > b.MinDelay {
		delay += time.Duration(b.Rng.Int63n(int64(b.MaxDelay - b.MinDelay)))
	}
	time.Sleep(delay)

	amount = b.MinBid + b.Rng.Float64()*(b.MaxBid-b.MinBid)
	return amount, true
}

// safeSend sends a bid into ch, but never panics even if the channel
// has been closed by the engine.
func safeSend(ctx context.Context, ch chan<- auction.Bid, b auction.Bid) bool {
	defer func() {
		// Recover if the channel was closed between select
		// and the actual send.
		_ = recover()
	}()
	select {
	case <-ctx.Done():
		return false // simulation has ended
	case ch <- b:
		return true
	}
}

// Start launches one goroutine per bidder to send bids for each auction ID.
// It stops sending when ctx is cancelled or when the engine's channel is closed.
func Start(ctx context.Context, bidders []*Bidder, auctionIDs []string, eng *auction.Engine) {
	for _, bidder := range bidders {
		b := bidder
		go func() {
			for _, aid := range auctionIDs {
				select {
				case <-ctx.Done():
					return
				default:
					go func() {
						if amt, ok := b.SimulateBid(); ok {
							safeSend(ctx, eng.BidChannel(), auction.Bid{
								AuctionID: aid,
								BidderID:  b.ID,
								Amount:    amt,
								Time:      time.Now(), // arrival time
							})
						}
					}()
				}
			}
		}()
	}
}
