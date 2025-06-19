import React, { useState, useCallback } from 'react'
import { useDropzone } from 'react-dropzone'
import { Upload as UploadIcon, FileText, Settings, Wand2, Download, CheckCircle, AlertCircle, Clock } from 'lucide-react'
import config from '../config'

interface UploadedFile {
  file: File
  jobId?: string
  status?: 'pending' | 'processing' | 'completed' | 'failed'
  modType?: string
  errorMessage?: string
  processedUrl?: string
}

interface Preset {
  id: string
  name: string
  description: string
  credit_cost: number
}

const Upload: React.FC = () => {
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([])
  const [presets, setPresets] = useState<Preset[]>([])
  const [selectedPreset, setSelectedPreset] = useState<string>('')
  const [customPrompt, setCustomPrompt] = useState<string>('')
  const [isProcessing, setIsProcessing] = useState(false)

  // Fetch presets when component mounts
  React.useEffect(() => {
    fetchPresets()
  }, [])

  const fetchPresets = async () => {
    try {
      const response = await fetch(`${config.apiUrl}/api/v1/presets/`)
      const data = await response.json()
      setPresets(data.presets || [])
    } catch (error) {
      console.error('Failed to fetch presets:', error)
    }
  }

  const onDrop = useCallback(async (acceptedFiles: File[]) => {
    
    for (const file of acceptedFiles) {
      // Add file to state immediately
      setUploadedFiles(prev => [...prev, { file, status: 'pending' }])
      
      try {
        // Upload file
        const formData = new FormData()
        formData.append('mod_file', file)
        
        const response = await fetch(`${config.apiUrl}/api/v1/mods/upload`, {
          method: 'POST',
          body: formData,
        })
        
        const result = await response.json()
        
        if (response.ok) {
          // Update file with job info
          setUploadedFiles(prev => 
            prev.map(f => 
              f.file === file 
                ? { ...f, jobId: result.job_id, status: result.status, modType: result.mod_type }
                : f
            )
          )
        } else {
          setUploadedFiles(prev => 
            prev.map(f => 
              f.file === file 
                ? { ...f, status: 'failed', errorMessage: result.error }
                : f
            )
          )
        }
      } catch (error) {
        console.error('Upload failed:', error)
        setUploadedFiles(prev => 
          prev.map(f => 
            f.file === file 
              ? { ...f, status: 'failed', errorMessage: 'Upload failed' }
              : f
          )
        )
      }
    }
  }, [])

  const processWithAI = async (fileIndex: number) => {
    const file = uploadedFiles[fileIndex]
    if (!file.jobId) return

    setIsProcessing(true)
    
    try {
      const response = await fetch(`${config.apiUrl}/api/v1/mods/jobs/${file.jobId}/process`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          preset_id: selectedPreset || 'minecraft_balance',
          prompt: customPrompt || 'Enhance this mod with balanced improvements',
          model_config: 'default'
        }),
      })
      
      const result = await response.json()
      
      if (response.ok) {
        // Update status to processing
        setUploadedFiles(prev => 
          prev.map((f, i) => 
            i === fileIndex 
              ? { ...f, status: 'processing' }
              : f
          )
        )
        
        // Poll for completion (in a real app, use WebSocket)
        pollJobStatus(file.jobId!, fileIndex)
      } else {
        setUploadedFiles(prev => 
          prev.map((f, i) => 
            i === fileIndex 
              ? { ...f, status: 'failed', errorMessage: result.error }
              : f
          )
        )
      }
    } catch (error) {
      console.error('Processing failed:', error)
    }
    
    setIsProcessing(false)
  }

  const pollJobStatus = async (jobId: string, fileIndex: number) => {
    const poll = async () => {
      try {
        const response = await fetch(`${config.apiUrl}/api/v1/mods/jobs/${jobId}`)
        const job = await response.json()
        
        setUploadedFiles(prev => 
          prev.map((f, i) => 
            i === fileIndex 
              ? { 
                  ...f, 
                  status: job.status,
                  processedUrl: job.processed_url,
                  errorMessage: job.error_message 
                }
              : f
          )
        )
        
        if (job.status === 'processing') {
          setTimeout(poll, 3000) // Poll every 3 seconds
        }
      } catch (error) {
        console.error('Failed to poll job status:', error)
      }
    }
    
    poll()
  }

  const downloadProcessedFile = async (fileIndex: number) => {
    const file = uploadedFiles[fileIndex]
    if (!file.jobId) return

    try {
      const response = await fetch(`${config.apiUrl}/api/v1/mods/jobs/${file.jobId}/download`)
      const result = await response.json()
      
      if (response.ok && result.download_url) {
        // Open download URL
        window.open(result.download_url, '_blank')
      } else {
        alert('Download failed: ' + (result.error || 'Unknown error'))
      }
    } catch (error) {
      console.error('Download failed:', error)
      alert('Download failed')
    }
  }

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'application/java-archive': ['.jar'],
      'application/zip': ['.zip'],
      'application/json': ['.json'],
      'text/plain': ['.mcmeta']
    },
    maxSize: 100 * 1024 * 1024, // 100MB
  })

  const getStatusIcon = (status?: string) => {
    switch (status) {
      case 'pending': return <Clock className="w-5 h-5 text-yellow-500" />
      case 'processing': return <Settings className="w-5 h-5 text-blue-500 animate-spin" />
      case 'completed': return <CheckCircle className="w-5 h-5 text-green-500" />
      case 'failed': return <AlertCircle className="w-5 h-5 text-red-500" />
      default: return <FileText className="w-5 h-5 text-gray-500" />
    }
  }

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      <div>
        <h1 className="text-3xl font-bold mb-2">Upload & Enhance Your Mods</h1>
        <p className="text-gray-600">Upload .jar, .zip, .json, or .mcmeta files to enhance them with AI</p>
      </div>

      {/* Upload Area */}
      <div className="card">
        <div 
          {...getRootProps()} 
          className={`border-2 border-dashed rounded-lg p-12 text-center cursor-pointer transition-colors ${
            isDragActive 
              ? 'border-primary-500 bg-primary-50' 
              : 'border-gray-300 hover:border-gray-400'
          }`}
        >
          <input {...getInputProps()} />
          <UploadIcon className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          {isDragActive ? (
            <p className="text-primary-600 font-medium">Drop the files here...</p>
          ) : (
            <>
              <p className="text-gray-600 font-medium mb-2">
                Drag & drop your mod files here, or click to browse
              </p>
              <p className="text-sm text-gray-500">
                Supports .jar, .zip, .json, .mcmeta files up to 100MB
              </p>
            </>
          )}
        </div>
      </div>

      {/* AI Enhancement Settings */}
      {uploadedFiles.length > 0 && (
        <div className="card">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <Wand2 className="w-5 h-5" />
            AI Enhancement Settings
          </h2>
          
          <div className="grid md:grid-cols-2 gap-6">
            {/* Preset Selection */}
            <div>
              <label className="block text-sm font-medium mb-2">Choose Enhancement Preset</label>
              <select 
                value={selectedPreset} 
                onChange={(e) => setSelectedPreset(e.target.value)}
                className="input"
              >
                <option value="">Select a preset...</option>
                {presets.map(preset => (
                  <option key={preset.id} value={preset.id}>
                    {preset.name} ({preset.credit_cost} credits) - {preset.description}
                  </option>
                ))}
              </select>
            </div>

            {/* Custom Prompt */}
            <div>
              <label className="block text-sm font-medium mb-2">Custom Instructions (Optional)</label>
              <textarea
                value={customPrompt}
                onChange={(e) => setCustomPrompt(e.target.value)}
                placeholder="Describe specific changes you want..."
                className="input min-h-[100px] resize-y"
              />
            </div>
          </div>
        </div>
      )}

      {/* Uploaded Files */}
      {uploadedFiles.length > 0 && (
        <div className="card">
          <h2 className="text-xl font-semibold mb-4">Uploaded Files</h2>
          <div className="space-y-4">
            {uploadedFiles.map((uploadedFile, index) => (
              <div key={index} className="border border-gray-200 rounded-lg p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {getStatusIcon(uploadedFile.status)}
                    <div>
                      <p className="font-medium">{uploadedFile.file.name}</p>
                      <p className="text-sm text-gray-500">
                        {(uploadedFile.file.size / (1024 * 1024)).toFixed(1)} MB
                        {uploadedFile.modType && ` • ${uploadedFile.modType}`}
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    {uploadedFile.status === 'pending' && uploadedFile.jobId && (
                      <button
                        onClick={() => processWithAI(index)}
                        disabled={isProcessing}
                        className="btn-primary flex items-center gap-2"
                      >
                        <Wand2 className="w-4 h-4" />
                        Enhance with AI
                      </button>
                    )}
                    
                    {uploadedFile.status === 'completed' && uploadedFile.processedUrl && (
                      <button
                        onClick={() => downloadProcessedFile(index)}
                        className="btn-secondary flex items-center gap-2"
                      >
                        <Download className="w-4 h-4" />
                        Download Enhanced
                      </button>
                    )}
                  </div>
                </div>
                
                {uploadedFile.errorMessage && (
                  <div className="mt-2 p-2 bg-red-50 border border-red-200 rounded text-red-700 text-sm">
                    {uploadedFile.errorMessage}
                  </div>
                )}
                
                {uploadedFile.status === 'processing' && (
                  <div className="mt-2 p-2 bg-blue-50 border border-blue-200 rounded text-blue-700 text-sm">
                    AI is enhancing your mod... This may take a few minutes.
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default Upload
