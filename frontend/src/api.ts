import type { Event, DayInfo, LunarDay, Settings } from './types'

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
          SetDayBackground(date: string, color: string): Promise<void>
          ExportData(): Promise<string>
          ImportData(data: string): Promise<void>
          GetSettings(): Promise<Settings>
          SaveSettings(settings: Settings): Promise<void>
          HideWindow(): Promise<void>
          ShowWindow(): Promise<void>
          QuitApp(): Promise<void>
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

export async function setDayBackground(date: string, color: string): Promise<void> {
  return api!.SetDayBackground(date, color)
}

export async function exportData(): Promise<string> {
  return api!.ExportData()
}

export async function importData(data: string): Promise<void> {
  return api!.ImportData(data)
}

export async function getSettings(): Promise<Settings> {
  return api!.GetSettings()
}

export async function saveSettings(settings: Settings): Promise<void> {
  return api!.SaveSettings(settings)
}

export async function hideWindow(): Promise<void> {
  return api!.HideWindow()
}

export async function showWindow(): Promise<void> {
  return api!.ShowWindow()
}

export async function quitApp(): Promise<void> {
  return api!.QuitApp()
}
