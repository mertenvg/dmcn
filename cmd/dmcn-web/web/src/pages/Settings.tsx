import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../lib/hooks/useAuth';
import { useKeys } from '../lib/hooks/useKeys';
import { toHex } from '../lib/crypto/keys';

export function Settings() {
  const { address } = useAuth();
  const { keys } = useKeys();

  // Compute fingerprint: first 20 bytes of SHA-256(Ed25519Public || X25519Public)
  const [fingerprint, setFingerprint] = useState('');

  if (keys && !fingerprint) {
    const data = new Uint8Array(64);
    data.set(keys.ed25519Public, 0);
    data.set(keys.x25519Public, 32);
    crypto.subtle.digest('SHA-256', data).then(hash => {
      setFingerprint(toHex(new Uint8Array(hash).slice(0, 20)).toUpperCase());
    });
  }

  return (
    <div style={{ maxWidth: 600, margin: '0 auto', padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1>Settings</h1>
        <Link to="/inbox" style={{ padding: '6px 12px', background: '#eee', borderRadius: 4, textDecoration: 'none', color: '#333' }}>Back</Link>
      </div>

      <div style={{ background: 'white', padding: 24, borderRadius: 8 }}>
        <div style={{ marginBottom: 16 }}>
          <label style={{ fontWeight: 600, display: 'block', marginBottom: 4 }}>Address</label>
          <div style={{ fontFamily: 'monospace' }}>{address}</div>
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ fontWeight: 600, display: 'block', marginBottom: 4 }}>Fingerprint</label>
          <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>{fingerprint || 'Computing...'}</div>
        </div>
        <div>
          <label style={{ fontWeight: 600, display: 'block', marginBottom: 4 }}>Security</label>
          <p style={{ color: '#666', fontSize: 14 }}>Private keys are stored encrypted and never leave your browser. All encryption and signing happens client-side.</p>
        </div>
      </div>
    </div>
  );
}
