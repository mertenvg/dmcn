import { useState, useCallback, useEffect } from 'react';
import { listMessages, ackMessage, type EnvelopeItem } from '../api/client';
import { decryptEnvelope } from '../crypto/decrypt';
import { decodeEncryptedEnvelope, decodeSignedMessage } from '../crypto/protobuf';
import { fromBase64, fromHex } from '../crypto/keys';
import { useKeys } from './useKeys';

export interface DecryptedMessage {
  hash: string;
  senderAddress: string;
  recipientAddress: string;
  subject: string;
  body: string;
  sentAt: number;
}

export function useMessages() {
  const [messages, setMessages] = useState<DecryptedMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { keys } = useKeys();

  const fetchAndDecrypt = useCallback(async () => {
    if (!keys) return;
    setLoading(true);
    setError(null);
    try {
      const { envelopes } = await listMessages();
      const decrypted: DecryptedMessage[] = [];

      for (const item of envelopes) {
        try {
          const envBytes = fromBase64(item.data);
          const envProto = await decodeEncryptedEnvelope(envBytes);

          // Convert to decrypt-compatible format
          const envelope = {
            recipients: envProto.recipients.map(r => ({
              recipientXPub: new Uint8Array(r.recipientXPub),
              ephemeralXPub: new Uint8Array(r.ephemeralXPub),
              wrappedCek: new Uint8Array(r.wrappedCek),
              cekNonce: new Uint8Array(r.cekNonce),
              cekTag: new Uint8Array(r.cekTag),
            })),
            encryptedPayload: new Uint8Array(envProto.encryptedPayload),
            payloadNonce: new Uint8Array(envProto.payloadNonce),
            payloadTag: new Uint8Array(envProto.payloadTag),
          };

          const signedMsgBytes = await decryptEnvelope(envelope, keys.x25519Private, keys.x25519Public);
          const signedMsg = await decodeSignedMessage(signedMsgBytes);
          const pt = signedMsg.plaintext;

          decrypted.push({
            hash: item.hash,
            senderAddress: pt.senderAddress,
            recipientAddress: pt.recipientAddress,
            subject: pt.subject,
            body: new TextDecoder().decode(pt.body.content),
            sentAt: Number(pt.sentAt),
          });
        } catch (e) {
          console.error('Failed to decrypt envelope:', item.hash, e);
        }
      }

      setMessages(decrypted);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to fetch messages');
    } finally {
      setLoading(false);
    }
  }, [keys]);

  const acknowledge = useCallback(async (hash: string) => {
    await ackMessage(hash);
    setMessages(prev => prev.filter(m => m.hash !== hash));
  }, []);

  return { messages, loading, error, fetchAndDecrypt, acknowledge };
}
