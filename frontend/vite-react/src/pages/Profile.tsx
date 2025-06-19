import React from 'react'
import { useAuth } from '../contexts/AuthContext'

const Profile: React.FC = () => {
  const { user } = useAuth()

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Profile</h1>
      <div className="grid md:grid-cols-2 gap-8">
        <div className="card">
          <h2 className="text-xl font-semibold mb-4">Account Information</h2>
          {user && (
            <div className="space-y-3">
              <div>
                <label className="text-sm font-medium text-gray-500">Email</label>
                <p className="text-gray-900">{user.email}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-gray-500">Display Name</label>
                <p className="text-gray-900">{user.displayName}</p>
              </div>
              <div>
                <label className="text-sm font-medium text-gray-500">Plan</label>
                <p className="text-gray-900 capitalize">{user.plan}</p>
              </div>
            </div>
          )}
        </div>
        
        <div className="card">
          <h2 className="text-xl font-semibold mb-4">Credits & Usage</h2>
          {user && (
            <div className="space-y-3">
              <div>
                <label className="text-sm font-medium text-gray-500">Available Credits</label>
                <p className="text-2xl font-bold text-blue-600">{user.credits}</p>
              </div>
              <button className="btn-primary w-full">
                Purchase More Credits
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Profile
