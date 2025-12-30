import { writable } from 'svelte/store';

// Authentication store state
function createAuthStore() {
  const { subscribe, set, update } = writable({
    user: null,
    isAuthenticated: false,
    isLoading: true,
  });

  return {
    subscribe,

    // Register a new user
    async register(username, email, password, passwordConfirmation) {
      try {
        const response = await fetch('http://localhost:3000/api/auth/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include', // Send cookies
          body: JSON.stringify({
            username,
            email,
            password,
            password_confirmation: passwordConfirmation,
          }),
        });

        const data = await response.json();

        if (!data.success) {
          throw new Error(data.error.message || 'ثبت‌نام ناموفق بود');
        }

        set({
          user: data.data.user,
          isAuthenticated: true,
          isLoading: false,
        });

        return { success: true };
      } catch (error) {
        return { success: false, error: error.message };
      }
    },

    // Login existing user
    async login(email, password) {
      try {
        const response = await fetch('http://localhost:3000/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include', // Send cookies
          body: JSON.stringify({ email, password }),
        });

        const data = await response.json();

        if (!data.success) {
          throw new Error(data.error.message || 'ورود ناموفق بود');
        }

        set({
          user: data.data.user,
          isAuthenticated: true,
          isLoading: false,
        });

        return { success: true };
      } catch (error) {
        return { success: false, error: error.message };
      }
    },

    // Check if user is already authenticated
    async checkAuth() {
      try {
        const response = await fetch('http://localhost:3000/api/auth/me', {
          credentials: 'include', // Send cookies
        });

        if (response.ok) {
          const data = await response.json();
          set({
            user: data.data.user,
            isAuthenticated: true,
            isLoading: false,
          });
        } else {
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
          });
        }
      } catch (error) {
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      }
    },

    // Logout user
    async logout() {
      try {
        await fetch('http://localhost:3000/api/auth/logout', {
          method: 'POST',
          credentials: 'include',
        });
      } catch (error) {
        console.error('Logout failed:', error);
      } finally {
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      }
    },
  };
}

export const authStore = createAuthStore();
