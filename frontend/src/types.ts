import type * as Model from '../bindings/desktopcal/internal/model/models'

/**
 * 前端 DTO 直接派生自 Wails 生成的模型，后端字段变化经 `make bindings`
 * 重新生成后会自动反映到这里，不存在手写副本漂移的问题。
 */
export type Event = Model.Event
export type LunarDay = Model.LunarDay
export type Settings = Model.Settings

/** 去掉数组元素上的 null —— 生成器对 Go 的 []T 一律保守地标成 (T | null)[]。 */
type NonNullElements<T> = T extends (infer U | null)[] ? U[] : never

export interface DayInfo extends Omit<Model.DayInfo, 'events'> {
  /** 后端为 *LunarDay，空占位格会返回 null */
  lunar: Model.LunarDay | null
  /** 收窄为 Event[]：后端返回的是 []Event，元素不会是 null */
  events: NonNullElements<Model.DayInfo['events']>
}
