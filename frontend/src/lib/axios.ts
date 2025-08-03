import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Create axios instance
const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
axiosInstance.interceptors.request.use(
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
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;

// API endpoints
export const API_ENDPOINTS = {
  // Auth
  login: '/auth/login',
  logout: '/auth/logout',
  refresh: '/auth/refresh',
  
  // Users
  users: '/users',
  profile: '/profile',
  
  // Inventory
  inventory: '/inventory',
  materials: '/materials',
  warehouses: '/warehouses',
  stockMovements: '/stock-movements',
  
  // Orders
  orders: '/orders',
  quotes: '/quotes',
  inquiries: '/inquiries',
  
  // Customers
  customers: '/customers',
  suppliers: '/suppliers',
  
  // Products
  products: '/products',
  categories: '/categories',
  
  // Reports
  reports: '/reports',
  analytics: '/analytics',
};