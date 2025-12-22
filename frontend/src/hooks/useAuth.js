import { useState } from 'react';
import api from '../lib/api';

export const useAuth = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const login = async (email, password) => {
    try {
      setLoading(true);
      setError('');
      const response = await api.post('/auth/login', { email_or_username: email, password });
      
      // Backend returns: { access_token, refresh_token, user }
      if (response.data.access_token) {
        localStorage.setItem('token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
        
        // Transfer guest cart after login
        await transferGuestCartAfterLogin();
        
        return response.data;
      }
      return null;
    } catch (err) {
      setError(err.response?.data?.error || 'Login failed. Please try again.');
      return null;
    } finally {
      setLoading(false);
    }
  };

  const transferGuestCartAfterLogin = async () => {
    try {
      const guestCart = localStorage.getItem('readagain_guest_cart');
      if (!guestCart) return;
      
      const guestItems = JSON.parse(guestCart);
      if (guestItems.length === 0) return;
      
      // Backend expects: { book_ids: [1, 2, 3] }
      const bookIds = guestItems.map(item => item.book_id || item.id);
      await api.post('/cart/merge', { book_ids: bookIds });
      
      // Clear guest cart
      localStorage.removeItem('readagain_guest_cart');
    } catch (err) {
      console.error('Error transferring guest cart:', err);
    }
  };

  const signup = async (userData) => {
    try {
      setLoading(true);
      setError('');
      const response = await api.post('/auth/register', userData);
      
      // Backend returns: { message, user }
      if (response.data.user) {
        return true;
      }
      return false;
    } catch (err) {
      setError(err.response?.data?.error || 'Signup failed. Please try again.');
      return false;
    } finally {
      setLoading(false);
    }
  };

  const forgotPassword = async (email) => {
    try {
      setLoading(true);
      setError('');
      await api.post('/auth/forgot-password', { email });
      return true;
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to send reset email.');
      return false;
    } finally {
      setLoading(false);
    }
  };

  const resetPassword = async (token, password) => {
    try {
      setLoading(true);
      setError('');
      await api.post('/auth/reset-password', { token, password });
      return true;
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to reset password.');
      return false;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('permissions');
    window.location.href = '/';
  };

  const isAuthenticated = () => {
    return !!localStorage.getItem('token');
  };

  const getUser = () => {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  };

  const getPermissions = () => {
    const permissions = localStorage.getItem('permissions');
    return permissions ? JSON.parse(permissions) : [];
  };

  return {
    login,
    signup,
    forgotPassword,
    resetPassword,
    logout,
    isAuthenticated,
    getUser,
    getPermissions,
    loading,
    error
  };
};
