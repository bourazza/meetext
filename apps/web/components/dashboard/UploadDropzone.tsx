'use client'

import React, { useRef, useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
  UploadCloud, FileText, Music, Video, ShieldCheck,
  CheckCircle2, Loader2, Clock
} from 'lucide-react'

type UploadState = 'idle' | 'uploading' | 'processing' | 'complete'

interface Props {
  onFileSelected?: (file: File) => void
  disabled?: boolean
}

/**
 * UploadDropzone — MVP supported formats:
 *   ✅ PDF  (fully supported — text extraction → Ollama AI)
 *   🔜 Audio / Video  (coming soon — requires Whisper + n8n integration)
 */
export function UploadDropzone({ onFileSelected, disabled }: Props) {
  const [state, setState] = useState<UploadState>('idle')
  const [progress, setProgress] = useState(0)
  const [dragging, setDragging] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFile = useCallback(
    (file: File | null | undefined) => {
      if (!file || disabled) return

      // Reject audio/video — not supported yet
      if (
        file.type.startsWith('audio/') ||
        file.type.startsWith('video/') ||
        /\.(mp3|wav|mp4|mov|avi|webm|ogg)$/i.test(file.name)
      ) {
        // Let parent handle the toast; just return
        onFileSelected?.(file)
        return
      }

      if (file.type !== 'application/pdf' && !/\.pdf$/i.test(file.name)) {
        onFileSelected?.(file) // parent will show unsupported error
        return
      }

      onFileSelected?.(file)
    },
    [disabled, onFileSelected]
  )

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setDragging(false)
    handleFile(e.dataTransfer.files?.[0])
  }

  return (
    <div className="w-full">
      {/* Hidden real file input — PDF only */}
      <input
        ref={fileInputRef}
        type="file"
        accept="application/pdf"
        className="hidden"
        onChange={(e) => handleFile(e.target.files?.[0])}
      />

      <AnimatePresence mode="wait">
        {state === 'idle' && (
          <motion.div
            key="idle"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="flex flex-col items-center"
          >
            {/* Drop Zone */}
            <div
              onDragOver={(e) => { e.preventDefault(); setDragging(true) }}
              onDragLeave={() => setDragging(false)}
              onDrop={handleDrop}
              onClick={() => fileInputRef.current?.click()}
              className={`group w-full cursor-pointer rounded-2xl border-2 border-dashed py-20 text-center transition-all duration-200 ${
                dragging
                  ? 'border-indigo-500 bg-indigo-50/30'
                  : 'border-zinc-200 bg-zinc-50/50 hover:border-indigo-400 hover:bg-indigo-50/10'
              }`}
            >
              <div className="mb-4 inline-flex items-center justify-center rounded-2xl bg-white p-4 text-zinc-700 shadow-sm border border-zinc-100 group-hover:scale-105 transition-all">
                <UploadCloud className="h-8 w-8 text-indigo-500" />
              </div>
              <h3 className="mb-1 text-base font-bold text-zinc-900 group-hover:text-indigo-600 transition-colors">
                Drop your PDF here or click to browse
              </h3>
              <p className="mb-6 text-xs text-zinc-400">PDF meeting documents up to 250 MB</p>

              {/* Format badges */}
              <div className="flex flex-wrap items-center justify-center gap-3">
                {/* Supported */}
                <span className="flex items-center gap-1.5 rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-1.5 text-xs font-medium text-emerald-700 shadow-sm">
                  <FileText className="h-3 w-3" />
                  PDF
                  <span className="ml-1 rounded bg-emerald-200/60 px-1 py-px text-[9px] font-semibold uppercase">Supported</span>
                </span>

                {/* Coming Soon */}
                <span className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-zinc-100 px-3 py-1.5 text-xs text-zinc-400 shadow-sm cursor-default select-none">
                  <Music className="h-3 w-3" />
                  Audio
                  <span className="ml-1 flex items-center gap-0.5 rounded bg-amber-100 px-1 py-px text-[9px] font-semibold text-amber-600">
                    <Clock className="h-2 w-2" /> Soon
                  </span>
                </span>
                <span className="flex items-center gap-1.5 rounded-lg border border-zinc-200 bg-zinc-100 px-3 py-1.5 text-xs text-zinc-400 shadow-sm cursor-default select-none">
                  <Video className="h-3 w-3" />
                  Video
                  <span className="ml-1 flex items-center gap-0.5 rounded bg-amber-100 px-1 py-px text-[9px] font-semibold text-amber-600">
                    <Clock className="h-2 w-2" /> Soon
                  </span>
                </span>
              </div>
            </div>

            <div className="mt-4 flex items-center gap-2 text-xs font-medium text-zinc-400">
              <ShieldCheck className="h-4 w-4 text-emerald-500" />
              Enterprise-grade encryption · HIPAA compliant
            </div>
          </motion.div>
        )}

        {state === 'uploading' && (
          <motion.div
            key="uploading"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0 }}
            className="flex flex-col items-center py-16"
          >
            <FileText className="mb-4 h-10 w-10 text-indigo-400" />
            <div className="mb-2 text-sm font-semibold text-zinc-900">Uploading PDF...</div>
            <div className="mb-2 w-full max-w-md overflow-hidden rounded-full bg-zinc-100">
              <motion.div
                className="h-2 bg-indigo-600 rounded-full"
                initial={{ width: 0 }}
                animate={{ width: `${progress}%` }}
                transition={{ ease: 'linear' }}
              />
            </div>
            <div className="text-xs text-zinc-400">{progress}%</div>
          </motion.div>
        )}

        {state === 'processing' && (
          <motion.div
            key="processing"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0 }}
            className="flex flex-col items-center py-16 text-center"
          >
            <Loader2 className="mb-4 h-8 w-8 animate-spin text-indigo-600" />
            <h3 className="mb-1 text-lg font-semibold text-zinc-900">AI is analysing your document…</h3>
            <p className="text-sm text-zinc-400">Extracting tasks, decisions, risks, and generating documentation.</p>
          </motion.div>
        )}

        {state === 'complete' && (
          <motion.div
            key="complete"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0 }}
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
            <p className="text-sm text-zinc-400">Your meeting documentation is ready below.</p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
