import { createContext, useContext, useState, useCallback, useEffect, ReactNode, createElement } from 'react';
import type { IdentityKeyPair } from '../crypto/keys';
import { useAuth } from './useAuth';

interface KeysContextValue {
  keys: IdentityKeyPair | null;
  setKeys: (keys: IdentityKeyPair) => void;
  clearKeys: () => void;
}

const KeysContext = createContext<KeysContextValue | null>(null);

export function KeysProvider({ children }: { children: ReactNode }) {
  const [keys, setKeysState] = useState<IdentityKeyPair | null>(null);
  const { isAuthenticated } = useAuth();

  const setKeys = useCallback((k: IdentityKeyPair) => {
    setKeysState(k);
  }, []);

  const clearKeys = useCallback(() => {
    setKeysState(null);
  }, []);

  // Clear keys when logged out or tab closes
  useEffect(() => {
    if (!isAuthenticated) clearKeys();
  }, [isAuthenticated, clearKeys]);

  useEffect(() => {
    const handleUnload = () => clearKeys();
    window.addEventListener('beforeunload', handleUnload);
    return () => window.removeEventListener('beforeunload', handleUnload);
  }, [clearKeys]);

  return createElement(KeysContext.Provider, { value: { keys, setKeys, clearKeys } }, children);
}

export function useKeys(): KeysContextValue {
  const ctx = useContext(KeysContext);
  if (!ctx) throw new Error('useKeys must be used within KeysProvider');
  return ctx;
}
