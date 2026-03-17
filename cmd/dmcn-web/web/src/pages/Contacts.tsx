import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useContacts } from '../lib/hooks/useContacts';
import { lookupIdentity } from '../lib/api/client';

export function Contacts() {
  const { contacts, addContact, removeContact } = useContacts();
  const [address, setAddress] = useState('');
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const identity = await lookupIdentity(address);
      addContact({ address, name: name || address, fingerprint: identity.fingerprint });
      setAddress('');
      setName('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add contact');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 600, margin: '0 auto', padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1>Contacts</h1>
        <Link to="/inbox" style={{ padding: '6px 12px', background: '#eee', borderRadius: 4, textDecoration: 'none', color: '#333' }}>Back</Link>
      </div>

      <form onSubmit={handleAdd} style={{ marginBottom: 24, padding: 16, background: 'white', borderRadius: 8 }}>
        <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
          <input type="text" value={address} onChange={e => setAddress(e.target.value)} placeholder="bob@dmcn.me" required style={{ flex: 1, padding: 8, borderRadius: 4, border: '1px solid #ccc' }} />
          <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Name (optional)" style={{ flex: 1, padding: 8, borderRadius: 4, border: '1px solid #ccc' }} />
          <button type="submit" disabled={loading} style={{ padding: '8px 16px', borderRadius: 4, background: '#0066cc', color: 'white', border: 'none', cursor: 'pointer' }}>
            {loading ? '...' : 'Add'}
          </button>
        </div>
        {error && <p style={{ color: 'red', fontSize: 14 }}>{error}</p>}
      </form>

      {contacts.length === 0 && <p style={{ color: '#666', textAlign: 'center' }}>No contacts</p>}

      {contacts.map(c => (
        <div key={c.address} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: 12, borderBottom: '1px solid #eee' }}>
          <div>
            <div style={{ fontWeight: 600 }}>{c.name}</div>
            <div style={{ color: '#666', fontSize: 14 }}>{c.address}</div>
            <div style={{ color: '#999', fontSize: 12, fontFamily: 'monospace' }}>{c.fingerprint}</div>
          </div>
          <button onClick={() => removeContact(c.address)} style={{ padding: '4px 8px', border: '1px solid #ccc', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 12 }}>Remove</button>
        </div>
      ))}
    </div>
  );
}
