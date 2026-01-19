import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './LogViewer.css';
import { formatTimestamp } from '../utils/time';
import { mockLogs } from '../mockData';

function LogViewer() {
  const [logs, setLogs] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredLogs, setFilteredLogs] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [logsPerPage] = useState(10);
  const [sortOrder, setSortOrder] = useState('asc');

  useEffect(() => {
    const LOG_URL = process.env.REACT_APP_LOG_SERVER_URL || 'http://localhost:8081/logs';
    axios.get(LOG_URL)
      .then(response => {
        setLogs(response.data);
        setFilteredLogs(response.data);
      })
      .catch(error => {
        console.error('Error fetching logs:', error);
        // Fallback to mock data if backend is unavailable
        console.log('Using mock data for demo purposes');
        setLogs(mockLogs);
        setFilteredLogs(mockLogs);
      });
  }, []);

  useEffect(() => {
    const filtered = logs.filter(log =>
      log.message.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredLogs(filtered);
  }, [searchTerm, logs]);

  const handleSort = () => {
    const sortedLogs = [...filteredLogs].sort((a, b) => {
      if (sortOrder === 'asc') {
        return a.timestamp - b.timestamp;
      } else {
        return b.timestamp - a.timestamp;
      }
    });
    setFilteredLogs(sortedLogs);
    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
  };

  const indexOfLastLog = currentPage * logsPerPage;
  const indexOfFirstLog = indexOfLastLog - logsPerPage;
  const currentLogs = filteredLogs.slice(indexOfFirstLog, indexOfLastLog);

  const paginate = (pageNumber) => setCurrentPage(pageNumber);

  return (
    <div className="log-viewer">
      <h2>Audit History</h2>
      <input
        type="text"
        placeholder="Search audit logs..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        className="search-bar"
      />
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <button onClick={handleSort} className="sort-button">
          Sort by Time ({sortOrder})
        </button>
        <span style={{ color: '#94a3b8', fontSize: '0.9rem' }}>Showing {filteredLogs.length} entries</span>
      </div>
      <ul className="log-list">
        {currentLogs.map((log, index) => (
          <li key={index} className="log-item">
            <strong>{formatTimestamp(log.timestamp)}</strong>
            <span>{log.message}</span>
          </li>
        ))}
      </ul>
      <div className="pagination">
        {Array.from({ length: Math.ceil(filteredLogs.length / logsPerPage) }, (_, i) => (
          <button
            key={i + 1}
            onClick={() => paginate(i + 1)}
            className={currentPage === i + 1 ? 'active' : ''}
          >
            {i + 1}
          </button>
        ))}
      </div>
    </div>
  );
}

export default LogViewer;