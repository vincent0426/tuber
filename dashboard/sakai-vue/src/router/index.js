import { createRouter, createWebHashHistory } from 'vue-router';
import AppLayout_Customer from '@/layout/AppLayout_Customer.vue';
import AppLayout_Driver from '@/layout/AppLayout_Driver.vue';
import store from '@/store';
const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/customer',
            // With sidebar & topbar
            component: AppLayout_Customer,
            children: [
                {
                    path: '/customer/home',
                    name: 'CustomerHome',
                    component: () => import('@/views/pages/customer/CustomerHome.vue')
                },
                {
                    path: '/customer/history',
                    name: 'CustomerHistory',
                    component: () => import('@/views/pages/Empty.vue')
                },
                {
                    path: '/customer/search', // /search?q=yourSearchQuery
                    name: 'CustomerSearch',
                    component: () => import('@/views/pages/Empty.vue')
                },
                {
                    path: '/customer/setting',
                    name: 'CutomerSetting',
                    component: () => import('@/views/pages/customer/CustomerSetting.vue')
                },
                {
                    path: '/customer/mytrip',
                    name: 'MyTrip',
                    component: () => import('@/views/pages/customer/MyTrip.vue')
                },
                {
                    path: '/customer/favorite',
                    name: 'Favorite',
                    component: () => import('@/views/pages/customer/Favorite.vue')
                }
            ]
        },
        {
            path: '/driver',
            // With sidebar & topbar
            component: AppLayout_Driver,
            children: [
                {
                    path: '/driver/home',
                    name: 'DriverHome',
                    component: () => import('@/views/pages/driver/DriverHome.vue')
                },
                {
                    path: '/driver/history',
                    name: 'DriverHistory',
                    component: () => import('@/views/pages/Empty.vue')
                },
                {
                    path: '/driver/create',
                    name: 'DriverCreate',
                    component: () => import('@/views/pages/Empty.vue')
                },
                {
                    path: '/driver/setting',
                    name: 'DriverSetting',
                    component: () => import('@/views/pages/Empty.vue')
                },
                {
                    path: '/driver/profile',
                    name: 'DriverProfile',
                    component: () => import('@/views/pages/driver/DriverProfile.vue')
                },
                {
                    path: '/driver/applylist',
                    name: 'ApplyList',
                    component: () => import('@/views/pages/driver/ApplyList.vue')
                },
            ]
        },
        {
            path: '/auth/login',
            name: 'login',
            component: () => import('@/views/pages/auth/Login.vue')
        },
        {
            path: '/auth/access',
            name: 'accessDenied',
            component: () => import('@/views/pages/auth/Access.vue')
        },
        {
            path: '/auth/register',
            name: 'Register',
            component: () => import('@/views/pages/auth/Access.vue')
        },
        {
            path: '/chat/:userId',
            name: 'Chatroom',
            component: () => import('@/views/pages/utils/Chat.vue')
        },
        {
            path: '/onTrip/:tripId',
            name: 'OnTrip',
            component: () => import('@/views/pages/utils/OnTrip.vue')
        },
        {
            path: '/TripDetail/:tripId',
            name: 'TripDetail',
            component: () => import('@/views/pages/utils/TripDetail.vue')
        }
    ],
});


// Navigation guard to check authentication before each navigation
router.beforeEach(async (to, from, next) => {
    if (to.name !== 'login' && !store.getters.login) {
        try {
            console.log('check login');
            await store.dispatch("checkLogin");
            next();
        } catch {
            next({ name: "login" });
        }
    } else {
    console.log('have login');
      next(); // Continue with the navigation
    }
});

export default router;
