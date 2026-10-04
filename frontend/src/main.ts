import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

// 非破坏性错误提示：向页面追加一个可关闭的浮层，而不是清空整个 body。
// （清空 body 会让应用彻底不可用，且丢失用户当前编辑内容。）
function showErr(msg: string) {
  const id = '__app_error__'
  let el = document.getElementById(id)
  if (!el) {
    el = document.createElement('div')
    el.id = id
    el.title = '点击关闭'
    el.style.cssText = [
      'position:fixed',
      'left:0',
      'right:0',
      'bottom:0',
      'z-index:99999',
      'max-height:40vh',
      'overflow:auto',
      'background:rgba(153,0,0,0.95)',
      'color:#fff',
      'padding:12px 16px',
      'font-size:13px',
      'font-family:ui-monospace,SFMono-Regular,Menlo,monospace',
      'white-space:pre-wrap',
      'cursor:pointer',
      'border-top:2px solid #f66',
    ].join(';')
    el.addEventListener('click', () => el?.remove())
    document.body.appendChild(el)
  }
  el.textContent = '发生错误（点击关闭）：\n' + msg
}

window.addEventListener('error', (e) => {
  showErr(e.message + '\n' + (e.error?.stack ?? ''))
})
window.addEventListener('unhandledrejection', (e) => {
  const reason = e.reason as { stack?: string; message?: string } | undefined
  showErr(String(reason?.stack ?? reason?.message ?? reason))
})

const app = createApp(App)
app.use(createPinia())
app.mount('#app')
