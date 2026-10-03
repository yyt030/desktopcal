import * as App from '../bindings/desktopcal/app'
import type { DayInfo, Event, LunarDay, Settings } from './types'

export async function getMonthDays(year: number, month: number): Promise<DayInfo[]> {
  return App.GetMonthDays(year, month) as unknown as DayInfo[]
}

export async function getEventsByDate(date: string): Promise<Event[]> {
  return App.GetEventsByDate(date) as unknown as Event[]
}

export async function createEvent(event: Event): Promise<Event> {
  return App.CreateEvent(event) as unknown as Event
}

export async function updateEvent(event: Event): Promise<void> {
  return App.UpdateEvent(event)
}

export async function deleteEvent(id: string): Promise<void> {
  return App.DeleteEvent(id)
}

export async function getTodayLunar(): Promise<LunarDay> {
  return App.GetTodayLunar() as unknown as LunarDay
}

export async function getTodaySolarTerm(): Promise<string> {
  return App.GetTodaySolarTerm()
}

export async function setDayBackground(date: string, color: string): Promise<void> {
  return App.SetDayBackground(date, color)
}

export async function exportData(): Promise<string> {
  return App.ExportData()
}

export async function importData(data: string): Promise<void> {
  return App.ImportData(data)
}

export async function getSettings(): Promise<Settings> {
  return App.GetSettings() as unknown as Settings
}

export async function saveSettings(settings: Settings): Promise<void> {
  return App.SaveSettings(settings)
}

export async function hideWindow(): Promise<void> {
  return App.HideWindow()
}

export async function showWindow(): Promise<void> {
  return App.ShowWindow()
}

export async function quitApp(): Promise<void> {
  return App.QuitApp()
}

export async function setWindowSize(width: number, height: number): Promise<void> {
  return App.SetWindowSize(width, height)
}

export async function setWindowAlwaysOnTop(onTop: boolean): Promise<void> {
  return App.SetWindowAlwaysOnTop(onTop)
}
