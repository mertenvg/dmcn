import { useState, useCallback } from 'react';

export interface Contact {
  address: string;
  name: string;
  fingerprint: string;
}

export function useContacts() {
  const [contacts, setContacts] = useState<Contact[]>(() => {
    try {
      const saved = localStorage.getItem('dmcn_contacts');
      return saved ? JSON.parse(saved) : [];
    } catch {
      return [];
    }
  });

  const addContact = useCallback((contact: Contact) => {
    setContacts(prev => {
      const updated = [...prev.filter(c => c.address !== contact.address), contact];
      localStorage.setItem('dmcn_contacts', JSON.stringify(updated));
      return updated;
    });
  }, []);

  const removeContact = useCallback((address: string) => {
    setContacts(prev => {
      const updated = prev.filter(c => c.address !== address);
      localStorage.setItem('dmcn_contacts', JSON.stringify(updated));
      return updated;
    });
  }, []);

  return { contacts, addContact, removeContact };
}
