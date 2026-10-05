import * as App from '../bindings/desktopcal/app'
import type { DayInfo, Event, LunarDay, Settings } from './types'

// 生成器对指针/切片一律附上 `| null`（Go 侧可能为空），此处 await 后用普通
// `as` 收窄。相比 `as unknown as` 仍保留结构重叠检查：后端签名发生不兼容
// 变化时会在这一层直接报错，而不是静默通过。

export async function getMonthDays(year: number, month: number): Promise<DayInfo[]> {
  return (await App.GetMonthDays(year, month)) as DayInfo[]
}

export async function getEventsByDate(date: string): Promise<Event[]> {
  return (await App.GetEventsByDate(date)) as Event[]
}

export async function createEvent(event: Event): Promise<Event> {
  return (await App.CreateEvent(event)) as Event
}

export async function updateEvent(event: Event): Promise<void> {
  return App.UpdateEvent(event)
}

export async function deleteEvent(id: string): Promise<void> {
  return App.DeleteEvent(id)
}

export async function getTodayLunar(): Promise<LunarDay> {
  return (await App.GetTodayLunar()) as LunarDay
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
  return (await App.GetSettings()) as Settings
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

export async function minimiseWindow(): Promise<void> {
  return App.MinimiseWindow()
}

export async function maximiseWindow(): Promise<void> {
  return App.MaximiseWindow()
}

export async function closeWindow(): Promise<void> {
  return App.CloseWindow()
}

export async function getVersion(): Promise<string> {
  return App.GetVersion()
}

export async function getWindowPosition(): Promise<{ x: number; y: number }> {
  const [x, y] = await App.GetWindowPosition()
  return { x, y }
}

export async function setWindowPosition(x: number, y: number): Promise<void> {
  return App.SetWindowPosition(x, y)
}

export async function getWindowSize(): Promise<{ width: number; height: number }> {
  const [width, height] = await App.GetWindowSize()
  return { width, height }
}

export async function getScreenWorkArea(): Promise<{ x: number; y: number; width: number; height: number }> {
  const [x, y, width, height] = await App.GetScreenWorkArea()
  return { x, y, width, height }
}
