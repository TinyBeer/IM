import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '@/store';
import Home from '../views/Home.vue';
import userRoutes from './module/user';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'Home',
    // meta: {
    //   auth: true, // 本路由需要验证
    // },
    component: Home,
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue'),
  },
  ...userRoutes,
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

// GOOD
router.beforeEach((to, from, next) => {
  if (to.meta.auth) {
    // 判断是否登录
    if (store.state.token) {
      // 需要判断token的有效性
      next();
    } else {
      router.push({ name: 'login' });
    }
  } else {
    next();
  }
});

export default router;
