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
