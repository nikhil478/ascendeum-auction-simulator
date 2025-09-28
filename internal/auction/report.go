package auction

import (
	"encoding/json"
	"os"
	"time"

	"github.com/nikhil478/ascendeum-action-simulator/internal/resource"
)

type Summary struct {
	T0          time.Time        `json:"start"`
	TEnd        time.Time        `json:"end"`
	DurationSec float64          `json:"duration_seconds"`
	NumAuctions int             `json:"num_auctions"`
	Resource    resource.Metadata `json:"resource"`
	Auctions    []string         `json:"auctions"`
}

func NewSummary(t0, tEnd time.Time, n int, r resource.Metadata) *Summary {
	return &Summary{
		T0:          t0,
		TEnd:        tEnd,
		DurationSec: tEnd.Sub(t0).Seconds(),
		NumAuctions: n,
		Resource:    r,
	}
}

func (s *Summary) Add(a string) { s.Auctions = append(s.Auctions, a) }

func saveJSON(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}
