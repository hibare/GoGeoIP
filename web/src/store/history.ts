import { defineStore } from "pinia"
import { ref, computed } from "vue"

export interface LookupHistory {
  id: string
  ip: string
  organization: string
  location: string
  countryCode?: string
  timestamp: number
  first_lookup: number
  times_looked: number
}

const STORAGE_KEY = "waypoint_lookup_history"
const MAX_ENTRIES = 100

export const useHistoryStore = defineStore("history", () => {
  const history = ref<LookupHistory[]>([])

  const sortedHistory = computed(() => {
    return [...history.value].sort((a, b) => b.timestamp - a.timestamp)
  })

  function loadFromStorage() {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        history.value = JSON.parse(stored)
      } catch {
        history.value = []
      }
    }
  }

  function saveToStorage() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history.value))
  }

  function addEntry(ip: string, organization: string, location: string, countryCode?: string) {
    const existingIndex = history.value.findIndex((h) => h.ip === ip)

    if (existingIndex !== -1) {
      const existing = history.value[existingIndex]
      existing.timestamp = Date.now()
      existing.organization = organization
      existing.location = location
      existing.countryCode = countryCode
      existing.times_looked += 1
      history.value = [...history.value]
    } else {
      const now = Date.now()
      const entry: LookupHistory = {
        id: crypto.randomUUID(),
        ip,
        organization,
        location,
        countryCode,
        timestamp: now,
        first_lookup: now,
        times_looked: 1,
      }
      history.value.unshift(entry)

      if (history.value.length > MAX_ENTRIES) {
        history.value = history.value.slice(0, MAX_ENTRIES)
      }
    }

    saveToStorage()
  }

  function removeEntry(id: string) {
    history.value = history.value.filter((h) => h.id !== id)
    saveToStorage()
  }

  function clearHistory() {
    history.value = []
    saveToStorage()
  }

  function hasIp(ip: string): boolean {
    return history.value.some((h) => h.ip === ip)
  }

  loadFromStorage()

  return {
    history: sortedHistory,
    addEntry,
    removeEntry,
    clearHistory,
    hasIp,
    loadFromStorage,
  }
})
