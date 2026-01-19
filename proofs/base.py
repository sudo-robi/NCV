from substrateinterface import SubstrateInterface
import time
import logging

class ProofResult:
    def __init__(self, proof_type, success, message, evidence=None):
        self.proof_type = proof_type
        self.success = success
        self.message = message
        self.evidence = evidence or {}
        self.timestamp = time.time()

    def to_dict(self):
        return {
            "proof_type": self.proof_type,
            "success": self.success,
            "message": self.message,
            "evidence": self.evidence,
            "timestamp": self.timestamp
        }

class BaseProof:
    def __init__(self, name):
        self.name = name

    def run(self, ref_conn: SubstrateInterface, nut_conn: SubstrateInterface):
        raise NotImplementedError
