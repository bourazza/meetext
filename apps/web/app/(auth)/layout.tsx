export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/30 px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary">Meetext</h1>
          <p className="text-muted-foreground mt-1 text-sm">AI-powered meeting intelligence</p>
        </div>
        {children}
      </div>
    </div>
  )
}
