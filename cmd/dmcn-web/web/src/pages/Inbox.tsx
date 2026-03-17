import { useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useMessages } from '../lib/hooks/useMessages';
import { useAuth } from '../lib/hooks/useAuth';
import { logout as apiLogout } from '../lib/api/client';

export function Inbox() {
  const { messages, loading, error, fetchAndDecrypt } = useMessages();
  const { address, clearSession } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    fetchAndDecrypt();
  }, [fetchAndDecrypt]);

  const handleLogout = async () => {
    try { await apiLogout(); } catch { /* ignore */ }
    clearSession();
    navigate('/login');
  };

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1>Inbox</h1>
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <span style={{ color: '#666', fontSize: 14 }}>{address}</span>
          <Link to="/compose" style={{ padding: '6px 12px', background: '#0066cc', color: 'white', borderRadius: 4, textDecoration: 'none' }}>Compose</Link>
          <Link to="/contacts" style={{ padding: '6px 12px', background: '#eee', borderRadius: 4, textDecoration: 'none', color: '#333' }}>Contacts</Link>
          <Link to="/settings" style={{ padding: '6px 12px', background: '#eee', borderRadius: 4, textDecoration: 'none', color: '#333' }}>Settings</Link>
          <button onClick={handleLogout} style={{ padding: '6px 12px', background: '#eee', border: 'none', borderRadius: 4, cursor: 'pointer' }}>Logout</button>
        </div>
      </div>

      <button onClick={fetchAndDecrypt} disabled={loading} style={{ marginBottom: 16, padding: '6px 12px', border: '1px solid #ccc', borderRadius: 4, background: 'white', cursor: 'pointer' }}>
        {loading ? 'Refreshing...' : 'Refresh'}
      </button>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {messages.length === 0 && !loading && <p style={{ color: '#666', textAlign: 'center', padding: 40 }}>No messages</p>}

      <div>
        {messages.map(msg => (
          <div key={msg.hash} style={{ padding: 16, borderBottom: '1px solid #eee', cursor: 'pointer' }} onClick={() => navigate(`/message/${msg.hash}`)}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <strong>{msg.senderAddress}</strong>
              <span style={{ color: '#666', fontSize: 12 }}>{new Date(msg.sentAt * 1000).toLocaleString()}</span>
            </div>
            <div style={{ fontWeight: 600, marginTop: 4 }}>{msg.subject}</div>
            <div style={{ color: '#666', marginTop: 4, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{msg.body}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
