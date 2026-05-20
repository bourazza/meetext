'use client'

import React, { useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { UploadCloud, File, Loader2, CheckCircle2, ShieldCheck, Music, Video, FileText, Link as LinkIcon } from 'lucide-react'

type UploadState = 'idle' | 'uploading' | 'processing' | 'complete'

export function UploadDropzone() {
  const [state, setState] = useState<UploadState>('idle')
  const [progress, setProgress] = useState(0)

  const simulateUpload = useCallback(() => {
    setState('uploading')
    setProgress(0)
    
    const interval = setInterval(() => {
      setProgress((p) => {
        if (p >= 100) {
          clearInterval(interval)
          setState('processing')
          setTimeout(() => {
            setState('complete')
            setTimeout(() => setState('idle'), 3000)
          }, 2000)
          return 100
        }
        return p + 5
      })
    }, 100)
  }, [])

  const handleDragOver = (e: React.DragEvent) => e.preventDefault()
  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    if (state === 'idle') simulateUpload()
  }

  return (
    <div className="w-full">
      <AnimatePresence mode="wait">
        {state === 'idle' && (
          <motion.div
            key="idle"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="flex flex-col items-center"
          >
            <div
              onDragOver={handleDragOver}
              onDrop={handleDrop}
              onClick={simulateUpload}
              className="group w-full cursor-pointer rounded-2xl border-2 border-dashed border-zinc-200 bg-zinc-50/50 py-24 text-center transition-colors hover:border-zinc-300 hover:bg-zinc-50"
            >
              <div className="mb-6 inline-flex items-center justify-center rounded-2xl bg-zinc-100 p-5 text-zinc-900 transition-transform group-hover:scale-105">
                <UploadCloud className="h-10 w-10" />
              </div>
              <h3 className="mb-2 text-2xl font-semibold text-zinc-950">Drop your meeting file here</h3>
              <p className="mb-10 text-base text-zinc-500">Audio, video, PDF, or paste a meeting URL</p>

              <div className="mb-10 flex flex-wrap justify-center gap-4 px-6">
                <button className="flex items-center gap-2 rounded-lg border border-zinc-200 bg-white px-5 py-2.5 text-sm font-medium text-zinc-700 shadow-sm transition hover:bg-zinc-50 hover:text-zinc-900" onClick={(e) => e.stopPropagation()}>
                  <Music className="h-4 w-4" /> Audio
                </button>
                <button className="flex items-center gap-2 rounded-lg border border-zinc-200 bg-white px-5 py-2.5 text-sm font-medium text-zinc-700 shadow-sm transition hover:bg-zinc-50 hover:text-zinc-900" onClick={(e) => e.stopPropagation()}>
                  <Video className="h-4 w-4" /> Video
                </button>
                <button className="flex items-center gap-2 rounded-lg border border-zinc-200 bg-white px-5 py-2.5 text-sm font-medium text-zinc-700 shadow-sm transition hover:bg-zinc-50 hover:text-zinc-900" onClick={(e) => e.stopPropagation()}>
                  <FileText className="h-4 w-4" /> PDF
                </button>
                <button className="flex items-center gap-2 rounded-lg border border-zinc-200 bg-white px-5 py-2.5 text-sm font-medium text-zinc-700 shadow-sm transition hover:bg-zinc-50 hover:text-zinc-900" onClick={(e) => e.stopPropagation()}>
                  <LinkIcon className="h-4 w-4" /> Paste URL
                </button>
              </div>

              <button className="inline-flex items-center gap-2 rounded-lg bg-zinc-950 px-8 py-3.5 text-base font-medium text-white shadow-sm transition hover:bg-zinc-800">
                <File className="h-5 w-5" /> Browse Files
              </button>
            </div>
            
            <div className="mt-8 flex items-center gap-2 text-sm font-medium text-zinc-500">
              <ShieldCheck className="h-4 w-4" />
              Your data is secure and encrypted
            </div>
          </motion.div>
        )}

        {state === 'uploading' && (
          <motion.div
            key="uploading"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="flex flex-col items-center py-16"
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
