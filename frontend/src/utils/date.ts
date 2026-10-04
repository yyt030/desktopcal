// 日期工具：统一使用「本地时区解析」，避免 new Date('YYYY-MM-DD') 被当作 UTC
// 而在非东八区环境导致取到的日期偏移一天。

/** 解析 'YYYY-MM-DD' 为本地时区的 Date */
export function parseLocalDate(date: string): Date {
  const [y, m, d] = date.split('-').map(Number)
  return new Date(y || 1970, (m || 1) - 1, d || 1)
}

/** 格式化为 'YYYY-MM-DD'（本地时区） */
export function formatDate(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}
