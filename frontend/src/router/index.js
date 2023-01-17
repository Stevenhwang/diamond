import Vue from 'vue'
import Router from 'vue-router'

// vue router 版本升级 https://github.com/vuejs/vue-router/issues/2881#issuecomment-520554378
const originalPush = Router.prototype.push
Router.prototype.push = function push(location, onResolve, onReject) {
  if (onResolve || onReject) return originalPush.call(this, location, onResolve, onReject)
  return originalPush.call(this, location).catch(err => err)
}

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [{
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: 'Dashboard', icon: 'dashboard' }
    }]
  },

  {
    path: '/terminal',
    component: () => import('@/views/terminal/index'),
    hidden: true
  },

  {
    path: '/grafana',
    component: Layout,
    children: [
      {
        path: 'https://grafana.com/',
        meta: { title: 'Grafana', icon: 'link' }
      }
    ]
  },

]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/server',
    name: 'server',
    component: Layout,
    children: [
      {
        path: 'index',
        // name: 'server',
        component: () => import('@/views/server/index'),
        meta: { title: 'Server', icon: 'el-icon-cpu' }
      }
    ]
  },
  {
    path: '/task',
    name: 'task',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/task/index'),
        meta: { title: 'Task', icon: 'el-icon-bangzhu' }
      }
    ]
  },
  {
    path: '/history',
    name: 'history',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/history/index'),
        meta: { title: 'History', icon: 'el-icon-data-board' }
      }
    ]
  },
  {
    path: '/credential',
    name: 'credential',
    component: Layout,
    children: [
      {
        path: 'index',
        // name: 'credential',
        component: () => import('@/views/credential/index'),
        meta: { title: 'Credential', icon: 'el-icon-key' }
      }
    ]
  },
  {
    path: '/record',
    name: 'record',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/record/index'),
        meta: { title: 'Record', icon: 'el-icon-notebook-2' }
      }
    ]
  },
  {
    path: '/user',
    name: 'user',
    component: Layout,
    children: [
      {
        path: 'index',
        // name: 'user',
        component: () => import('@/views/user/index'),
        meta: { title: 'User', icon: 'user' }
      }
    ]
  },

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
