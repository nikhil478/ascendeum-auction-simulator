package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/nikhil478/ascendeum-action-simulator/internal/auction"
	"github.com/nikhil478/ascendeum-action-simulator/internal/bidder"
	"github.com/nikhil478/ascendeum-action-simulator/internal/resource"
)

func main() {
	numAuctions := flag.Int("auctions", 40, "number of concurrent auctions")
	numBidders := flag.Int("bidders", 100, "number of bidders")
	timeout := flag.Duration("timeout", 5*time.Second, "auction timeout")
	outDir := flag.String("out", "outputs", "output directory")
	seed := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout+1*time.Second)
	defer cancel()

	resMeta := resource.Capture(runtime.NumCPU())
	engine := auction.NewEngine(*numAuctions, *timeout)
	engine.Start(ctx, runtime.NumCPU())

	bidders := bidder.NewBidders(*numBidders, *seed)

	ids := make([]string, 0, *numAuctions)
	ids = append(ids, engine.AuctionList()...)

	bidder.Start(ctx, bidders, ids, engine)

	<-ctx.Done()
	if err := engine.GenerateReport(*outDir, resMeta); err != nil {
		log.Fatalf("report: %v", err)
	}
}
