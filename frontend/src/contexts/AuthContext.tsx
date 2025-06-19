import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'

// Mock user interface for MVP
interface User {
  id: string
  email: string
  displayName: string
  credits: number
  plan: string
}

interface AuthContextType {
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  loading: boolean
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
  const [loading, setLoading] = useState(true)

  // Mock authentication - in production, integrate with Firebase Auth
  const login = async (email: string, _password: string) => {
    setLoading(true)
    try {
      // Mock login - replace with Firebase Auth
      const mockUser: User = {
        id: '1',
        email: email,
        displayName: email.split('@')[0],
        credits: 10,
        plan: 'free'
      }
      setUser(mockUser)
      localStorage.setItem('modforge_user', JSON.stringify(mockUser))
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const logout = async () => {
    setUser(null)
    localStorage.removeItem('modforge_user')
  }

  // Check for existing session on mount
  useEffect(() => {
    const storedUser = localStorage.getItem('modforge_user')
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser))
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        localStorage.removeItem('modforge_user')
      }
    }
    setLoading(false)
  }, [])

  const value: AuthContextType = {
    user,
    login,
    logout,
    loading
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}
