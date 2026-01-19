import json
import http.server
import socketserver
import requests
import os

class SaboteurProxy(http.server.BaseHTTPRequestHandler):
    TARGET_RPC = "https://polkadot-rpc.publicnode.com"
    FORGE_HEIGHT = None
    FORGE_ROOT = "0x" + "f" * 64
    STUCK_HEADER = None

    def do_POST(self):
        is_dormant = os.environ.get("SABOTEUR_DORMANT") == "true"
        is_poisoned = os.environ.get("SABOTEUR_INCORRECT_DATA") == "true"
        
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        payload = json.loads(post_data)

        # Intercept and forge if necessary
        if payload.get('method') == 'chain_getHeader':
            if is_dormant and SaboteurProxy.STUCK_HEADER:
                print(f"[SABOTEUR] DORMANT MODE: Returning stuck header (Block {int(SaboteurProxy.STUCK_HEADER['number'], 16)})")
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps({"jsonrpc": "2.0", "result": SaboteurProxy.STUCK_HEADER, "id": payload.get('id')}).encode())
                return

        # Forward request to real RPC
        try:
            response = requests.post(self.TARGET_RPC, json=payload, timeout=10)
            resp_json = response.json()
        except Exception as e:
            print(f"[SABOTEUR] RPC ERROR: {e}")
            self.send_response(500)
            self.end_headers()
            return

        if payload.get('method') == 'chain_getHeader':
            header = resp_json.get('result', {})
            if header:
                height = int(header['number'], 16) if isinstance(header['number'], str) else header['number']
                
                if is_dormant and not SaboteurProxy.STUCK_HEADER:
                    # To guarantee a "stale" result, we'll return a header that is 30 blocks behind
                    behind_height = max(0, height - 30)
                    print(f"[SABOTEUR] DORMANT MODE ACTIVATED: Simulating node stuck at block {behind_height} (actual {height})")
                    # Modify the header copy
                    stale_header = json.loads(json.dumps(header))
                    stale_header['number'] = hex(behind_height)
                    # We should also change the hash to keep it consistent-ish, but NCV Freshness mostly cares about 'number'
                    SaboteurProxy.STUCK_HEADER = stale_header
                
                if is_poisoned:
                    if self.FORGE_HEIGHT and height == self.FORGE_HEIGHT:
                        print(f"[SABOTEUR] POISONING BLOCK {height}")
                        header['stateRoot'] = self.FORGE_ROOT
                    elif not self.FORGE_HEIGHT:
                        print(f"[SABOTEUR] POISONING HEAD (Block {height})")
                        header['stateRoot'] = self.FORGE_ROOT

        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(resp_json).encode())

def start_saboteur(port=9944, forge_height=None):
    SaboteurProxy.FORGE_HEIGHT = forge_height
    # Allow address reuse to avoid port conflicts
    socketserver.TCPServer.allow_reuse_address = True
    with socketserver.TCPServer(("", port), SaboteurProxy) as httpd:
        print(f"Saboteur active on port {port}")
        httpd.serve_forever()

if __name__ == "__main__":
    start_saboteur()
