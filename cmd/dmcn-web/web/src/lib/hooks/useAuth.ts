import { createContext, useContext, useState, useCallback, useEffect, ReactNode, createElement } from 'react';
import { setSessionToken, getSessionToken } from '../api/client';

interface AuthState {
  isAuthenticated: boolean;
  address: string | null;
  sessionToken: string | null;
}

interface AuthContextValue extends AuthState {
  setSession: (address: string, token: string) => void;
  clearSession: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(() => {
    const saved = localStorage.getItem('dmcn_session');
    if (saved) {
      try {
        const { address, sessionToken } = JSON.parse(saved);
        setSessionToken(sessionToken);
        return { isAuthenticated: true, address, sessionToken };
      } catch {
        // ignore
      }
    }
    return { isAuthenticated: false, address: null, sessionToken: null };
  });

  const setSession = useCallback((address: string, token: string) => {
    setSessionToken(token);
    localStorage.setItem('dmcn_session', JSON.stringify({ address, sessionToken: token }));
    setState({ isAuthenticated: true, address, sessionToken: token });
  }, []);

  const clearSession = useCallback(() => {
    setSessionToken(null);
    localStorage.removeItem('dmcn_session');
    localStorage.removeItem('dmcn_encrypted_payload');
    setState({ isAuthenticated: false, address: null, sessionToken: null });
  }, []);

  return createElement(AuthContext.Provider, { value: { ...state, setSession, clearSession } }, children);
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
