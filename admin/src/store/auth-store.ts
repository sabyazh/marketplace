import { create } from 'zustand';
import { apiPost } from '@/lib/api-client';
import type { User } from '@/types';

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  hydrate: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: null,
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiPost<{
        access_token: string;
        refresh_token: string;
        user: {
          id: string;
          full_name: string;
          email: string;
          phone: string;
          role: string;
          avatar_url: string;
        };
      }>('/auth/login/email', {
        email,
        password,
      });
      const { access_token, user: authUser } = response.data;

      if (authUser.role !== 'admin') {
        set({ isLoading: false, error: 'Access denied. Admin privileges required.' });
        return;
      }

      // Map backend auth response to admin User type
      const user: User = {
        id: authUser.id,
        email: authUser.email,
        phone: authUser.phone || '',
        first_name: authUser.full_name.split(' ')[0] || '',
        last_name: authUser.full_name.split(' ').slice(1).join(' ') || '',
        role: authUser.role as User['role'],
        status: 'active',
        avatar_url: authUser.avatar_url || undefined,
        email_verified: true,
        phone_verified: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      localStorage.setItem('admin_token', access_token);
      localStorage.setItem('admin_user', JSON.stringify(user));

      set({
        token: access_token,
        user,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (err: unknown) {
      const message =
        err instanceof Error
          ? err.message
          : (err as { response?: { data?: { message?: string } } })?.response?.data?.message ||
            'Login failed. Please check your credentials.';
      set({ isLoading: false, error: message });
      throw new Error(message);
    }
  },

  logout: () => {
    localStorage.removeItem('admin_token');
    localStorage.removeItem('admin_user');
    set({
      token: null,
      user: null,
      isAuthenticated: false,
      error: null,
    });
    window.location.href = '/login';
  },

  hydrate: () => {
    if (typeof window === 'undefined') return;

    const token = localStorage.getItem('admin_token');
    const userStr = localStorage.getItem('admin_user');

    if (token && userStr) {
      try {
        const user = JSON.parse(userStr) as User;
        set({ token, user, isAuthenticated: true });
      } catch {
        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user');
      }
    }
  },
}));
