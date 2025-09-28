package auction

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/nikhil478/ascendeum-action-simulator/internal/resource"
)

type Engine struct {
	mu       sync.RWMutex
	auctions map[string]*LiveAuction
	bidCh    chan Bid
	startT   time.Time
	timeout  time.Duration
}

func NewEngine(num int, timeout time.Duration) *Engine {
	start := time.Now()
	e := &Engine{
		auctions: make(map[string]*LiveAuction, num),
		bidCh:    make(chan Bid, num*100),
		startT:   start,
		timeout:  timeout,
	}
	for i := 0; i < num; i++ {
		id := fmt.Sprintf("auction-%04d", i+1)
		attrs := generateAttributes() // implement as needed
		e.auctions[id] = NewLiveAuction(id, attrs, start, start.Add(timeout))
	}
	return e
}

func (e *Engine) AuctionList() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ids := make([]string, 0, len(e.auctions))
	for id := range e.auctions {
		ids = append(ids, id)
	}
	return ids
}

func (e *Engine) BidChannel() chan<- Bid { return e.bidCh }


// SubmitBid is called by bidders; it stamps arrival time immediately.
func (e *Engine) SubmitBid(b Bid) {
	b.Time = time.Now()
	select {
	case e.bidCh <- b:
		// ok
	default:
		// buffer full: drop or handle as needed
	}
}

// Start launches worker goroutines that process bids until ctx is done.
func (e *Engine) Start(ctx context.Context, workers int) {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case b := <-e.bidCh:
					e.placeBid(b)
				}
			}
		}()
	}
	// Wait for ctx cancel then close auctions and wait workers
	go func() {
		<-ctx.Done()
		e.closeAll()
		wg.Wait()
	}()
}

func (e *Engine) placeBid(b Bid) {
	e.mu.RLock()
	a, ok := e.auctions[b.AuctionID]
	e.mu.RUnlock()
	if ok {
		a.AcceptBid(b)
	}
}

func (e *Engine) closeAll() {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, a := range e.auctions {
		a.Close()
	}
}

func (e *Engine) AuctionIDs() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ids := make([]string, 0, len(e.auctions))
	for id := range e.auctions {
		ids = append(ids, id)
	}
	return ids
}

// GenerateReport writes each auction's result and the global summary.
func (e *Engine) GenerateReport(outDir string, resMeta resource.Metadata) error {
	end := time.Now()
	summary := NewSummary(e.startT, end, len(e.auctions), resMeta)

	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, a := range e.auctions {
		result := a.Snapshot()
		if err := saveJSON(filepath.Join(outDir, result.AuctionID+".json"), result); err != nil {
			return err
		}
		summary.Add(result.AuctionID + ".json")
	}
	return saveJSON(filepath.Join(outDir, "global_summary.json"), summary)
}

// Example placeholder: produce 20 attributes per auction.
func generateAttributes() map[string]interface{} {
	attrs := make(map[string]interface{})
	for i := 0; i < 20; i++ {
		attrs[fmt.Sprintf("attr_%02d", i+1)] = i
	}
	return attrs
}
