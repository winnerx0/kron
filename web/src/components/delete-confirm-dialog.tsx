'use client'

import { X } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  jobName: string
  onConfirm: () => void
  loading?: boolean
}

export function DeleteConfirmDialog({ open, onOpenChange, jobName, onConfirm, loading = false }: Props) {
  if (!open) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 animate-fade-in"
      style={{ backgroundColor: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)' }}
      onClick={e => { if (e.target === e.currentTarget) onOpenChange(false) }}
    >
      <div className="w-full max-w-sm bg-card border border-border rounded-lg shadow-2xl animate-slide-up flex flex-col">
        {/* Header */}
        <div className="flex items-start justify-between px-6 pt-5 pb-4 border-b border-border shrink-0">
          <div>
            <h2 className="text-base font-semibold tracking-tight text-card-foreground">
              Delete job?
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              This action cannot be undone
            </p>
          </div>
          <button
            onClick={() => onOpenChange(false)}
            disabled={loading}
            className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors disabled:opacity-50"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Body */}
        <div className="px-6 py-4 flex-1">
          <p className="text-sm text-foreground">
            Are you sure you want to delete <span className="font-medium">{jobName}</span>? Any scheduled executions will be canceled.
          </p>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-2 px-6 py-4 border-t border-border shrink-0">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onOpenChange(false)}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={onConfirm}
            disabled={loading}
          >
            {loading ? 'Deleting…' : 'Delete'}
          </Button>
        </div>
      </div>
    </div>
  )
}
