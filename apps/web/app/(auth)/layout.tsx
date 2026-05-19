import { UploadCloud } from 'lucide-react'

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-zinc-50/50 px-4 py-12 text-zinc-950 sm:px-6 lg:px-8">
      <div className="mb-8 flex flex-col items-center">
        <div className="mb-4 flex h-10 w-10 items-center justify-center rounded bg-zinc-900 text-white shadow-sm">
          <UploadCloud className="h-6 w-6" />
        </div>
        <div className="text-2xl font-semibold tracking-tight text-zinc-950">Meetext</div>
      </div>
      <div className="w-full max-w-[440px]">
        {children}
      </div>
    </main>
  )
}
