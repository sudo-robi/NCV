// Mock verification logs for demo purposes when backend is unavailable
export const mockLogs = [
    {
        "proof_type": "Head Correctness",
        "success": false,
        "message": "STATE ROOT MISMATCH at block 29578443",
        "evidence": {
            "block_number": 29578443,
            "ref_root": "0xc8d57c40eba3de124c16f7aeee3530e3685cc13b979ee948a757b65c790a1468",
            "nut_root": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
            "ref_height": 29578443,
            "nut_height": 29578444
        },
        "timestamp": 1768814306.787851
    },
    {
        "proof_type": "Freshness",
        "success": false,
        "message": "STALE DATA DETECTED. Node is 35 blocks behind (~210s)",
        "evidence": {
            "ref_height": 29578441,
            "nut_height": 29578406,
            "estimated_lag_seconds": 210,
            "threshold_seconds": 12
        },
        "timestamp": 1768814104.234567
    },
    {
        "proof_type": "Execution Correctness",
        "success": true,
        "message": "Execution results (fee) match: 159909413",
        "evidence": {
            "ref_fee": 159909413,
            "nut_fee": 159909413,
            "ref_weight": {
                "ref_time": 492762000,
                "proof_size": 10779
            },
            "nut_weight": {
                "ref_time": 492762000,
                "proof_size": 10779
            }
        },
        "timestamp": 1768814319.173694
    },
    {
        "proof_type": "Freshness",
        "success": true,
        "message": "Node is fresh. Node is 0 blocks behind (~0s)",
        "evidence": {
            "ref_height": 29578447,
            "nut_height": 29578447,
            "estimated_lag_seconds": 0,
            "threshold_seconds": 12
        },
        "timestamp": 1768814322.144780
    },
    {
        "proof_type": "Head Correctness",
        "success": true,
        "message": "State roots match at block 29578440",
        "evidence": {
            "block_number": 29578440,
            "ref_root": "0xa1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
            "nut_root": "0xa1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
            "ref_height": 29578440,
            "nut_height": 29578440
        },
        "timestamp": 1768814200.123456
    }
];
