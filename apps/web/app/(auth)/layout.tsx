import { Compass } from 'lucide-react'

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-zinc-50/40 px-4 py-12 text-zinc-950 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* Background soft glow accents */}
      <div className="absolute top-0 left-1/4 w-[500px] h-[500px] bg-indigo-500/5 rounded-full blur-[120px] pointer-events-none" />
      <div className="absolute bottom-0 right-1/4 w-[400px] h-[400px] bg-purple-500/5 rounded-full blur-[100px] pointer-events-none" />

      <div className="mb-8 flex flex-col items-center relative z-10">
        <div className="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-600 text-white shadow-lg shadow-indigo-600/25">
          <Compass className="h-6 w-6 animate-pulse" />
        </div>
        <div className="text-3xl font-extrabold tracking-tight text-zinc-900 flex items-center gap-1.5">
          Meetext
          <span className="rounded bg-indigo-950 px-1.5 py-0.5 text-[10px] font-semibold text-indigo-400 border border-indigo-900">
            MVP
          </span>
        </div>
      </div>
      <div className="w-full max-w-[480px] relative z-10">
        {children}
      </div>
    </main>
  )
}
