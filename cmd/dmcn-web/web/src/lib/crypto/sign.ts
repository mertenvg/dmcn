import { importEd25519PrivateKey, importEd25519PublicKey } from './keys';

export async function sign(ed25519Seed: Uint8Array, data: Uint8Array): Promise<Uint8Array> {
  const key = await importEd25519PrivateKey(ed25519Seed);
  const sig = await crypto.subtle.sign('Ed25519', key, data);
  return new Uint8Array(sig);
}

export async function verify(ed25519Pub: Uint8Array, data: Uint8Array, signature: Uint8Array): Promise<boolean> {
  const key = await importEd25519PublicKey(ed25519Pub);
  return crypto.subtle.verify('Ed25519', key, signature, data);
}
