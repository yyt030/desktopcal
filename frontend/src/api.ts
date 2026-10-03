import type { Event, DayInfo, LunarDay } from './types'

declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetMonthDays(year: number, month: number): Promise<DayInfo[]>
          GetEventsByDate(date: string): Promise<Event[]>
          CreateEvent(event: Event): Promise<Event>
          UpdateEvent(event: Event): Promise<void>
          DeleteEvent(id: string): Promise<void>
          GetTodayLunar(): Promise<LunarDay>
          GetTodaySolarTerm(): Promise<string>
        }
      }
    }
  }
}

const api = window.go?.main?.App

export async function getMonthDays(year: number, month: number): Promise<DayInfo[]> {
  return api!.GetMonthDays(year, month)
}

export async function getEventsByDate(date: string): Promise<Event[]> {
  return api!.GetEventsByDate(date)
}

export async function createEvent(event: Event): Promise<Event> {
  return api!.CreateEvent(event)
}

export async function updateEvent(event: Event): Promise<void> {
  return api!.UpdateEvent(event)
}

export async function deleteEvent(id: string): Promise<void> {
  return api!.DeleteEvent(id)
}

export async function getTodayLunar(): Promise<LunarDay> {
  return api!.GetTodayLunar()
}

export async function getTodaySolarTerm(): Promise<string> {
  return api!.GetTodaySolarTerm()
}
