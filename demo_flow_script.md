# Demo Flow Script

This script outlines the steps to demonstrate the Node Correctness Verification (NCV) system.

## Prerequisites
- Python 3 installed
- Virtual environment set up
- Dependencies installed (`substrate-interface`, `requests`)
- Backend and frontend components ready

## Steps

### 1. Start the Saboteur Proxy
```bash
python saboteur.py
```
This proxy intercepts RPC calls and injects forged data to simulate failures.

### 2. Run the Challenger
```bash
python challenger.py <REF_URL> <NUT_URL>
```
Replace `<REF_URL>` and `<NUT_URL>` with the reference and node URLs.

### 3. View Logs
Check the `logs/` directory for JSON evidence files. Use the following command to view the latest log:
```bash
cat logs/evidence_<timestamp>.json
```

### 4. Validate Results
Ensure the challenger detects the injected failures and generates appropriate proofs. Look for the following in the logs:
- `Head Correctness`: State root mismatch detected.
- `Execution Correctness`: Fee estimation mismatch detected.
- `Freshness Proof`: Node lag detected.

### 5. Stop the Saboteur Proxy
```bash
pkill -f saboteur.py
```
Ensure the proxy is stopped to restore normal operation.

## Notes
- Use the `demo.py` script to automate the entire flow.
- Ensure the backend and frontend are running before starting the demo.