import { toBase64, fromBase64 } from './keys';

export interface EncryptedBundle {
  salt: string;       // base64
  nonce: string;      // base64 (12 bytes IV)
  ciphertext: string; // base64
  tag: string;        // base64 (16 bytes)
}

const INFO_STRING = new TextEncoder().encode('dmcn-webkeys-v1');

export async function encryptKeys(keyBytes: Uint8Array, passphrase: string): Promise<EncryptedBundle> {
  const salt = crypto.getRandomValues(new Uint8Array(32));

  // Import passphrase as key material
  const passphraseKey = await crypto.subtle.importKey(
    'raw',
    new TextEncoder().encode(passphrase),
    'HKDF',
    false,
    ['deriveKey']
  );

  // Derive AES key using HKDF-SHA256
  const aesKey = await crypto.subtle.deriveKey(
    { name: 'HKDF', hash: 'SHA-256', salt, info: INFO_STRING },
    passphraseKey,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt']
  );

  const nonce = crypto.getRandomValues(new Uint8Array(12));
  const encrypted = await crypto.subtle.encrypt(
    { name: 'AES-GCM', iv: nonce, tagLength: 128 },
    aesKey,
    keyBytes
  );

  // Web Crypto appends tag to ciphertext — split them
  const encryptedBytes = new Uint8Array(encrypted);
  const ciphertext = encryptedBytes.slice(0, encryptedBytes.length - 16);
  const tag = encryptedBytes.slice(encryptedBytes.length - 16);

  return {
    salt: toBase64(salt),
    nonce: toBase64(nonce),
    ciphertext: toBase64(ciphertext),
    tag: toBase64(tag),
  };
}

export async function decryptKeys(bundle: EncryptedBundle, passphrase: string): Promise<Uint8Array> {
  const salt = fromBase64(bundle.salt);
  const nonce = fromBase64(bundle.nonce);
  const ciphertext = fromBase64(bundle.ciphertext);
  const tag = fromBase64(bundle.tag);

  const passphraseKey = await crypto.subtle.importKey(
    'raw',
    new TextEncoder().encode(passphrase),
    'HKDF',
    false,
    ['deriveKey']
  );

  const aesKey = await crypto.subtle.deriveKey(
    { name: 'HKDF', hash: 'SHA-256', salt, info: INFO_STRING },
    passphraseKey,
    { name: 'AES-GCM', length: 256 },
    false,
    ['decrypt']
  );

  // Reassemble ciphertext + tag (Web Crypto expects them concatenated)
  const combined = new Uint8Array(ciphertext.length + tag.length);
  combined.set(ciphertext);
  combined.set(tag, ciphertext.length);

  const decrypted = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: nonce, tagLength: 128 },
    aesKey,
    combined
  );

  return new Uint8Array(decrypted);
}
