## High-Level Architecture

### 1. Bidder Service
Simulates virtual bidders and their bidding behavior.

- Generates and submits bids based on auction attributes.
- Sends bids to a shared channel to communicate with the Auction Service.

---

### 2. Auction Service
Manages the lifecycle and state of each auction.

- Creates and initializes auctions with attributes and **individual deadlines**.
- Processes incoming bids **serially per auction** using mutex locks to ensure order.
- Supports **parallel processing across different auctions**.
- Ignores bids arriving after the auction deadline.
- Persists auction data and final results.

---

### 3. Simulation Service
Coordinates the auction simulation end-to-end.

- Parses runtime flags (number of auctions, bidders, timeouts, etc.).
- Requests the Auction Service to create auctions.
- Requests the Bidder Service to initialize bidders and provides the auction list.

---

## Interaction Flow

1. **Channel Setup** – Auction Service exposes a **bid channel** for receiving bids.
2. **Channel Distribution** – Simulation Service passes this channel to the Bidder Service.
3. **Bid Submission** – Multiple bidder goroutines submit bids concurrently:
   - Bids for the **same auction** are processed **serially** via a mutex.
   - Bids for **different auctions** are processed **in parallel**, leveraging multiple workers.
   - Bids arriving **after the auction deadline** are ignored.

> **Note:**  
> The channel simulates asynchronous network calls within a single Go process. It provides decoupled behavior similar to real HTTP or gRPC requests, without external infrastructure.
