import axios from 'axios'

// Initialize Axios instance with backend base URL
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add request interceptor to include JWT token if available
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Fetch all products
export const fetchProducts = () => {
  return apiClient.get('/products')
}

// Fetch a single product by ID
export const fetchProduct = (id) => {
  return apiClient.get(`/products/${id}`)
}

export default apiClient
