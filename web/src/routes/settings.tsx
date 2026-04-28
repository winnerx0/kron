import { useState } from 'react'
import { useTheme } from '@/lib/theme'
import { Sun, Moon, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

function Section({ title, description, children }: { title: string; description: string; children: React.ReactNode }) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 py-8 border-b border-border last:border-0">
      <div>
        <h2 className="text-sm font-semibold tracking-tight">{title}</h2>
        <p className="text-sm text-muted-foreground mt-1">{description}</p>
      </div>
      <div className="lg:col-span-2 space-y-4">{children}</div>
    </div>
  )
}

export function SettingsPage() {
  const { theme, toggle } = useTheme()
  const [apiUrl, setApiUrl] = useState(import.meta.env.VITE_API_URL || 'http://localhost:5000')
  const [saved, setSaved] = useState(false)

  const handleSave = () => {
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div className="max-w-2xl space-y-0 animate-fade-in">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold tracking-tight">Settings</h1>
        <p className="text-sm text-muted-foreground mt-1">Configure your Kron instance</p>
      </div>

      <div className="bg-card border border-border rounded-lg px-6 divide-y divide-border">
        {/* Appearance */}
        <Section title="Appearance" description="Adjust the visual theme of the dashboard.">
          <div className="flex items-center gap-3">
            <button
              onClick={() => theme === 'dark' && toggle()}
              className={`flex-1 flex items-center gap-3 px-4 py-3 rounded-lg border text-sm transition-all ${
                theme === 'light'
                  ? 'border-primary bg-primary/5 font-medium'
                  : 'border-border text-muted-foreground hover:border-foreground/40'
              }`}
            >
              <Sun className="w-4 h-4" />
              Light
              {theme === 'light' && <Check className="w-3.5 h-3.5 ml-auto" />}
            </button>
            <button
              onClick={() => theme === 'light' && toggle()}
              className={`flex-1 flex items-center gap-3 px-4 py-3 rounded-lg border text-sm transition-all ${
                theme === 'dark'
                  ? 'border-primary bg-primary/5 font-medium'
                  : 'border-border text-muted-foreground hover:border-foreground/40'
              }`}
            >
              <Moon className="w-4 h-4" />
              Dark
              {theme === 'dark' && <Check className="w-3.5 h-3.5 ml-auto" />}
            </button>
          </div>
        </Section>

        {/* API */}
        <Section title="API Connection" description="The URL of your running Kron backend server.">
          <div className="space-y-2">
            <label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">API URL</label>
            <Input
              value={apiUrl}
              onChange={e => setApiUrl(e.target.value)}
              placeholder="http://localhost:5000"
              className="h-9 text-sm font-mono"
            />
            <p className="text-xs text-muted-foreground">
              Changes take effect on next page reload. Set <code className="font-mono bg-muted px-1 rounded">VITE_API_URL</code> env var for a permanent override.
            </p>
          </div>
          <Button size="sm" className="h-8 text-xs gap-2" onClick={handleSave}>
            {saved ? <><Check className="w-3.5 h-3.5" /> Saved</> : 'Save'}
          </Button>
        </Section>

        {/* About */}
        <Section title="About" description="Information about this Kron instance.">
          <div className="space-y-2 text-sm">
            {[
              ['Version', 'v0.1.0'],
              ['API URL', apiUrl],
              ['Refresh interval', '5 seconds'],
            ].map(([k, v]) => (
              <div key={k} className="flex items-center justify-between py-2 border-b border-border last:border-0">
                <span className="text-muted-foreground">{k}</span>
                <span className="font-mono text-xs">{v}</span>
              </div>
            ))}
          </div>
        </Section>
      </div>
    </div>
  )
}
