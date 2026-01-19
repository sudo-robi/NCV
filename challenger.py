import json
import time
import logging
from substrateinterface import SubstrateInterface
from proofs.head_correctness import HeadCorrectnessProof
from proofs.execution_correctness import ExecutionCorrectnessProof
from proofs.freshness import FreshnessProof

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler("logs/challenger.log"),
        logging.StreamHandler()
    ]
)

class Challenger:
    def __init__(self, ref_url, nut_url):
        self.ref_url = ref_url
        self.nut_url = nut_url
        self.proofs = [
            HeadCorrectnessProof(),
            ExecutionCorrectnessProof(),
            FreshnessProof()
        ]
        self.evidence_log = []

    def run_cycle(self):
        logging.info(f"Starting verification cycle. NUT: {self.nut_url}")
        
        try:
            ref_conn = SubstrateInterface(url=self.ref_url)
            nut_conn = SubstrateInterface(url=self.nut_url)
        except Exception as e:
            logging.error(f"Failed to connect to nodes: {e}")
            return

        cycle_results = []
        for proof in self.proofs:
            logging.info(f"Running proof: {proof.name}...")
            result = proof.run(ref_conn, nut_conn)
            
            if not result.success:
                logging.warning(f"PROOF FAILED: {result.proof_type} - {result.message}")
            else:
                logging.info(f"Proof success: {result.proof_type}")
                
            cycle_results.append(result.to_dict())
        
        self.save_evidence(cycle_results)
        return cycle_results

    def save_evidence(self, results):
        timestamp = int(time.time())
        filename = f"logs/evidence_{timestamp}.json"
        with open(filename, "w") as f:
            json.dump(results, f, indent=4)
        logging.info(f"Evidence saved to {filename}")

if __name__ == "__main__":
    import sys
    # Allow passing URLs via CLI
    REF = sys.argv[1] if len(sys.argv) > 1 else "wss://rpc.polkadot.io"
    NUT = sys.argv[2] if len(sys.argv) > 2 else "wss://polkadot.api.onfinality.io/public-ws"
    
    challenger = Challenger(REF, NUT)
    challenger.run_cycle()
