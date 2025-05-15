# ⚖️ Balancer

**Balancer** is a lightweight load balancing service written in Go. It supports multiple traffic distribution algorithms, runtime health checks, and request proxying to healthy backend services.

---

## 🚀 Features

- 🌀 Round Robin
- ⚖️ Weighted Round Robin
- 📉 Least Connections
- 🩺 Periodic health checks
- ⚙️ Configuration via environment variables
- ✅ Includes test coverage for core logic

---

## 📊 Load Balancing Algorithms

### 🔁 Round Robin

Cycles through the list of available hosts, forwarding each request to the next host in order.

### ⚖️ Weighted Round Robin

Distributes traffic proportionally based on predefined weights for each host. The algorithm ensures that no host is selected more often than its assigned weight allows.

### 📉 Least Connections

Sends traffic to the host with the fewest active connections — ideal for uneven or unpredictable workloads.

---

## ⚙️ Configuration

Balancer is configured via environment variables. You can use Docker Compose to run the service

### Example: `docker-compose.yml`

```yaml
version: '3.9'

services:
  balancer:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: balancer
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - HOST=0.0.0.0
      - PORT=8080
      - METHOD=round_robin           # Options: round_robin, weighted_round_robin, least_connections
      - SERVERS=http://host1:8123,http://host2:8123
      - WEIGHTS=1,2                  # Required for weighted round robin
      - HEALTHCHECK_INTERVAL=60     # Interval in seconds
      - HEALTHCHECK_TIMEOUT=15      # Timeout in seconds
      - WITH_LOG=true               # Enable or disable logging
