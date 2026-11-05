import { useSyncExternalStore } from 'react'

import type { UploadProgress } from '@/upload'

const NOTIFY_MS = 120

let transfers: UploadProgress[] = []
let lastNotify = 0
let timer: ReturnType<typeof setTimeout> | null = null
const listeners = new Set<() => void>()

function getTransfers(): UploadProgress[] {
  return transfers
}

function subscribe(listener: () => void): () => void {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

function flush(): void {
  timer = null
  lastNotify = Date.now()
  for (const listener of listeners) listener()
}

function notify(): void {
  if (listeners.size === 0) return
  const wait = lastNotify + NOTIFY_MS - Date.now()
  if (wait <= 0) {
    if (timer) clearTimeout(timer)
    flush()
    return
  }
  if (!timer) timer = setTimeout(flush, wait)
}

function sameList(a: UploadProgress[], b: UploadProgress[]): boolean {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    const x = a[i]!
    const y = b[i]!
    if (
      x.id !== y.id ||
      x.sent !== y.sent ||
      x.total !== y.total ||
      x.status !== y.status ||
      x.speed !== y.speed ||
      x.error !== y.error
    )
      return false
  }
  return true
}

export function updateTransfers(
  update: (current: UploadProgress[]) => UploadProgress[],
): void {
  const next = update(transfers)
  if (next === transfers || sameList(next, transfers)) return
  const wasEmpty = transfers.length === 0
  transfers = next
  if (next.length === 0 && !wasEmpty) {
    if (timer) clearTimeout(timer)
    flush()
    return
  }
  notify()
}

export function clearTransfers(): void {
  updateTransfers(() => [])
}

export function useTransfers(): UploadProgress[] {
  return useSyncExternalStore(subscribe, getTransfers)
}
