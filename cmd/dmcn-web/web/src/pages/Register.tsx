import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../lib/hooks/useAuth';
import { useKeys } from '../lib/hooks/useKeys';
import { register } from '../lib/api/client';
import { generateIdentityKeyPair, toBase64 } from '../lib/crypto/keys';
import { encryptKeys } from '../lib/crypto/keystore';
import { encodeIdentitySignableBytes, encodeIdentityRecord } from '../lib/crypto/protobuf';
import { sign } from '../lib/crypto/sign';

export function Register() {
  const [address, setAddress] = useState('');
  const [passphrase, setPassphrase] = useState('');
  const [confirmPassphrase, setConfirmPassphrase] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const { setSession } = useAuth();
  const { setKeys } = useKeys();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (passphrase !== confirmPassphrase) {
      setError('Passphrases do not match');
      return;
    }
    setLoading(true);
    setError('');

    try {
      // Generate key pairs
      const keys = await generateIdentityKeyPair();

      // Encrypt keys with passphrase
      const keyData = JSON.stringify({
        ed25519_public: toBase64(keys.ed25519Public),
        ed25519_private: toBase64(keys.ed25519Private),
        x25519_public: toBase64(keys.x25519Public),
        x25519_private: toBase64(keys.x25519Private),
        device_id: toBase64(keys.deviceId),
        created_at: keys.createdAt,
      });
      const encryptedPayload = await encryptKeys(new TextEncoder().encode(keyData), passphrase);

      // Create and sign identity record
      const now = keys.createdAt;
      const signableBytes = await encodeIdentitySignableBytes({
        version: 1,
        address,
        ed25519PublicKey: keys.ed25519Public,
        x25519PublicKey: keys.x25519Public,
        createdAt: now,
        expiresAt: 0,
        relayHints: [],
        verificationTier: 0,
        bridgeCapability: false,
      });

      const seed = keys.ed25519Private.slice(0, 32);
      const selfSignature = await sign(seed, signableBytes);

      // Encode full identity record with signature
      const identityRecordBytes = await encodeIdentityRecord({
        version: 1,
        address,
        ed25519PublicKey: keys.ed25519Public,
        x25519PublicKey: keys.x25519Public,
        createdAt: now,
        expiresAt: 0,
        relayHints: [],
        verificationTier: 0,
        bridgeCapability: false,
        selfSignature,
      });

      // Register with server
      const { session_token } = await register({
        address,
        ed25519_pub: toBase64(keys.ed25519Public),
        x25519_pub: toBase64(keys.x25519Public),
        encrypted_payload: encryptedPayload,
        identity_record: toBase64(identityRecordBytes),
        self_signature: toBase64(selfSignature),
      });

      setKeys(keys);
      setSession(address, session_token);
      navigate('/inbox');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '80px auto', padding: 24 }}>
      <h1 style={{ marginBottom: 24 }}>Create Account</h1>
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
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 4, fontWeight: 600 }}>Confirm Passphrase</label>
          <input
            type="password"
            value={confirmPassphrase}
            onChange={e => setConfirmPassphrase(e.target.value)}
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
          {loading ? 'Creating account...' : 'Create Account'}
        </button>
      </form>
      <p style={{ marginTop: 16, textAlign: 'center' }}>
        <Link to="/login">Already have an account?</Link>
      </p>
    </div>
  );
}
