import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

function showErr(msg: string) {
  document.body.innerHTML =
    '<pre style="color:#fff;background:#900;padding:20px;font-size:16px;white-space:pre-wrap;margin:0;height:100vh">' +
    msg + '</pre>'
}
window.addEventListener('error', (e) => showErr('ERROR: ' + e.message + '\n' + (e.error?.stack || '')))
window.addEventListener('unhandledrejection', (e) => {
  const r: any = e.reason
  showErr('REJECT: ' + String(r?.stack || r?.message || r))
})

const app = createApp(App)
app.use(createPinia())
app.mount('#app')

