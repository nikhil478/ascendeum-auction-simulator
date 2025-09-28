# Assumptions for Auction Simulator

### 1. Resource Standardization

* **Reference:** Requirement point 2 – *“Provide a mechanism to standardize resources with respect to the vCPU and RAM available.”*
* **Reference:** Evaluation Criteria point 3 – *“Resource Standardization: A clear approach to standardizing vCPUs and RAM usage.”*

**Simulation assumptions:**

* Standardization applies at the **Auction Service level only**:  
  * Limit goroutines per auction relative to CPU cores.  
  * Buffer channels or in-memory objects sized relative to available RAM.  
* System-level profiling or Docker-based resource limits are **not required**.  
* Goal is to maintain **predictable and controlled concurrency/memory usage**.

---

### 2. Bidder Behavior

**Simulation assumptions:**

* No rate limiting per bidder.  
* A bidder may submit multiple bids in rapid succession; this scenario is **not considered**.  
* Implementing a rate limiter is **not required**.

---

> I would love to have a conversation on each of these assumptions. If we decide something needs implementation, I have kept it **as simple and concise as possible** for this simulation.