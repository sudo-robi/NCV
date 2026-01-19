from .base import BaseProof, ProofResult
from substrateinterface import Keypair

class ExecutionCorrectnessProof(BaseProof):
    def __init__(self):
        super().__init__("Execution Correctness")
        # Alice on Substrate (Polkadot address format handled by library)
        self.kp = Keypair.create_from_uri('//Alice')

    def run(self, ref_conn, nut_conn):
        try:
            # Use a valid SS58 address that matches the chain's format
            # For Polkadot/Substrate, Alice is a good default for testing if we are on a dev node,
            # but for public RPC, we just need ANY valid address.
            dest_address = self.kp.ss58_address
            
            call = ref_conn.compose_call(
                call_module='Balances',
                call_function='transfer_keep_alive',
                call_params={
                    'dest': dest_address,
                    'value': 10**10 
                }
            )

            # query_info is a state-less check that verifies if the node can correctly 
            # compute the weight and fee for a given extrinsic.
            ref_info = ref_conn.get_payment_info(call=call, keypair=self.kp)
            nut_info = nut_conn.get_payment_info(call=call, keypair=self.kp)

            ref_fee = int(ref_info.get('partialFee', 0))
            nut_fee = int(nut_info.get('partialFee', 0))

            # Allow for some minor drift if the nodes are at different heights (fees change with congestion)
            # but ideally they should be very close.
            success = (ref_fee == nut_fee)
            msg = f"Execution results (fee) match: {ref_fee}" if success else f"EXECUTION MISMATCH: Ref {ref_fee}, Nut {nut_fee}"

            return ProofResult(
                self.name,
                success,
                msg,
                {
                    "ref_fee": ref_fee,
                    "nut_fee": nut_fee,
                    "ref_weight": ref_info.get('weight'),
                    "nut_weight": nut_info.get('weight')
                }
            )
        except Exception as e:
            return ProofResult(self.name, False, f"Error during execution simulation: {str(e)}")
