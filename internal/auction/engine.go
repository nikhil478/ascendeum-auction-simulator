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
}

func NewEngine(num int, timeout time.Duration) *Engine {
	start := time.Now()
	e := &Engine{
		auctions: make(map[string]*LiveAuction, num),
		bidCh:    make(chan Bid, num*100),
		startT:   start,
	}
	for i := 0; i < num; i++ {
		id := fmt.Sprintf("auction-%04d", i+1)
		e.auctions[id] = NewLiveAuction(id, start, start.Add(timeout))
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

// Start runs worker goroutines and closes bidCh when context is done.
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
				case b, ok := <-e.bidCh:
					if !ok {
						return
					}
					e.placeBid(b)
				}
			}
		}()
	}
	go func() {
		<-ctx.Done()
		close(e.bidCh)
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

func (e *Engine) GenerateReport(outDir string, resMeta resource.Metadata) error {
	end := time.Now()
	summary := NewSummary(e.startT, end, len(e.auctions), resMeta)

	for _, a := range e.auctions {
		if err := saveJSON(filepath.Join(outDir, a.AuctionID+".json"), a.Auction); err != nil {
			return err
		}
		summary.Add(a.AuctionID + ".json")
	}
	return saveJSON(filepath.Join(outDir, "global_summary.json"), summary)
}
