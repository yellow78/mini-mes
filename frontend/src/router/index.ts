// Vue Router 路由設定
import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import LotView from '../views/LotView.vue'
import SpcView from '../views/SpcView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/',     component: DashboardView, meta: { title: 'Dashboard' } },
    { path: '/lots', component: LotView,       meta: { title: 'Lot / WIP' } },
    { path: '/spc',  component: SpcView,       meta: { title: 'SPC' }       },
  ],
})

router.afterEach((to) => {
  document.title = `${to.meta.title ?? 'MES'} — Mini-MES`
})

export default router
