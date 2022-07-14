import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import TeamsView from '../views/TeamsView.vue'
import TeammatesView from '../views/TeammatesView.vue'
import HistoryView from '../views/HistoryView.vue'
import RegisterView from '../views/RegisterView.vue'
import LogoutView from '../views/LogoutView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/logout',
    name: 'logout',
    component: LogoutView
  },
  {
    path: '/teams',
    name: 'teams',
    component: TeamsView
  },
  {
    path: '/teams/:teamId/teammates',
    name: 'teammates',
    component: TeammatesView,
    props: ({params}) => ({teamId: Number.parseInt(params.teamId, 10) || 0})
  },
  {
    path: '/teams/:teamId/history',
    name: 'history',
    component: HistoryView,
    props: ({params}) => ({teamId: Number.parseInt(params.teamId, 10) || 0})
  },
  {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/AboutView.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),//createWebHistory(process.env.BASE_URL),
  routes
})

export default router
