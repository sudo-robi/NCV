import json
import os
import glob
from http.server import HTTPServer, BaseHTTPRequestHandler

class LogHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/logs':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.send_header('Access-Control-Allow-Origin', '*')
            self.end_headers()
            
            log_dir = 'logs'
            files = sorted(glob.glob(os.path.join(log_dir, 'evidence_*.json')), reverse=True)
            
            all_logs = []
            for file_path in files[:20]: 
                try:
                    with open(file_path, 'r') as f:
                        data = json.load(f)
                        if isinstance(data, list):
                            all_logs.extend(data)
                        else:
                            all_logs.append(data)
                except Exception as e:
                    print(f"Error reading {file_path}: {e}")
            
            all_logs.sort(key=lambda x: x.get('timestamp', 0), reverse=True)
            
            self.wfile.write(json.dumps(all_logs).encode())
        else:
            self.send_response(404)
            self.end_headers()

def run_server(port=8081):
    server_address = ('', port)
    httpd = HTTPServer(server_address, LogHandler)
    print(f"Log Server running on port {port}")
    httpd.serve_forever()

if __name__ == "__main__":
    run_server()
