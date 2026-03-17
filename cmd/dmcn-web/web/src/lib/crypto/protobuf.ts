import protobuf from 'protobufjs';

let root: protobuf.Root | null = null;

// Load proto definitions defined inline for PoC.
// protobufjs encodes fields in field number order by default,
// which matches Go's deterministic encoding.
async function getRoot(): Promise<protobuf.Root> {
  if (root) return root;

  root = protobuf.Root.fromJSON({
    nested: {
      dmcn: {
        nested: {
          identity: {
            nested: {
              VerificationTier: {
                values: { UNSPECIFIED: 0, PROVIDER_HOSTED: 1, DOMAIN_DNS: 2, DANE: 3 },
              },
              IdentityRecord: {
                fields: {
                  version: { type: 'uint32', id: 1 },
                  address: { type: 'string', id: 2 },
                  ed25519PublicKey: { type: 'bytes', id: 3 },
                  x25519PublicKey: { type: 'bytes', id: 4 },
                  createdAt: { type: 'int64', id: 5 },
                  expiresAt: { type: 'int64', id: 6 },
                  relayHints: { type: 'string', id: 7, rule: 'repeated' },
                  verificationTier: { type: 'VerificationTier', id: 8 },
                  selfSignature: { type: 'bytes', id: 10 },
                  bridgeCapability: { type: 'bool', id: 11 },
                },
              },
            },
          },
          message: {
            nested: {
              MessageBody: {
                fields: {
                  contentType: { type: 'string', id: 1 },
                  content: { type: 'bytes', id: 2 },
                },
              },
              AttachmentRecord: {
                fields: {
                  attachmentId: { type: 'bytes', id: 1 },
                  filename: { type: 'string', id: 2 },
                  contentType: { type: 'string', id: 3 },
                  sizeBytes: { type: 'uint64', id: 4 },
                  contentHash: { type: 'bytes', id: 5 },
                  content: { type: 'bytes', id: 6 },
                },
              },
              PlaintextMessage: {
                fields: {
                  version: { type: 'uint32', id: 1 },
                  messageId: { type: 'bytes', id: 2 },
                  threadId: { type: 'bytes', id: 3 },
                  senderAddress: { type: 'string', id: 4 },
                  senderPublicKey: { type: 'bytes', id: 5 },
                  recipientAddress: { type: 'string', id: 6 },
                  sentAt: { type: 'int64', id: 7 },
                  subject: { type: 'string', id: 8 },
                  body: { type: 'MessageBody', id: 9 },
                  attachments: { type: 'AttachmentRecord', id: 10, rule: 'repeated' },
                  replyToId: { type: 'bytes', id: 11 },
                },
              },
              SignedMessage: {
                fields: {
                  plaintext: { type: 'PlaintextMessage', id: 1 },
                  senderSignature: { type: 'bytes', id: 2 },
                },
              },
              RecipientRecord: {
                fields: {
                  deviceId: { type: 'bytes', id: 1 },
                  recipientXPub: { type: 'bytes', id: 2 },
                  ephemeralXPub: { type: 'bytes', id: 3 },
                  wrappedCek: { type: 'bytes', id: 4 },
                  cekNonce: { type: 'bytes', id: 5 },
                  cekTag: { type: 'bytes', id: 6 },
                },
              },
              EncryptedEnvelope: {
                fields: {
                  version: { type: 'uint32', id: 1 },
                  messageId: { type: 'bytes', id: 2 },
                  recipients: { type: 'RecipientRecord', id: 3, rule: 'repeated' },
                  encryptedPayload: { type: 'bytes', id: 4 },
                  payloadNonce: { type: 'bytes', id: 5 },
                  payloadTag: { type: 'bytes', id: 6 },
                  payloadSizeClass: { type: 'uint32', id: 7 },
                  createdAt: { type: 'int64', id: 8 },
                  ratchetPubKey: { type: 'bytes', id: 9 },
                },
              },
            },
          },
        },
      },
    },
  });

  return root;
}

export async function encodeIdentityRecord(record: {
  version: number;
  address: string;
  ed25519PublicKey: Uint8Array;
  x25519PublicKey: Uint8Array;
  createdAt: number;
  expiresAt: number;
  relayHints: string[];
  verificationTier: number;
  bridgeCapability: boolean;
  selfSignature?: Uint8Array;
}): Promise<Uint8Array> {
  const root = await getRoot();
  const IdentityRecord = root.lookupType('dmcn.identity.IdentityRecord');
  const msg = IdentityRecord.create(record);
  return IdentityRecord.encode(msg).finish();
}

// Encode identity record WITHOUT selfSignature (for signing)
export async function encodeIdentitySignableBytes(record: {
  version: number;
  address: string;
  ed25519PublicKey: Uint8Array;
  x25519PublicKey: Uint8Array;
  createdAt: number;
  expiresAt: number;
  relayHints: string[];
  verificationTier: number;
  bridgeCapability: boolean;
}): Promise<Uint8Array> {
  const root = await getRoot();
  const IdentityRecord = root.lookupType('dmcn.identity.IdentityRecord');
  const msg = IdentityRecord.create({
    ...record,
    selfSignature: undefined,
  });
  return IdentityRecord.encode(msg).finish();
}

export async function encodePlaintextMessage(msg: {
  version: number;
  messageId: Uint8Array;
  threadId: Uint8Array;
  senderAddress: string;
  senderPublicKey: Uint8Array;
  recipientAddress: string;
  sentAt: number;
  subject: string;
  body: { contentType: string; content: Uint8Array };
  replyToId?: Uint8Array;
}): Promise<Uint8Array> {
  const root = await getRoot();
  const PlaintextMessage = root.lookupType('dmcn.message.PlaintextMessage');
  const encoded = PlaintextMessage.create(msg);
  return PlaintextMessage.encode(encoded).finish();
}

export async function encodeSignedMessage(msg: {
  plaintext: {
    version: number;
    messageId: Uint8Array;
    threadId: Uint8Array;
    senderAddress: string;
    senderPublicKey: Uint8Array;
    recipientAddress: string;
    sentAt: number;
    subject: string;
    body: { contentType: string; content: Uint8Array };
    replyToId?: Uint8Array;
  };
  senderSignature: Uint8Array;
}): Promise<Uint8Array> {
  const root = await getRoot();
  const SignedMessage = root.lookupType('dmcn.message.SignedMessage');
  const encoded = SignedMessage.create(msg);
  return SignedMessage.encode(encoded).finish();
}

export async function encodeEncryptedEnvelope(env: {
  version: number;
  messageId: Uint8Array;
  recipients: Array<{
    deviceId: Uint8Array;
    recipientXPub: Uint8Array;
    ephemeralXPub: Uint8Array;
    wrappedCek: Uint8Array;
    cekNonce: Uint8Array;
    cekTag: Uint8Array;
  }>;
  encryptedPayload: Uint8Array;
  payloadNonce: Uint8Array;
  payloadTag: Uint8Array;
  payloadSizeClass: number;
  createdAt: number;
  ratchetPubKey: Uint8Array;
}): Promise<Uint8Array> {
  const root = await getRoot();
  const EncryptedEnvelope = root.lookupType('dmcn.message.EncryptedEnvelope');
  const encoded = EncryptedEnvelope.create(env);
  return EncryptedEnvelope.encode(encoded).finish();
}

export async function decodeSignedMessage(data: Uint8Array): Promise<{
  plaintext: {
    version: number;
    messageId: Uint8Array;
    threadId: Uint8Array;
    senderAddress: string;
    senderPublicKey: Uint8Array;
    recipientAddress: string;
    sentAt: number;
    subject: string;
    body: { contentType: string; content: Uint8Array };
  };
  senderSignature: Uint8Array;
}> {
  const root = await getRoot();
  const SignedMessage = root.lookupType('dmcn.message.SignedMessage');
  const decoded = SignedMessage.decode(data);
  return decoded as any;
}

export async function decodeEncryptedEnvelope(data: Uint8Array): Promise<{
  version: number;
  messageId: Uint8Array;
  recipients: Array<{
    deviceId: Uint8Array;
    recipientXPub: Uint8Array;
    ephemeralXPub: Uint8Array;
    wrappedCek: Uint8Array;
    cekNonce: Uint8Array;
    cekTag: Uint8Array;
  }>;
  encryptedPayload: Uint8Array;
  payloadNonce: Uint8Array;
  payloadTag: Uint8Array;
  payloadSizeClass: number;
  createdAt: number;
  ratchetPubKey: Uint8Array;
}> {
  const root = await getRoot();
  const EncryptedEnvelope = root.lookupType('dmcn.message.EncryptedEnvelope');
  const decoded = EncryptedEnvelope.decode(data);
  return decoded as any;
}
