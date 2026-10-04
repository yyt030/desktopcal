export interface Event {
  id: string
  title: string
  description: string
  date: string
  startTime: string
  endTime: string
  color: string
  completed: boolean
  createdAt: string
  updatedAt: string
}

export interface LunarDay {
  year: number
  month: number
  day: number
  isLeapMonth: boolean
  dayName: string
  monthName: string
  animal: string
  ganZhi: string
}

export interface DayInfo {
  solarDate: string
  lunar: LunarDay
  solarTerm: string
  events: Event[]
  isToday: boolean
  isCurrentMonth: boolean
  isEmpty: boolean
  bgColor: string
}

export interface Settings {
  opacity: number
  width: number
  height: number
  theme: string
  alwaysTop: boolean
  fontSize: number
  firstDayOfWeek: number
  showLunar: boolean
  showSolarTerm: boolean
  layoutMode: string
}
