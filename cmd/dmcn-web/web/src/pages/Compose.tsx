import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../lib/hooks/useAuth';
import { useKeys } from '../lib/hooks/useKeys';
import { lookupIdentity, sendMessage } from '../lib/api/client';
import { encryptMessage } from '../lib/crypto/encrypt';
import { encodeSignedMessage, encodeEncryptedEnvelope, encodePlaintextMessage } from '../lib/crypto/protobuf';
import { sign } from '../lib/crypto/sign';
import { toBase64, fromBase64 } from '../lib/crypto/keys';

export function Compose() {
  const [to, setTo] = useState('');
  const [subject, setSubject] = useState('');
  const [body, setBody] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const { address } = useAuth();
  const { keys } = useKeys();
  const navigate = useNavigate();

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!keys || !address) return;
    setLoading(true);
    setError('');

    try {
      // Look up recipient
      const recipient = await lookupIdentity(to);
      const recipientX25519 = fromBase64(recipient.x25519_pub);

      // Create message UUID
      const messageId = crypto.getRandomValues(new Uint8Array(16));
      messageId[6] = (messageId[6] & 0x0f) | 0x40;
      messageId[8] = (messageId[8] & 0x3f) | 0x80;

      const threadId = crypto.getRandomValues(new Uint8Array(16));
      threadId[6] = (threadId[6] & 0x0f) | 0x40;
      threadId[8] = (threadId[8] & 0x3f) | 0x80;

      const now = Math.floor(Date.now() / 1000);

      // Encode plaintext message (for signing)
      const plaintextMsg = {
        version: 1,
        messageId,
        threadId,
        senderAddress: address,
        senderPublicKey: keys.ed25519Public,
        recipientAddress: to,
        sentAt: now,
        subject,
        body: { contentType: 'text/plain', content: new TextEncoder().encode(body) },
      };
      const plaintextBytes = await encodePlaintextMessage(plaintextMsg);

      // Sign
      const seed = keys.ed25519Private.slice(0, 32);
      const msgSignature = await sign(seed, plaintextBytes);

      // Encode signed message
      const signedMsgBytes = await encodeSignedMessage({
        plaintext: plaintextMsg,
        senderSignature: msgSignature,
      });

      // Encrypt for recipient (and self)
      const envelope = await encryptMessage(signedMsgBytes, messageId, now, [
        { deviceId: new Uint8Array(16), x25519Pub: recipientX25519 },
        { deviceId: keys.deviceId, x25519Pub: keys.x25519Public },
      ]);

      // Encode envelope to protobuf
      const envBytes = await encodeEncryptedEnvelope(envelope);

      // Compute SHA-256 of envelope bytes and sign it
      const envHashBuf = await crypto.subtle.digest('SHA-256', envBytes);
      const envHash = new Uint8Array(envHashBuf);
      const envSignature = await sign(seed, envHash);

      // Send
      await sendMessage({
        sender_address: address,
        sender_signature: toBase64(envSignature),
        envelope: toBase64(envBytes),
        recipient_address: to,
      });

      navigate('/inbox');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to send message');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 600, margin: '0 auto', padding: 24 }}>
      <h1 style={{ marginBottom: 24 }}>Compose</h1>
      <form onSubmit={handleSend}>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 4, fontWeight: 600 }}>To</label>
          <input type="text" value={to} onChange={e => setTo(e.target.value)} placeholder="bob@dmcn.me" required style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc' }} />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 4, fontWeight: 600 }}>Subject</label>
          <input type="text" value={subject} onChange={e => setSubject(e.target.value)} required style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc' }} />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 4, fontWeight: 600 }}>Message</label>
          <textarea value={body} onChange={e => setBody(e.target.value)} required rows={10} style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', resize: 'vertical' }} />
        </div>
        {error && <p style={{ color: 'red', marginBottom: 16 }}>{error}</p>}
        <div style={{ display: 'flex', gap: 12 }}>
          <button type="submit" disabled={loading} style={{ padding: '10px 24px', borderRadius: 4, background: '#0066cc', color: 'white', border: 'none', cursor: 'pointer', fontWeight: 600 }}>
            {loading ? 'Sending...' : 'Send'}
          </button>
          <button type="button" onClick={() => navigate('/inbox')} style={{ padding: '10px 24px', borderRadius: 4, background: '#eee', border: 'none', cursor: 'pointer' }}>Cancel</button>
        </div>
      </form>
    </div>
  );
}
