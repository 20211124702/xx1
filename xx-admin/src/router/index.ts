import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Login',
    component: () => import('../views/login.vue')
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('../views/home.vue'),
    children: [
      {
        path: 'user',
        name: 'User',
        component: () => import('../views/user.vue')
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('../views/role.vue')
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('../views/menu.vue')
      },
      {
        path: 'table',
        name: 'Table',
        component: () => import('../views/table.vue')
      }
    ]
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/register.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const isLogin = !!localStorage.getItem('token')
  if ((to.path !== '/' && to.path !== '/register') && !isLogin) {
    next('/')
  } else if ((to.path === '/' ) && isLogin) {
    next('/home')
  } else {
    next()
  }
})

export default router;