import React from 'react'
import { UploadDropzone } from '@/components/dashboard/UploadDropzone'

export default function DashboardPage() {
  return (
    <div className="mx-auto max-w-5xl px-8 py-12">
      <div className="mb-12 max-w-2xl">
        <h1 className="mb-3 text-3xl font-semibold tracking-tight text-zinc-950">
          Upload New Meeting
        </h1>
        <p className="text-lg text-zinc-500">
          Drop your recording here, and Meetext will generate a full transcript, extract tasks, and write your project documentation.
        </p>
      </div>
      
      <div className="flex justify-center">
        <UploadDropzone />
      </div>
    </div>
  )
}
