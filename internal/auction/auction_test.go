package auction

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLiveAuction_AcceptBid(t *testing.T) {
	start := time.Now()
	deadline := start.Add(1 * time.Hour)
	auction := NewLiveAuction("1", map[string]interface{}{"item": "book"}, start, deadline)

	bid1 := Bid{BidderID: "Alice", Amount: 100, Time: start.Add(10 * time.Minute)}
	bid2 := Bid{BidderID: "Bob", Amount: 150, Time: start.Add(20 * time.Minute)}
	bid3 := Bid{BidderID: "Charlie", Amount: 150, Time: start.Add(5 * time.Minute)} // same amount, earlier

	// Accept bids
	auction.AcceptBid(bid1)
	auction.AcceptBid(bid2)
	auction.AcceptBid(bid3)

	snap := auction.Snapshot()

	if len(snap.Bids) != 3 {
		t.Errorf("expected 3 bids, got %d", len(snap.Bids))
	}

	if snap.Winner.BidderID != "Charlie" {
		t.Errorf("expected winner to be Charlie, got %v", snap.Winner.BidderID)
	}
}

func TestLiveAuction_AcceptBid_AfterDeadline(t *testing.T) {
	start := time.Now()
	deadline := start.Add(30 * time.Minute)
	auction := NewLiveAuction("2", nil, start, deadline)

	bid := Bid{BidderID: "Alice", Amount: 50, Time: start.Add(31 * time.Minute)}
	auction.AcceptBid(bid)

	snap := auction.Snapshot()
	if len(snap.Bids) != 0 {
		t.Errorf("expected 0 bids because bid was after deadline, got %d", len(snap.Bids))
	}

	if snap.Winner != nil {
		t.Errorf("expected no winner, got %v", snap.Winner)
	}
}

func TestLiveAuction_Close(t *testing.T) {
	start := time.Now()
	deadline := start.Add(1 * time.Hour)
	auction := NewLiveAuction("3", nil, start, deadline)

	auction.Close()
	snap := auction.Snapshot()

	if snap.Status != "closed" {
		t.Errorf("expected auction status 'closed', got %s", snap.Status)
	}

	if snap.ClosedAt.IsZero() {
		t.Errorf("expected ClosedAt to be set, got zero value")
	}
}

// TestParallelProcessingDifferentAuctions ensures that bids for different auctions can be processed in parallel.
func TestParallelProcessingDifferentAuctions(t *testing.T) {
	engine := NewEngine(2, 1*time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine.Start(ctx, 2) // start 2 workers

	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()
	var auction1Time, auction2Time time.Time

	// Bid for auction-0001
	go func() {
		defer wg.Done()
		engine.SubmitBid(Bid{AuctionID: "auction-0001", BidderID: "Alice", Amount: 100})
		auction1Time = time.Now()
	}()

	// Bid for auction-0002
	go func() {
		defer wg.Done()
		engine.SubmitBid(Bid{AuctionID: "auction-0002", BidderID: "Bob", Amount: 150})
		auction2Time = time.Now()
	}()

	wg.Wait()

	if auction1Time.Sub(start) > 50*time.Millisecond && auction2Time.Sub(start) > 50*time.Millisecond {
		t.Errorf("Expected at least one auction to process quickly in parallel, but both took too long")
	}
}

// TestSerialProcessingSameAuction ensures that bids for the same auction are processed serially.
func TestSerialProcessingSameAuction(t *testing.T) {
	engine := NewEngine(1, 1*time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine.Start(ctx, 1) // single worker

	auctionID := "auction-0001"
	numBids := 5

	// Submit bids sequentially (to test mutex-based serial processing)
	for i := 0; i < numBids; i++ {
		bid := Bid{
			AuctionID: auctionID,
			BidderID:  string(rune('A' + i)),
			Amount:    float64(100 + i),
			Time:      time.Now(),
		}
		engine.SubmitBid(bid)
	}

	// Wait for all bids to be processed
	time.Sleep(100 * time.Millisecond)

	// Collect bids from snapshot
	snap := engine.auctions[auctionID].Snapshot()
	if len(snap.Bids) != numBids {
		t.Fatalf("expected %d bids, got %d", numBids, len(snap.Bids))
	}

	// Verify that the bid order matches submission order
	for i, b := range snap.Bids {
		expected := string(rune('A' + i))
		if b.BidderID != expected {
			t.Errorf("expected bid %d to be %s, got %s", i, expected, b.BidderID)
		}
	}
}
