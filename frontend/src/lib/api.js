import axios from 'axios';

// Get API base URL from environment variable or default to localhost
let API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000';

// Remove trailing slash to prevent double slashes
API_BASE_URL = API_BASE_URL.replace(/\/$/, '');

// Force HTTPS in production to prevent mixed content errors
if (window.location.protocol === 'https:') {
  // Always use HTTPS when site is accessed over HTTPS
  if (API_BASE_URL.startsWith('http://')) {
    API_BASE_URL = API_BASE_URL.replace('http://', 'https://');
  }
}

// Create axios instance with base configuration
const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  timeout: 50000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Handle 401 errors with automatic token refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      const refreshToken = localStorage.getItem('refresh_token');
      const token = localStorage.getItem('token');
      const isLoginPage = window.location.pathname === '/login';
      const isAuthMeEndpoint = error.config?.url?.includes('/auth/me');
      const isCartEndpoint = error.config?.url?.includes('/cart');
      
      // Try to refresh token if we have a refresh token
      if (refreshToken && token && !isLoginPage) {
        originalRequest._retry = true;
        
        try {
          const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {
            refresh_token: refreshToken
          });
          
          if (response.data.access_token) {
            // Update tokens
            localStorage.setItem('token', response.data.access_token);
            if (response.data.refresh_token) {
              localStorage.setItem('refresh_token', response.data.refresh_token);
            }
            
            // Retry original request with new token
            originalRequest.headers.Authorization = `Bearer ${response.data.access_token}`;
            return api(originalRequest);
          }
        } catch (refreshError) {
          console.error('Token refresh failed:', refreshError);
          // Refresh failed, clear tokens and redirect
          localStorage.removeItem('token');
          localStorage.removeItem('refresh_token');
          localStorage.removeItem('user');
          localStorage.removeItem('permissions');
          if (!isLoginPage) {
            window.location.href = '/login';
          }
          return Promise.reject(refreshError);
        }
      }
      
      // No refresh token or refresh failed
      // Don't redirect for /auth/me (handled by ProtectedRoute)
      // Don't clear token for cart operations (might be guest cart transfer)
      if (token && !isLoginPage && !isAuthMeEndpoint && !isCartEndpoint) {
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        localStorage.removeItem('permissions');
        window.location.href = '/login';
      }
    }
    
    return Promise.reject(error);
  }
);

export default api;
export { API_BASE_URL };
