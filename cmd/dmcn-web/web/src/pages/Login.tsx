import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../lib/hooks/useAuth';
import { useKeys } from '../lib/hooks/useKeys';
import { login, loginVerify } from '../lib/api/client';
import { decryptKeys } from '../lib/crypto/keystore';
import { sign } from '../lib/crypto/sign';
import { fromBase64, toBase64 } from '../lib/crypto/keys';
import type { IdentityKeyPair } from '../lib/crypto/keys';

export function Login() {
  const [address, setAddress] = useState('');
  const [passphrase, setPassphrase] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const { setSession } = useAuth();
  const { setKeys } = useKeys();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      // Step 1: Get encrypted payload + challenge
      const loginResp = await login(address);

      // Step 2: Decrypt keys with passphrase
      const keyBytes = await decryptKeys(loginResp.encrypted_payload, passphrase);

      // Parse key bytes (JSON format matching the encrypted payload structure)
      const keyData = JSON.parse(new TextDecoder().decode(keyBytes));
      const keys: IdentityKeyPair = {
        ed25519Public: fromBase64(keyData.ed25519_public),
        ed25519Private: fromBase64(keyData.ed25519_private),
        x25519Public: fromBase64(keyData.x25519_public),
        x25519Private: fromBase64(keyData.x25519_private),
        deviceId: fromBase64(keyData.device_id),
        createdAt: keyData.created_at,
      };

      // Step 3: Sign challenge nonce
      const nonce = fromBase64(loginResp.challenge_nonce);
      // Ed25519 seed is first 32 bytes of the 64-byte private key
      const seed = keys.ed25519Private.slice(0, 32);
      const signature = await sign(seed, nonce);

      // Step 4: Verify with server
      const { session_token } = await loginVerify(
        address,
        toBase64(signature),
        loginResp.challenge_nonce
      );

      // Store encrypted payload in localStorage for future fast loads
      localStorage.setItem('dmcn_encrypted_payload', JSON.stringify({
        version: loginResp.version,
        encrypted_payload: loginResp.encrypted_payload,
      }));

      setKeys(keys);
      setSession(address, session_token);
      navigate('/inbox');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '80px auto', padding: 24 }}>
      <h1 style={{ marginBottom: 24 }}>DMCN Mail</h1>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 4, fontWeight: 600 }}>Address</label>
          <input
            type="text"
            value={address}
            onChange={e => setAddress(e.target.value)}
            placeholder="alice@dmcn.me"
            required
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc' }}
          />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 4, fontWeight: 600 }}>Passphrase</label>
          <input
            type="password"
            value={passphrase}
            onChange={e => setPassphrase(e.target.value)}
            required
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc' }}
          />
        </div>
        {error && <p style={{ color: 'red', marginBottom: 16 }}>{error}</p>}
        <button
          type="submit"
          disabled={loading}
          style={{ width: '100%', padding: 10, borderRadius: 4, background: '#0066cc', color: 'white', border: 'none', cursor: 'pointer', fontWeight: 600 }}
        >
          {loading ? 'Signing in...' : 'Sign In'}
        </button>
      </form>
      <p style={{ marginTop: 16, textAlign: 'center' }}>
        <Link to="/register">Create account</Link>
      </p>
    </div>
  );
}
