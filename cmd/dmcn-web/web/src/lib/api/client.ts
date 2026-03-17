let sessionToken: string | null = null;

export function setSessionToken(token: string | null) {
  sessionToken = token;
}

export function getSessionToken(): string | null {
  return sessionToken;
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (sessionToken) {
    headers['Authorization'] = `Bearer ${sessionToken}`;
  }

  const res = await fetch(path, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

// Auth API
export interface EncryptedPayload {
  salt: string;
  nonce: string;
  ciphertext: string;
  tag: string;
}

export interface RegisterRequest {
  address: string;
  ed25519_pub: string;
  x25519_pub: string;
  encrypted_payload: EncryptedPayload;
  identity_record: string;
  self_signature: string;
}

export interface LoginResponse {
  version: number;
  ed25519_pub: string;
  encrypted_payload: EncryptedPayload;
  challenge_nonce: string;
}

export interface SessionResponse {
  session_token: string;
}

export function register(req: RegisterRequest): Promise<SessionResponse> {
  return request('POST', '/api/v1/register', req);
}

export function login(address: string): Promise<LoginResponse> {
  return request('POST', '/api/v1/login', { address });
}

export function loginVerify(address: string, challengeSignature: string, challengeNonce: string): Promise<SessionResponse> {
  return request('POST', '/api/v1/login/verify', {
    address,
    challenge_signature: challengeSignature,
    challenge_nonce: challengeNonce,
  });
}

export function logout(): Promise<void> {
  return request('POST', '/api/v1/logout');
}

// Identity API
export interface IdentityLookupResponse {
  address: string;
  ed25519_pub: string;
  x25519_pub: string;
  fingerprint: string;
  verification_tier: number;
}

export function lookupIdentity(address: string): Promise<IdentityLookupResponse> {
  return request('GET', `/api/v1/identity/lookup?address=${encodeURIComponent(address)}`);
}

// Messages API
export interface SendMessageRequest {
  sender_address: string;
  sender_signature: string;
  envelope: string;
}

export interface SendMessageResponse {
  envelope_hash: string;
}

export interface EnvelopeItem {
  hash: string;
  data: string;
}

export interface ListMessagesResponse {
  envelopes: EnvelopeItem[];
}

export function sendMessage(req: SendMessageRequest): Promise<SendMessageResponse> {
  return request('POST', '/api/v1/messages/send', req);
}

export function listMessages(): Promise<ListMessagesResponse> {
  return request('GET', '/api/v1/messages');
}

export function ackMessage(envelopeHash: string): Promise<void> {
  return request('POST', '/api/v1/messages/ack', { envelope_hash: envelopeHash });
}

// User payload API
export interface UpdatePayloadResponse {
  version: number;
}

export function updatePayload(encryptedPayload: EncryptedPayload): Promise<UpdatePayloadResponse> {
  return request('PUT', '/api/v1/user/payload', { encrypted_payload: encryptedPayload });
}
