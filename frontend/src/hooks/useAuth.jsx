import { createContext, useContext, useState, useEffect } from 'react'
import apiClient from '../services/api'

/**
 * AuthContext manages user authentication state and provides login/register/logout methods.
 */
const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // Initialize auth state from localStorage on mount
  useEffect(() => {
    const savedToken = localStorage.getItem('authToken')
    const savedUser = localStorage.getItem('authUser')

    if (savedToken && savedUser) {
      setToken(savedToken)
      setUser(JSON.parse(savedUser))
    }
    setLoading(false)
  }, [])

  // Register user
  const register = async (email, password) => {
    try {
      setLoading(true)
      setError(null)

      const response = await apiClient.post('/auth/register', {
        email,
        password,
      })

      const { user: userData, token: accessToken } = response.data

      // Store in state and localStorage
      setUser(userData)
      setToken(accessToken)
      localStorage.setItem('authToken', accessToken)
      localStorage.setItem('authUser', JSON.stringify(userData))

      return { success: true }
    } catch (err) {
      const errorMsg =
        err.response?.data?.error || 'Failed to register. Please try again.'
      setError(errorMsg)
      return { success: false, error: errorMsg }
    } finally {
      setLoading(false)
    }
  }

  // Login user
  const login = async (email, password) => {
    try {
      setLoading(true)
      setError(null)

      const response = await apiClient.post('/auth/login', {
        email,
        password,
      })

      const { user: userData, token: accessToken } = response.data

      // Store in state and localStorage
      setUser(userData)
      setToken(accessToken)
      localStorage.setItem('authToken', accessToken)
      localStorage.setItem('authUser', JSON.stringify(userData))

      return { success: true }
    } catch (err) {
      const errorMsg =
        err.response?.data?.error || 'Failed to login. Please try again.'
      setError(errorMsg)
      return { success: false, error: errorMsg }
    } finally {
      setLoading(false)
    }
  }

  // Logout user
  const logout = () => {
    setUser(null)
    setToken(null)
    setError(null)
    localStorage.removeItem('authToken')
    localStorage.removeItem('authUser')
  }

  const value = {
    user,
    token,
    loading,
    error,
    register,
    login,
    logout,
    isAuthenticated: !!user && !!token,
  }

  return (
    <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
  )
}

/**
 * Hook to use the auth context
 */
export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
