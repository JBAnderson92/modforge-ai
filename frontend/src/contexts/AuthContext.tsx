import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import config from '../config'

// User interface for authentication
interface User {
  id: string
  email: string
  display_name: string
  credits: number
  plan: string
}

interface AuthContextType {
  user: User | null
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, displayName: string) => Promise<void>
  logout: () => Promise<void>
  loading: boolean
  token: string | null
}

const AuthContext = createContext<AuthContextType | null>(null)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  // Real authentication with backend API
  const login = async (email: string, password: string) => {
    setLoading(true)
    try {
      const response = await fetch(`${config.apiUrl}/api/v1/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Login failed')
      }

      setUser(data.user)
      setToken(data.token)
      localStorage.setItem('modforge_user', JSON.stringify(data.user))
      localStorage.setItem('modforge_token', data.token)
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const register = async (email: string, password: string, displayName: string) => {
    setLoading(true)
    try {
      const response = await fetch(`${config.apiUrl}/api/v1/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          email, 
          password, 
          display_name: displayName 
        }),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Registration failed')
      }

      setUser(data.user)
      setToken(data.token)
      localStorage.setItem('modforge_user', JSON.stringify(data.user))
      localStorage.setItem('modforge_token', data.token)
    } catch (error) {
      console.error('Registration failed:', error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const logout = async () => {
    setUser(null)
    setToken(null)
    localStorage.removeItem('modforge_user')
    localStorage.removeItem('modforge_token')
  }

  // Check for existing session on mount
  useEffect(() => {
    const storedUser = localStorage.getItem('modforge_user')
    const storedToken = localStorage.getItem('modforge_token')
    
    if (storedUser && storedToken) {
      try {
        setUser(JSON.parse(storedUser))
        setToken(storedToken)
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        localStorage.removeItem('modforge_user')
        localStorage.removeItem('modforge_token')
      }
    }
    setLoading(false)
  }, [])

  const value: AuthContextType = {
    user,
    login,
    register,
    logout,
    loading,
    token
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}
