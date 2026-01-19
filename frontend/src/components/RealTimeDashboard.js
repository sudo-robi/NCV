import React, { useEffect, useState } from 'react';
import './RealTimeDashboard.css';
import io from 'socket.io-client';

const socket = io('http://localhost:4000'); // Replace with your backend WebSocket URL

function RealTimeDashboard() {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    // Listen for real-time log updates
    socket.on('logUpdate', (newLog) => {
      setLogs((prevLogs) => [newLog, ...prevLogs]);
    });

    return () => {
      socket.disconnect();
    };
  }, []);

  return (
    <div className="real-time-dashboard">
      <h2>Real-Time Dashboard</h2>
      <ul className="log-list">
        {logs.map((log, index) => (
          <li key={index} className="log-item">
            <strong>{log.timestamp}</strong>: {log.message}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default RealTimeDashboard;