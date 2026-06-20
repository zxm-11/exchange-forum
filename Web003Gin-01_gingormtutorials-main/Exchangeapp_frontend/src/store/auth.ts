import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from '../axios';

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const username = ref<string | null>(localStorage.getItem('username'));
  const unreadCount = ref(0);

  const isAuthenticated = computed(() => !!token.value);

  const parseUsername = (tokenStr: string): string => {
    try {
      const cleanToken = tokenStr.startsWith('Bearer ') ? tokenStr.slice(7) : tokenStr;
      const payload = JSON.parse(atob(cleanToken.split('.')[1]));
      return payload.username || '';
    } catch {
      return '';
    }
  };

  const login = async (user: string, password: string) => {
    try {
      const response = await axios.post('/auth/login', { username: user, password });
      token.value = response.data.token;
      localStorage.setItem('token', token.value || '');
      const u = response.data.username || parseUsername(response.data.token);
      username.value = u;
      localStorage.setItem('username', u);
    } catch (error) {
      throw new Error(`Login failed! ${error}`);
    }
  };

  const register = async (user: string, password: string) => {
    try {
      const response = await axios.post('/auth/register', { username: user, password });
      token.value = response.data.token;
      localStorage.setItem('token', token.value || '');
      const u = response.data.username || parseUsername(response.data.token);
      username.value = u;
      localStorage.setItem('username', u);
    } catch (error) {
      throw new Error(`Register failed! ${error}`);
    }
  };

  const logout = () => {
    token.value = null;
    username.value = null;
    unreadCount.value = 0;
    localStorage.removeItem('token');
    localStorage.removeItem('username');
  };

  const fetchUnreadCount = async () => {
    if (!token.value) return;
    try {
      const res = await axios.get('/notifications/unread_count');
      unreadCount.value = res.data.Unread_count || 0;
    } catch {
      // 静默失败
    }
  };

  return {
    token,
    username,
    unreadCount,
    isAuthenticated,
    login,
    register,
    logout,
    fetchUnreadCount
  };
});
