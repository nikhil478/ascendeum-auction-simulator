package auction

import (
	"context"
	"testing"
	"time"
)

func TestEngine_CreateAndListAuctions(t *testing.T) {
	engine := NewEngine(5, 1*time.Hour)

	ids := engine.AuctionList()
	if len(ids) != 5 {
		t.Errorf("expected 5 auctions, got %d", len(ids))
	}

	for _, id := range ids {
		if id == "" {
			t.Errorf("found empty auction ID")
		}
	}
}

func TestEngine_SubmitAndProcessBid(t *testing.T) {
	engine := NewEngine(2, 1*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	engine.Start(ctx, 2)

	auctionIDs := engine.AuctionList()
	if len(auctionIDs) == 0 {
		t.Fatal("no auctions to bid on")
	}
	aid := auctionIDs[0]

	engine.SubmitBid(Bid{AuctionID: aid, BidderID: "Alice", Amount: 100})
	engine.SubmitBid(Bid{AuctionID: aid, BidderID: "Bob", Amount: 150})
	engine.SubmitBid(Bid{AuctionID: aid, BidderID: "Charlie", Amount: 120})

	time.Sleep(50 * time.Millisecond)

	engine.mu.RLock()
	auction := engine.auctions[aid]
	engine.mu.RUnlock()

	result := auction.Snapshot()
	if result.Winner == nil {
		t.Fatal("expected a winner, got nil")
	}
	if result.Winner.BidderID != "Bob" {
		t.Errorf("expected Bob as winner, got %s", result.Winner.BidderID)
	}
}

func TestEngine_CloseAllAuctions(t *testing.T) {
	engine := NewEngine(3, 1*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	engine.Start(ctx, 2)

	for _, aid := range engine.AuctionList() {
		engine.SubmitBid(Bid{AuctionID: aid, BidderID: "Alice", Amount: 50})
	}

	time.Sleep(50 * time.Millisecond)

	cancel()
	time.Sleep(50 * time.Millisecond)


	engine.mu.RLock()
	defer engine.mu.RUnlock()
	for id, a := range engine.auctions {
		if a.Snapshot().Status != "closed" {
			t.Errorf("auction %s not closed", id)
		}
		if a.Snapshot().ClosedAt.IsZero() {
			t.Errorf("auction %s ClosedAt not set", id)
		}
	}
}

func TestEngine_BidOnNonexistentAuction(t *testing.T) {
	engine := NewEngine(1, 1*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	engine.Start(ctx, 1)

	// Submit bid to invalid auction
	engine.SubmitBid(Bid{AuctionID: "nonexistent", BidderID: "Alice", Amount: 100})

	time.Sleep(20 * time.Millisecond) // allow processing

	// Ensure nothing crashed and valid auction remains unaffected
	engine.mu.RLock()
	defer engine.mu.RUnlock()
	if len(engine.auctions) != 1 {
		t.Errorf("expected 1 auction, got %d", len(engine.auctions))
	}
}
