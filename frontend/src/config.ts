// Environment configuration for frontend
const config = {
  apiUrl: (import.meta as any).env.VITE_API_URL || 'http://localhost:8080',
  environment: (import.meta as any).env.MODE || 'development',
  version: '1.0.0'
}

export default config
