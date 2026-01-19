# Node Correctness Verification (NCV)

## Problem Statement and Motivation

Blockchain nodes can fail silently, continuing to respond to RPC calls and passing basic health checks while serving stale or incorrect data. This poses significant risks, including:
- **Data Integrity Issues**: Nodes may serve misleading or outdated data.
- **Network Vulnerabilities**: Lagging nodes can disrupt consensus mechanisms.
- **Operational Risks**: Undetected failures can lead to slashing penalties for validators.

**Motivation**: Existing monitoring tools focus on passive metrics like uptime and resource usage, which fail to detect these silent correctness issues. There is a need for an active verification system to ensure nodes are not just online but also correct and trustworthy.

---

## Explanation of the Solution

**Node Correctness Verification (NCV)** actively challenges blockchain nodes with verifiable proofs to ensure correctness, execution accuracy, and freshness. The system:
1. **Challenges Nodes**: Sends predefined challenges to the node under test (NUT).
2. **Compares Results**: Verifies the NUT's responses against trusted reference nodes.
3. **Generates Evidence**: Logs JSON-based proofs for debugging and slashing.

### How It Works
- **Head Correctness**: Verifies the `state_root` at a specific block height against multiple reference nodes.
- **Execution Correctness**: Simulates a transaction and verifies the computed fees and weights.
- **Freshness Proof**: Monitors block arrival times to detect stalling or lagging nodes.
- **Runtime Verification**: Ensures the node is running the correct WASM runtime version.

---

## Technologies, Tools, and Frameworks Used

- **Backend**: Go (Golang) for high-performance proof logic and modular design.
- **Frontend**: React for the log viewer and real-time dashboard.
- **APIs**: Polkadot RPC API for blockchain interactions.
- **Testing**: Unit and integration tests using Go's testing package.
- **Other Tools**:
  - `saboteur.py`: Simulates node failures for testing.
  - `demo.py`: Demonstrates the system's functionality end-to-end.

---

## Team Member Roles and Contributions

This project was developed by a **solo developer**. Contributions include:
- **Backend Development**: Designed and implemented modular proof logic, advanced proofs, and multi-chain support.
- **Frontend Development**: Built the log viewer and real-time dashboard for monitoring.
- **Testing and Validation**: Wrote unit and integration tests, simulated failure scenarios, and validated the system end-to-end.
- **Documentation**: Authored the README, demo flow script, and inline code comments.

---

## Future Scope and Roadmap

### Short-Term Goals
- **Multi-Chain Support**: Extend NCV to support Ethereum, Solana, and other blockchains.
- **Advanced Proofs**: Implement trustless storage proofs and Byzantine behavior detection.
- **Frontend Enhancements**: Add real-time updates and advanced filtering to the log viewer.

### Medium-Term Goals
- **Dashboard**: Build a comprehensive dashboard for real-time monitoring and visualization.
- **Integration with Validators**: Collaborate with blockchain validators to integrate NCV into their infrastructure.

### Long-Term Goals
- **Decentralized Proofs**: Enable decentralized verification using light clients and oracles.
- **Community Adoption**: Promote NCV as a standard for node health verification across multiple blockchains.

## Core Features (The Hackathon MVP)
Implemented for **Polkadot**:
1. **Head Correctness**: Compares `state_root` at a specific block height against a trusted reference.
2. **Execution Correctness**: Simulates a transaction (Balance transfer) and verifies the node computes the correct result (Weight/Fees).
3. **Freshness Proof**: Monitors block arrival times to detect silent stalling.

## How It Works

The NCV system operates by sending active challenges to the node under test (NUT) and comparing its responses to trusted reference nodes. Each challenge is designed to detect specific failure modes:

1. **Head Correctness**: Verifies the `state_root` at a specific block height against multiple reference nodes to ensure consensus.
2. **Execution Correctness**: Simulates a transaction and verifies the computed fees and weights match the expected values.
3. **Freshness Proof**: Monitors block arrival times to detect stalling or lagging nodes.

Each proof generates a JSON evidence log, which can be used for debugging or slashing.

## Testing

### Unit Tests
Unit tests validate the core proof logic, ensuring each function behaves as expected. Tests are located in the `tests/` directory and can be run using:
```bash
go test ./tests/backend_proofs_test.go
```

### Integration Tests
Integration tests validate the interaction between the backend and frontend. These tests use mock servers to simulate real-world scenarios. Run them with:
```bash
go test ./tests/integration_test.go
```

### Failure Simulation
The `demo.py` script simulates node failures using the `saboteur.py` proxy. This ensures the system can detect and respond to silent failures.

## Originality and Comparison with Existing Solutions

### Originality
- **Active Verification**: Unlike traditional monitoring tools that passively observe metrics like CPU usage or peer count, NCV actively challenges nodes with verifiable proofs.
- **Proof-Based Evidence**: Each failure generates a JSON proof bundle, providing actionable insights for debugging and slashing.
- **Silent Failure Detection**: NCV detects subtle issues like state root mismatches and execution errors that other tools might miss.

### Comparison with Existing Solutions
- **Traditional Monitoring Tools**:
  - Focus on passive metrics (e.g., uptime, resource usage).
  - Cannot detect silent failures or correctness issues.
- **NCV**:
  - Actively verifies node correctness, execution, and freshness.
  - Provides detailed evidence logs for every challenge.
- **Light Clients**:
  - Verify blockchain data but lack the ability to challenge full nodes.
  - NCV complements light clients by ensuring full node correctness.

## Roadmap

### Short-Term Goals
- **Multi-Chain Support**: Extend NCV to support Ethereum, Solana, and other blockchains.
- **Advanced Proofs**: Implement runtime verification and trustless storage proofs.
- **Frontend Enhancements**: Add real-time updates and advanced filtering to the log viewer.

### Medium-Term Goals
- **Dashboard**: Build a comprehensive dashboard for real-time monitoring and visualization.
- **Integration with Validators**: Collaborate with blockchain validators to integrate NCV into their infrastructure.

### Long-Term Goals
- **Decentralized Proofs**: Enable decentralized verification using light clients and oracles.
- **Community Adoption**: Promote NCV as a standard for node health verification across multiple blockchains.

## Real-World Examples

### Example 1: Polkadot Validator
A Polkadot validator uses NCV to ensure their node is always in sync and serving correct data. During a network partition, NCV detects a state root mismatch and alerts the operator, preventing potential slashing.

### Example 2: Blockchain Explorer
A blockchain explorer integrates NCV to verify the correctness of its backend nodes. This ensures users receive accurate data, even during network disruptions.

### Example 3: Enterprise Blockchain
An enterprise blockchain network adopts NCV to monitor its nodes. NCV's JSON evidence logs provide a clear audit trail for compliance and debugging.

## Future Work

- **Multi-Chain Support**: Extend the framework to support other blockchains like Ethereum and Solana.
- **Advanced Proofs**: Add new challenges for runtime verification and trustless storage.
- **Dashboard**: Build a real-time dashboard for monitoring node health.

## 🎥 Demo Flow

1. **Start the Saboteur Proxy**:
   ```bash
   python saboteur.py
   ```
   This proxy intercepts RPC calls and injects forged data to simulate failures.

2. **Run the Challenger**:
   ```bash
   python challenger.py <REF_URL> <NUT_URL>
   ```
   Replace `<REF_URL>` and `<NUT_URL>` with the reference and node URLs.

3. **View Logs**:
   Check the `logs/` directory for JSON evidence files.

4. **Validate Results**:
   Ensure the challenger detects the injected failures and generates appropriate proofs.

## Project Structure
- `challenger.py`: The main service that runs the challenge cycles.
- `proofs/`: Modular proof implementations.
- `saboteur.py`: A malicious proxy used for demoing "Silent Failure" detection.
- `demo.py`: Integrated script to run the full verification lifecycle.
- `logs/`: Evidence logs and JSON proofs.

## Setup and Run Instructions

### Prerequisites
- **Python 3**: Ensure Python 3 is installed on your system.
- **Go (Golang)**: Install Go for backend development.
- **Node.js**: Required for the frontend.

### 1. Backend Setup
1. Install Go dependencies:
   ```bash
   go mod tidy
   ```
2. Run the backend:
   ```bash
   go run challenger.go <nut_url> <ref_url1> [ref_url2...]
   ```
   Replace `<nut_url>` with the Node Under Test (NUT) URL and `<ref_url1>` with reference node URLs.

### 2. Frontend Setup
1. Navigate to the `frontend/` directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the frontend:
   ```bash
   npm start
   ```
   The frontend will be available at `http://localhost:3000`.

### 3. Demo Script
1. Activate the Python virtual environment:
   ```bash
   source venv/bin/activate
   ```
2. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   ```
3. Run the demo:
   ```bash
   python demo.py
   ```
   This will simulate node failures and demonstrate the system's functionality.

### 4. Testing
- **Unit Tests**:
  ```bash
  go test ./tests/backend_proofs_test.go
  ```
- **Integration Tests**:
  ```bash
  go test ./tests/integration_test.go
  ```

### Failure Simulation
The `demo.py` script simulates node failures using the `saboteur.py` proxy. This ensures the system can detect and respond to silent failures.

## 🏁 Quick Start / Demo

### 1. Setup
```bash
python3 -m venv venv
source venv/bin/activate
pip install substrate-interface requests
```

### 2. Run the "Oh Shit" Demo
This demo starts a **Saboteur Proxy** that intercepts Polkadot RPC calls and injects forged `state_root` data, simulating a silent correctness failure.

```bash
python demo.py
```

**Watch the Challenger detect the corruption instantly while the node process remains "Online".**

## Why this Wins
- **Active Verification**: Most monitoring is passive; PoH is proactive.
- **Evidence-Based**: Every failure generates a JSON proof bundle for slaving/slashing evidence.
- **Scalable**: Framework designed to support multiple chains and failure modes.

---
Built for the Hackathon. Build depth, not breadth.
