export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <main className="min-h-screen bg-[#f7f8fb] text-zinc-950">
      <div className="grid min-h-screen lg:grid-cols-[1fr_560px]">
        <section className="relative hidden overflow-hidden px-10 py-10 lg:flex lg:flex-col lg:justify-between">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,rgba(15,23,42,0.06)_1px,transparent_1px),linear-gradient(to_bottom,rgba(15,23,42,0.06)_1px,transparent_1px)] bg-[size:42px_42px]" />
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,rgba(16,185,129,0.16),transparent_34%),radial-gradient(circle_at_70%_80%,rgba(14,165,233,0.12),transparent_30%)]" />
          <div className="relative">
            <div className="inline-flex h-9 items-center rounded-md border border-zinc-200 bg-white/80 px-3 text-sm font-semibold shadow-sm">
              Meetext
            </div>
          </div>
          <div className="relative max-w-xl pb-10">
            <p className="mb-4 text-sm font-medium uppercase tracking-[0.18em] text-zinc-500">Meeting intelligence for freelancers</p>
            <h2 className="text-5xl font-semibold tracking-normal text-zinc-950">
              Turn client calls into clear docs, tasks, and project reports.
            </h2>
            <p className="mt-5 max-w-lg text-base leading-7 text-zinc-600">
              A calm workspace for capturing decisions, summarizing meetings, and keeping every client project moving.
            </p>
          </div>
        </section>

        <section className="flex min-h-screen items-center justify-center px-4 py-8 sm:px-6">
          <div className="w-full max-w-[440px]">
            <div className="mb-8 text-center lg:hidden">
              <div className="mb-3 text-2xl font-semibold">Meetext</div>
              <p className="text-sm text-zinc-600">AI-powered meeting intelligence</p>
            </div>
            {children}
          </div>
        </section>
      </div>
    </main>
  )
}
