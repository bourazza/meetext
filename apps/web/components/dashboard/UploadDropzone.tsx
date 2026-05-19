'use client'

import React, { useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { UploadCloud, File, Loader2, CheckCircle2, X } from 'lucide-react'

type UploadState = 'idle' | 'uploading' | 'processing' | 'complete'

export function UploadDropzone() {
  const [state, setState] = useState<UploadState>('idle')
  const [progress, setProgress] = useState(0)

  const simulateUpload = useCallback(() => {
    setState('uploading')
    setProgress(0)
    
    // Simulate upload progress
    const interval = setInterval(() => {
      setProgress((p) => {
        if (p >= 100) {
          clearInterval(interval)
          setState('processing')
          // Simulate processing time
          setTimeout(() => {
            setState('complete')
            setTimeout(() => setState('idle'), 3000) // Reset after 3s
          }, 2000)
          return 100
        }
        return p + 5
      })
    }, 100)
  }, [])

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    if (state === 'idle') {
      simulateUpload()
    }
  }

  return (
    <div className="w-full max-w-3xl rounded-2xl border border-zinc-200 bg-white p-8 shadow-sm">
      <AnimatePresence mode="wait">
        {state === 'idle' && (
          <motion.div
            key="idle"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onDragOver={handleDragOver}
            onDrop={handleDrop}
            className="group relative flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed border-zinc-200 bg-zinc-50/50 py-16 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
            onClick={simulateUpload}
          >
            <div className="mb-4 rounded-full bg-zinc-100 p-4 text-zinc-500 transition-transform group-hover:scale-105">
              <UploadCloud className="h-8 w-8" />
            </div>
            <h3 className="mb-1 text-lg font-semibold text-zinc-900">Upload a meeting recording</h3>
            <p className="text-sm text-zinc-500">Drag and drop audio, video, or PDF files here, or click to browse.</p>
          </motion.div>
        )}

        {state === 'uploading' && (
          <motion.div
            key="uploading"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="flex flex-col items-center py-12"
          >
            <File className="mb-6 h-12 w-12 text-zinc-400" />
            <div className="mb-2 text-sm font-medium text-zinc-900">Uploading meeting_recording.mp4...</div>
            <div className="mb-2 w-full max-w-md overflow-hidden rounded-full bg-zinc-100">
              <motion.div
                className="h-2 bg-zinc-900"
                initial={{ width: 0 }}
                animate={{ width: `${progress}%` }}
                transition={{ ease: 'linear' }}
              />
            </div>
            <div className="text-xs text-zinc-500">{progress}%</div>
          </motion.div>
        )}

        {state === 'processing' && (
          <motion.div
            key="processing"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="flex flex-col items-center py-16 text-center"
          >
            <Loader2 className="mb-4 h-8 w-8 animate-spin text-zinc-900" />
            <h3 className="mb-1 text-lg font-semibold text-zinc-900">Generating AI intelligence...</h3>
            <p className="text-sm text-zinc-500">Transcribing audio and extracting tasks, decisions, and goals.</p>
          </motion.div>
        )}

        {state === 'complete' && (
          <motion.div
            key="complete"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="flex flex-col items-center py-16 text-center"
          >
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: 'spring', stiffness: 200, damping: 20 }}
            >
              <CheckCircle2 className="mb-4 h-12 w-12 text-emerald-500" />
            </motion.div>
            <h3 className="mb-1 text-lg font-semibold text-zinc-900">Processing complete</h3>
            <p className="text-sm text-zinc-500">Your meeting documentation is ready.</p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
