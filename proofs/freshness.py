from .base import BaseProof, ProofResult
import time

class FreshnessProof(BaseProof):
    def __init__(self, threshold_seconds=12): # Polkadot target is 6s, 12s is 2 blocks
        super().__init__("Freshness")
        self.threshold = threshold_seconds

    def run(self, ref_conn, nut_conn):
        try:
            ref_head = ref_conn.get_block_header()
            nut_head = nut_conn.get_block_header()

            ref_num = ref_head['header']['number']
            nut_num = nut_head['header']['number']

            lag = ref_num - nut_num
            # In Polkadot, 1 block = 6 seconds roughly
            # So lag of 1 = 6s, lag of 2 = 12s
            estimated_lag_seconds = lag * 6
            
            success = (estimated_lag_seconds < self.threshold)
            msg = f"Node is {lag} blocks behind (~{estimated_lag_seconds}s)"
            
            if success:
                msg = "Node is fresh. " + msg
            else:
                msg = "STALE DATA DETECTED. " + msg

            return ProofResult(
                self.name,
                success,
                msg,
                {
                    "ref_height": ref_num,
                    "nut_height": nut_num,
                    "estimated_lag_seconds": estimated_lag_seconds,
                    "threshold_seconds": self.threshold
                }
            )
        except Exception as e:
            return ProofResult(self.name, False, f"Error during freshness check: {str(e)}")
