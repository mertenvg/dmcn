import { useParams, useNavigate } from 'react-router-dom';
import { useMessages } from '../lib/hooks/useMessages';

export function MessageView() {
  const { id } = useParams<{ id: string }>();
  const { messages, acknowledge } = useMessages();
  const navigate = useNavigate();
  const msg = messages.find(m => m.hash === id);

  if (!msg) {
    return (
      <div style={{ maxWidth: 600, margin: '40px auto', padding: 24 }}>
        <p>Message not found.</p>
        <button onClick={() => navigate('/inbox')} style={{ marginTop: 16, padding: '6px 12px', border: '1px solid #ccc', borderRadius: 4, background: 'white', cursor: 'pointer' }}>Back to Inbox</button>
      </div>
    );
  }

  const handleAck = async () => {
    await acknowledge(msg.hash);
    navigate('/inbox');
  };

  return (
    <div style={{ maxWidth: 600, margin: '0 auto', padding: 24 }}>
      <button onClick={() => navigate('/inbox')} style={{ marginBottom: 16, padding: '6px 12px', border: '1px solid #ccc', borderRadius: 4, background: 'white', cursor: 'pointer' }}>Back</button>
      <div style={{ background: 'white', padding: 24, borderRadius: 8, boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
        <h2 style={{ marginBottom: 8 }}>{msg.subject}</h2>
        <div style={{ color: '#666', marginBottom: 4 }}>From: {msg.senderAddress}</div>
        <div style={{ color: '#666', marginBottom: 16 }}>Date: {new Date(msg.sentAt * 1000).toLocaleString()}</div>
        <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.6 }}>{msg.body}</div>
      </div>
      <button onClick={handleAck} style={{ marginTop: 16, padding: '6px 12px', border: '1px solid #ccc', borderRadius: 4, background: 'white', cursor: 'pointer' }}>Acknowledge & Remove</button>
    </div>
  );
}
