# Ascendeum Auction Simulator

A concurrent auction simulator implemented in Go.

---

## Documentation

For detailed documentation, see:

- [Assumptions](docs/assumptions.md) – Simulation assumptions, auction timing, resource standardization, and bidder behavior.
- [Architecture](docs/architecture.md) – High-level design, services, and interaction flow.

---

## Run Locally (Go)

**Prerequisites:** Go 1.25.0+

```bash
# Run simulator
go run cmd/simulator/main.go

# Run with race detection
go run -race cmd/simulator/main.go

# Run tests
go test ./... -v
````

* Output files will be logged in the `outputs` folder at the project root.

---

## Run with Docker

**Prerequisites:** Docker

```bash
# Build Docker image
# Ensure there is no folder named 'outputs' in the root directory
docker build -t ascendeum-simulator .

# Run container
# Ensure there is an empty folder named 'outputs' in the root directory
docker run --rm -v $(pwd)/outputs:/outputs ascendeum-simulator
```

* Outputs will be logged in the `outputs` folder at the project root.

---

## Notes

* Command-line flags allow customizing auctions, bidders, timeout, and output directory.
* Example:

```bash
docker run --rm -v $(pwd)/outputs:/outputs ascendeum-simulator \
  -auctions 50 -bidders 200 -timeout 10s -out /outputs
```