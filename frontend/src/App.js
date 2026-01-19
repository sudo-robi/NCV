import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import LogViewer from './components/LogViewer';

const LOG_URL = 'http://localhost:8081/logs';

// Helper to format timestamp as Day Hour:Min:Sec
export const formatTimestamp = (ts) => {
  const date = new Date(ts * 1000);
  const options = {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  };
  return date.toLocaleString('en-GB', options);
};

function App() {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const response = await axios.get(LOG_URL);
        setLogs(response.data);
      } catch (error) {
        console.error('Error fetching logs:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchLogs();
    const interval = setInterval(fetchLogs, 5000);
    return () => clearInterval(interval);
  }, []);

  const getStatusColor = (success) => success ? '#10b981' : '#ef4444';

  return (
    <div className="ncv-dashboard">
      <header>
        <h1>NCV Sentinel</h1>
        <div className="system-status">
          <span className="pulse"></span>
          System Live | Monitoring Active
        </div>
      </header>

      <main>
        <section className="evidence-feed">
          <LogViewer />

          <h2>Active Verification Stream</h2>
          {loading ? (
            <div className="loading">Initializing verification engine...</div>
          ) : logs.length === 0 ? (
            <div className="empty">No challenges recorded yet. Running...</div>
          ) : (
            <div className="log-list">
              {logs.map((log, index) => (
                <div key={index} className={`log-card ${log.success ? 'pass' : 'fail'}`}>
                  <div className="log-header">
                    <span className="proof-type">{log.proof_type}</span>
                    <span className="timestamp">{formatTimestamp(log.timestamp)}</span>
                  </div>
                  <div className="log-body">
                    <p className="message">{log.message}</p>
                    {log.evidence && (
                      <pre className="evidence-bundle">
                        {JSON.stringify(log.evidence, null, 2)}
                      </pre>
                    )}
                  </div>
                  <div className="log-footer">
                    {log.success ? '✓ VERIFIED' : '⚠ FAILURE DETECTED'}
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

export default App;