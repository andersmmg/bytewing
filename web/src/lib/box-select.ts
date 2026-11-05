import {
  useCallback,
  useRef,
  type MouseEvent as ReactMouseEvent,
  type RefObject,
} from 'react'

// data-control marks a bare control outside any button
export const INTERACTIVE_SELECTOR = '[data-entry-path], button, a, input, textarea, select, [data-control]'

const MIN_MOVE = 4

/** how long a click after a drag is ignored */
const SWALLOW_MS = 300

export const FADE_MS = 150

type Box = { left: number; top: number; right: number; bottom: number }

type Target = { path: string } & Box

type Drag = {
  startX: number
  startY: number
  x: number
  y: number
  container: HTMLElement
  area: Box
  targets: Target[]
  base: string[]
  active: boolean
  frame: number
}

type Options = {
  snapshot: () => string[]
  /** runs every drag frame */
  apply: (paths: string[]) => void
  overlay: RefObject<HTMLElement | null>
}

export type BoxSelect = {
  onMouseDown: (e: ReactMouseEvent) => void
  click: (then: () => void) => void
}

/** measures every entry row in the container, in screen space */
function measure(container: HTMLElement): Target[] {
  const out: Target[] = []
  for (const el of container.querySelectorAll('[data-entry-path]')) {
    const path = el.getAttribute('data-entry-path')
    if (!path) continue
    const r = el.getBoundingClientRect()
    out.push({ path, left: r.left, top: r.top, right: r.right, bottom: r.bottom })
  }
  return out
}

/** the container's own rect, used to clip the painted box */
function areaOf(container: HTMLElement): Box {
  const r = container.getBoundingClientRect()
  return { left: r.left, top: r.top, right: r.right, bottom: r.bottom }
}

/** re-measures after the page moves under a drag */
function recapture(d: Drag) {
  d.area = areaOf(d.container)
  d.targets = measure(d.container)
}

/** returns the paths of every target the box overlap */
function touched(targets: Target[], box: Box): string[] {
  const out: string[] = []
  for (const t of targets) {
    if (t.left <= box.right && t.right >= box.left && t.top <= box.bottom && t.bottom >= box.top)
      out.push(t.path)
  }
  return out
}

/**
 * useBoxSelect drags a selection box over `[data-entry-path]` rows.
 */
export function useBoxSelect({ snapshot, apply, overlay }: Options): BoxSelect {
  const opts = useRef({ snapshot, apply })
  opts.current = { snapshot, apply }

  const drag = useRef<Drag | null>(null)
  const swallowUntil = useRef(0)

  /** draws the box and applies the selection it covers */
  const paint = useCallback((d: Drag) => {
    const box: Box = {
      left: Math.min(d.startX, d.x),
      top: Math.min(d.startY, d.y),
      right: Math.max(d.startX, d.x),
      bottom: Math.max(d.startY, d.y),
    }
    // clipped so it never paints over the sidebar or header
    const area = d.area
    const left = Math.max(box.left, area.left)
    const top = Math.max(box.top, area.top)
    const right = Math.min(box.right, area.right)
    const bottom = Math.min(box.bottom, area.bottom)
    const el = overlay.current
    const visible = left < right && top < bottom
    if (el) {
      if (visible) {
        // transform, not left/top: moving the box must not force layout
        el.style.width = `${right - left}px`
        el.style.height = `${bottom - top}px`
        el.style.transform = `translate3d(${left}px, ${top}px, 0)`
        el.style.opacity = '1'
      } else {
        el.style.opacity = '0'
      }
    }
    opts.current.apply([...d.base, ...touched(d.targets, box)])
  }, [overlay])

  /** tracks the pointer, coalescing paints into one animation frame */
  const move = useCallback(
    (e: globalThis.MouseEvent) => {
      const d = drag.current
      if (!d) return
      d.x = e.clientX
      d.y = e.clientY
      if (!d.active && Math.hypot(d.x - d.startX, d.y - d.startY) < MIN_MOVE) return
      d.active = true
      if (d.frame) return
      d.frame = requestAnimationFrame(() => {
        d.frame = 0
        const cur = drag.current
        if (cur) paint(cur)
      })
    },
    [paint],
  )

  /** re-measures when the page scrolls or resizes mid-drag */
  const scroll = useCallback(() => {
    const d = drag.current
    if (d) recapture(d)
  }, [])

  /** tears the drag down and commits the final selection */
  const end = useCallback(() => {
    const d = drag.current
    if (!d) return
    window.removeEventListener('mousemove', move)
    window.removeEventListener('mouseup', end)
    window.removeEventListener('blur', end)
    window.removeEventListener('scroll', scroll, true)
    window.removeEventListener('resize', scroll)
    if (d.frame) cancelAnimationFrame(d.frame)
    d.frame = 0
    drag.current = null
    if (!d.active) return
    recapture(d)
    paint(d)
    const el = overlay.current
    if (el) el.style.opacity = '0'
    swallowUntil.current = Date.now() + SWALLOW_MS
  }, [move, paint, overlay, scroll])

  /** starts a drag, unless the press belongs to a control or a scrollbar */
  const onMouseDown = useCallback(
    (e: ReactMouseEvent) => {
      if (e.button !== 0) return
      const target = e.target instanceof Element ? e.target : null
      if (target?.closest(INTERACTIVE_SELECTOR)) return
      const container = e.currentTarget
      if (!(container instanceof HTMLElement)) return
      // scrollbars dispatch to the scroller itself, so test the gutter all the way up
      for (let el: Element | null = target; el; el = el.parentElement) {
        const r = el.getBoundingClientRect()
        if (e.clientX - r.left >= el.clientWidth || e.clientY - r.top >= el.clientHeight) return
        if (el === container) break
      }
      e.preventDefault()
      const d: Drag = {
        startX: e.clientX,
        startY: e.clientY,
        x: e.clientX,
        y: e.clientY,
        container,
        area: areaOf(container),
        targets: [],
        base: e.shiftKey ? opts.current.snapshot() : [],
        active: false,
        frame: 0,
      }
      recapture(d)
      drag.current = d
      window.addEventListener('mousemove', move)
      window.addEventListener('mouseup', end)
      window.addEventListener('blur', end)
      window.addEventListener('scroll', scroll, true)
      window.addEventListener('resize', scroll)
    },
    [move, end, scroll],
  )

  /** runs a click handler unless that click ended a drag */
  const click = useCallback((then: () => void) => {
    if (Date.now() < swallowUntil.current) {
      swallowUntil.current = 0
      return
    }
    then()
  }, [])

  return { onMouseDown, click }
}
