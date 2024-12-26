import { writable } from 'svelte/store';
import { config } from '../config';

interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
}

function createAuthStore() {
  const { subscribe, set } = writable<AuthState>({
    isAuthenticated: false,
    token: null
  });

  return {
    subscribe,
     login: async (username: string, password: string) => {
    try {
      const response = await fetch(`${config.apiUrl}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password }),
        credentials: 'include'  // Important for cookie handling
      });

      if (response.ok) {
        const data = await response.json();
        // Store token in localStorage
        localStorage.setItem('token', data.token);
        
        // Set cookie with proper attributes
        document.cookie = `token=${data.token}; path=/; SameSite=Strict; Secure`;
        
        set({ isAuthenticated: true, token: data.token });
        return true;
      }
      return false;
    } catch (error) {
      console.error('Login error:', error);
        return false;
      }
    },
    logout: () => {
      localStorage.removeItem('token');
      document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
      set({ isAuthenticated: false, token: null });
    },
    initialize: () => {
      const token = localStorage.getItem('token');
      if (token) {
        set({ isAuthenticated: true, token });
      }
    }
  };
}

export const auth = createAuthStore();
