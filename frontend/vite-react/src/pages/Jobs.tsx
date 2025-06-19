import React, { useState, useEffect } from 'react'
import { Clock, CheckCircle, AlertCircle, Download, FileText, Wand2 } from 'lucide-react'
import config from '../config'

interface Job {
  id: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  mod_type: string
  original_filename: string
  original_file_size: number
  processed_url?: string
  tokens_used?: number
  credits_used?: number
  error_message?: string
  created_at: string
  updated_at: string
}

const Jobs: React.FC = () => {
  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [limit] = useState(10)

  useEffect(() => {
    fetchJobs()
  }, [page])

  const fetchJobs = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${config.apiUrl}/api/v1/mods/jobs?page=${page}&limit=${limit}`)
      const data = await response.json()
      setJobs(data.jobs || [])
    } catch (error) {
      console.error('Failed to fetch jobs:', error)
    } finally {
      setLoading(false)
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending': return <Clock className="w-5 h-5 text-yellow-500" />
      case 'processing': return <Wand2 className="w-5 h-5 text-blue-500 animate-spin" />
      case 'completed': return <CheckCircle className="w-5 h-5 text-green-500" />
      case 'failed': return <AlertCircle className="w-5 h-5 text-red-500" />
      default: return <FileText className="w-5 h-5 text-gray-500" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'bg-yellow-100 text-yellow-800'
      case 'processing': return 'bg-blue-100 text-blue-800'
      case 'completed': return 'bg-green-100 text-green-800'
      case 'failed': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const downloadProcessedFile = async (job: Job) => {
    try {
      const response = await fetch(`${config.apiUrl}/api/v1/mods/jobs/${job.id}/download`)
      const result = await response.json()
      
      if (response.ok && result.download_url) {
        window.open(result.download_url, '_blank')
      } else {
        alert('Download failed: ' + (result.error || 'Unknown error'))
      }
    } catch (error) {
      console.error('Download failed:', error)
      alert('Download failed')
    }
  }

  const formatFileSize = (bytes: number) => {
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  if (loading) {
    return (
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold mb-2">Your Processing Jobs</h1>
        <p className="text-gray-600">Track the status of your mod enhancement jobs</p>
      </div>

      {jobs.length === 0 ? (
        <div className="card text-center py-12">
          <FileText className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No jobs yet</h3>
          <p className="text-gray-500 mb-4">Upload and process your first mod to see jobs here</p>
          <a href="/upload" className="btn-primary">
            Upload Your First Mod
          </a>
        </div>
      ) : (
        <div className="space-y-4">
          {jobs.map((job) => (
            <div key={job.id} className="card">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  {getStatusIcon(job.status)}
                  <div>
                    <h3 className="font-medium text-gray-900">{job.original_filename}</h3>
                    <div className="flex items-center gap-4 text-sm text-gray-500">
                      <span>{formatFileSize(job.original_file_size)}</span>
                      <span>•</span>
                      <span className="capitalize">{job.mod_type}</span>
                      <span>•</span>
                      <span>Created {formatDate(job.created_at)}</span>
                    </div>
                  </div>
                </div>
                
                <div className="flex items-center gap-3">
                  <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(job.status)}`}>
                    {job.status.charAt(0).toUpperCase() + job.status.slice(1)}
                  </span>
                  
                  {job.status === 'completed' && job.processed_url && (
                    <button
                      onClick={() => downloadProcessedFile(job)}
                      className="btn-secondary flex items-center gap-2"
                    >
                      <Download className="w-4 h-4" />
                      Download
                    </button>
                  )}
                </div>
              </div>
              
              {job.error_message && (
                <div className="mt-3 p-3 bg-red-50 border border-red-200 rounded text-red-700 text-sm">
                  <strong>Error:</strong> {job.error_message}
                </div>
              )}
              
              {job.status === 'completed' && (
                <div className="mt-3 p-3 bg-green-50 border border-green-200 rounded text-green-700 text-sm">
                  <div className="flex items-center justify-between">
                    <span>✨ Enhancement completed successfully!</span>
                    <div className="text-xs">
                      {job.tokens_used && <span>Tokens: {job.tokens_used}</span>}
                      {job.credits_used && <span className="ml-2">Credits: {job.credits_used}</span>}
                    </div>
                  </div>
                </div>
              )}
              
              {job.status === 'processing' && (
                <div className="mt-3 p-3 bg-blue-50 border border-blue-200 rounded text-blue-700 text-sm">
                  🤖 AI is enhancing your mod... This may take a few minutes.
                </div>
              )}
            </div>
          ))}
        </div>
      )}
      
      {/* Pagination */}
      {jobs.length >= limit && (
        <div className="flex justify-center gap-2">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            className="btn-secondary disabled:opacity-50"
          >
            Previous
          </button>
          <span className="px-4 py-2 text-sm text-gray-600">
            Page {page}
          </span>
          <button
            onClick={() => setPage(p => p + 1)}
            disabled={jobs.length < limit}
            className="btn-secondary disabled:opacity-50"
          >
            Next
          </button>
        </div>
      )}
    </div>
  )
}

export default Jobs
