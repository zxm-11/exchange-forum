import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from '../axios';

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const username = ref<string | null>(localStorage.getItem('username'));

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
    localStorage.removeItem('token');
    localStorage.removeItem('username');
  };

  return {
    token,
    username,
    isAuthenticated,
    login,
    register,
    logout
  };
});
