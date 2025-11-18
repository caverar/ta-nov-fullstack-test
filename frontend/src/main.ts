import './assets/css/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import ui from '@nuxt/ui/vue-plugin'

import App from './App.vue'


const pinia = createPinia()
const app = createApp(App)
app.use(pinia)

app.use(createRouter({
  routes: [
    { path: '/', component: () => import('./pages/index.vue') },
    { path: '/ratings', component: () => import('./pages/ratings.vue') },
    // {
    //   path: '/settings',
    //   component: () => import('./pages/settings.vue'),
    //   children: [
    //     { path: '', component: () => import('./pages/settings/index.vue') },
    //     { path: 'members', component: () => import('./pages/settings/members.vue') },
    //     { path: 'notifications', component: () => import('./pages/settings/notifications.vue') },
    //     { path: 'security', component: () => import('./pages/settings/security.vue') },
    //   ]
    // }
  ],
  history: createWebHistory()
}))

app.use(ui)

app.mount('#app')
