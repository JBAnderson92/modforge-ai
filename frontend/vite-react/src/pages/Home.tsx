import React from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const Home: React.FC = () => {
  const { user } = useAuth()

  return (
    <div className="max-w-6xl mx-auto">
      {/* Hero Section */}
      <div className="text-center py-16">
        <h1 className="text-5xl font-bold text-gray-900 mb-6">
          Transform Your Game Mods with <span className="text-primary-600">AI</span>
        </h1>
        <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
          Upload your Minecraft, Skyrim, or Lua mods and let AI rebalance, translate, 
          enhance, and transform them into something extraordinary.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          {user ? (
            <Link to="/upload" className="btn-primary text-lg px-8 py-3">
              Start Modding with AI
            </Link>
          ) : (
            <button className="btn-primary text-lg px-8 py-3">
              Get Started Free
            </button>
          )}
          <button className="btn-secondary text-lg px-8 py-3">
            Watch Demo
          </button>
        </div>
      </div>

      {/* Features Grid */}
      <div className="grid md:grid-cols-3 gap-8 py-16">
        <div className="card text-center">
          <div className="w-12 h-12 bg-primary-100 rounded-lg mx-auto mb-4 flex items-center justify-center">
            <span className="text-primary-600 text-2xl">🎮</span>
          </div>
          <h3 className="text-xl font-semibold mb-3">Multi-Game Support</h3>
          <p className="text-gray-600">
            Works with Minecraft JSON mods, Skyrim ESP files, and Lua scripts. 
            More games coming soon.
          </p>
        </div>

        <div className="card text-center">
          <div className="w-12 h-12 bg-primary-100 rounded-lg mx-auto mb-4 flex items-center justify-center">
            <span className="text-primary-600 text-2xl">🤖</span>
          </div>
          <h3 className="text-xl font-semibold mb-3">AI-Powered Transformations</h3>
          <p className="text-gray-600">
            Rebalance stats, rewrite descriptions, translate content, and generate 
            new items with GPT-4.
          </p>
        </div>

        <div className="card text-center">
          <div className="w-12 h-12 bg-primary-100 rounded-lg mx-auto mb-4 flex items-center justify-center">
            <span className="text-primary-600 text-2xl">⚡</span>
          </div>
          <h3 className="text-xl font-semibold mb-3">Fast & Secure</h3>
          <p className="text-gray-600">
            Process mods in seconds with enterprise-grade security and 
            automatic virus scanning.
          </p>
        </div>
      </div>

      {/* Preset Examples */}
      <div className="py-16">
        <h2 className="text-3xl font-bold text-center mb-12">
          Popular AI Presets
        </h2>
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[
            { name: "Make Lore-Friendly", desc: "Rewrite names and descriptions", cost: 1 },
            { name: "Balance Items", desc: "Adjust stats and recipes", cost: 2 },
            { name: "Translate Mod", desc: "Convert to any language", cost: 1 },
            { name: "Add Content", desc: "Generate new items", cost: 3 },
          ].map((preset, index) => (
            <div key={index} className="card">
              <h4 className="font-semibold mb-2">{preset.name}</h4>
              <p className="text-gray-600 text-sm mb-3">{preset.desc}</p>
              <div className="flex justify-between items-center">
                <span className="text-primary-600 font-medium">{preset.cost} credit{preset.cost > 1 ? 's' : ''}</span>
                <button className="text-primary-600 hover:text-primary-700 font-medium">
                  Try it →
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* CTA Section */}
      {!user && (
        <div className="bg-primary-50 rounded-2xl p-8 text-center">
          <h2 className="text-2xl font-bold mb-4">Ready to Transform Your Mods?</h2>
          <p className="text-gray-600 mb-6">
            Join thousands of modders using AI to enhance their creations.
          </p>
          <button className="btn-primary">
            Start Free Trial
          </button>
        </div>
      )}
    </div>
  )
}

export default Home
