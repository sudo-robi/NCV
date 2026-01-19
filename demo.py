import subprocess
import time
import os
import signal
import json
import requests

def wait_for_saboteur(port=9944, timeout=10):
    start_time = time.time()
    while time.time() - start_time < timeout:
        try:
            requests.get(f"http://localhost:{port}", timeout=1)
        except:
            time.sleep(1)
            continue
        return True
    return False

def run_demo():
    print("=== NODE CORRECTNESS VERIFICATION ===")
    
    scenarios = [
        {"name": "dormant", "env": {"SABOTEUR_DORMANT": "true", "SABOTEUR_INCORRECT_DATA": "false"}},
        {"name": "poisoned", "env": {"SABOTEUR_DORMANT": "false", "SABOTEUR_INCORRECT_DATA": "true"}}
    ]

    for scenario in scenarios:
        print(f"\n--- Scenario: {scenario['name'].upper()} ---")
        
        # Clear logs
        for f in os.listdir("logs"):
            if f.endswith(".json") or f.endswith(".log"):
                os.remove(os.path.join("logs", f))

        print(f"1. Starting Saboteur Proxy...")
        env = os.environ.copy()
        env.update(scenario['env'])
        
        saboteur_proc = subprocess.Popen(["./venv/bin/python", "saboteur.py"], env=env)
        time.sleep(3) # Wait for port to bind

        if scenario['name'] == "dormant":
            print("   (Priming dormant state...)")
            try:
                # Make a request to trigger the freeze
                requests.post("http://localhost:9944", 
                             json={"jsonrpc": "2.0", "method": "chain_getHeader", "params": [], "id": 1},
                             timeout=5)
            except:
                pass
            print("   (Waiting 5s for node to stabilize...)")
            time.sleep(5)

        try:
            print("2. Running Challenger...")
            REF = "https://polkadot-rpc.publicnode.com"
            NUT = "http://localhost:9944"
            
            challenger_env = os.environ.copy()
            challenger_env["PYTHONPATH"] = "."
            result = subprocess.run(["./venv/bin/python", "challenger.py", REF, NUT], 
                               capture_output=True, text=True, env=challenger_env)
        
            print("\nChallenger Output Preview:")
            print("\n".join(result.stdout.splitlines()[-5:])) # Show last 5 lines
        
            log_files = sorted([f for f in os.listdir("logs") if f.startswith("evidence_")], reverse=True)
            if log_files:
                latest_log = os.path.join("logs", log_files[0])
                with open(latest_log, "r") as f:
                    evidence = json.load(f)
                    
                freshness = next((p for p in evidence if p['proof_type'] == "Freshness"), None)
                correctness = next((p for p in evidence if p['proof_type'] == "Head Correctness"), None)
                
                if scenario['name'] == "dormant":
                    if freshness and not freshness['success']:
                        print(f"\n✅ SUCCESS: NCV detected the DORMANT node (Staleness proof)!")
                        print(f"   Message: {freshness['message']}")
                    else:
                        print(f"\n❌ FAILURE: NCV did not detect the staleness.")
                        if freshness: print(f"   Debug: {freshness['message']}")
                else:
                    if correctness and not correctness['success']:
                        print(f"\n✅ SUCCESS: NCV detected the POISONED node (State Root Mismatch)!")
                    else:
                        print(f"\n❌ FAILURE: NCV did not detect the corruption.")

        finally:
            print("3. Cleaning up Saboteur...")
            saboteur_proc.terminate()
            saboteur_proc.wait()
            time.sleep(2) # Give port time to release

if __name__ == "__main__":
    run_demo()
