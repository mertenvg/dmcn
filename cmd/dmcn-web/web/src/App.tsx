import { HashRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './lib/hooks/useAuth';
import { KeysProvider } from './lib/hooks/useKeys';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Inbox } from './pages/Inbox';
import { Compose } from './pages/Compose';
import { MessageView } from './pages/MessageView';
import { Contacts } from './pages/Contacts';
import { Settings } from './pages/Settings';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

export function App() {
  return (
    <AuthProvider>
      <KeysProvider>
        <HashRouter>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/inbox" element={<ProtectedRoute><Inbox /></ProtectedRoute>} />
            <Route path="/compose" element={<ProtectedRoute><Compose /></ProtectedRoute>} />
            <Route path="/message/:id" element={<ProtectedRoute><MessageView /></ProtectedRoute>} />
            <Route path="/contacts" element={<ProtectedRoute><Contacts /></ProtectedRoute>} />
            <Route path="/settings" element={<ProtectedRoute><Settings /></ProtectedRoute>} />
            <Route path="*" element={<Navigate to="/inbox" replace />} />
          </Routes>
        </HashRouter>
      </KeysProvider>
    </AuthProvider>
  );
}
