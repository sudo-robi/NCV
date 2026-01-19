from .base import BaseProof, ProofResult

class HeadCorrectnessProof(BaseProof):
    def __init__(self):
        super().__init__("Head Correctness")

    def run(self, ref_conn, nut_conn):
        try:
            ref_head = ref_conn.get_block_header()
            nut_head = nut_conn.get_block_header()

            ref_num = ref_head['header']['number']
            nut_num = nut_head['header']['number']
            
            # If heights are different, we try to get the header at the same height if possible,
            # but for a "live" check, being at different heights is itself a piece of evidence.
            # However, to check *correctness* specifically, we need to compare the same block.
            
            target_num = min(ref_num, nut_num)
            
            # Fetch headers at the target number
            ref_at_target = ref_conn.get_block_header(block_number=target_num)
            nut_at_target = nut_conn.get_block_header(block_number=target_num)
            
            ref_root = ref_at_target['header']['stateRoot']
            nut_root = nut_at_target['header']['stateRoot']
            
            success = (ref_root == nut_root)
            msg = f"State root match at block {target_num}" if success else f"STATE ROOT MISMATCH at block {target_num}"
            
            return ProofResult(
                self.name, 
                success, 
                msg, 
                {
                    "block_number": target_num,
                    "ref_root": ref_root,
                    "nut_root": nut_root,
                    "ref_height": ref_num,
                    "nut_height": nut_num
                }
            )
        except Exception as e:
            return ProofResult(self.name, False, f"Error: {str(e)}")
