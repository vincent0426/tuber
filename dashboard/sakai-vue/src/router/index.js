import { createRouter, createWebHashHistory } from 'vue-router';
import AppLayout from '@/layout/AppLayout.vue';
import store from '@/store';
const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            // With sidebar & topbar
            component: AppLayout,
            children: [
                {
                    path: '/passenger/home',
                    name: 'PassengerHome',
                    component: () => import('@/views/pages/passenger/PassengerHome.vue')
                },
                {
                    path: '/passenger/history',
                    name: 'PassengerHistory',
                    component: () => import('@/views/pages/passenger/PassengerHistory.vue')
                },
                {
                    path: '/logo',
                    name: 'logo',
                    component: () => import('@/views/pages/logo.vue')
                },
                {
                    path: '/login',
                    name: 'UIlogin',
                    component: () => import('@/views/pages/login.vue')
                },
                {
                    path: '/profile',
                    name: 'profile',
                    component: () => import('@/views/pages/profile.vue')
                },
                {
                    path: '/uikit/formlayout',
                    name: 'formlayout',
                    component: () => import('@/views/uikit/FormLayout.vue')
                },
                {
                    path: '/passenger/search', // /search?q=yourSearchQuery
                    name: 'PassengerSearch',
                    component: () => import('@/views/pages/passenger/PassengerSearch.vue')
                },
                {
                    path: '/setting',
                    name: 'PassengerSetting',
                    component: () => import('@/views/pages/utils/Setting.vue')
                },
                {
                    path: '/passenger/mytrip',
                    name: 'MyTrip',
                    component: () => import('@/views/pages/passenger/MyTrip.vue')
                },
                {
                    path: '/passenger/favorite',
                    name: 'Favorite',
                    component: () => import('@/views/pages/passenger/Favorite.vue')
                },
                {
                    path: '/driver/home',
                    name: 'DriverHome',
                    component: () => import('@/views/pages/driver/DriverHome.vue')
                },
                {
                    path: '/driver/history',
                    name: 'DriverHistory',
                    component: () => import('@/views/pages/driver/DriverHistory.vue')
                },
                {
                    path: '/driver/create',
                    name: 'DriverCreate',
                    component: () => import('@/views/pages/driver/DriverCreate.vue')
                },
                {
                    path: '/setting',
                    name: 'DriverSetting',
                    component: () => import('@/views/pages/utils/Setting.vue')
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
                {
                    path: '/auth/chatList',
                    name: 'ChatList',
                    component: () => import('@/views/pages/auth/ChatList.vue')
                },
                {
                    path: '/chatroom/:tripId',
                    name: 'ChatRoom',
                    component: () => import('@/views/pages/auth/ChatRoom.vue'),
                    props: true // 將路由參數作為組件的 props 傳遞
                }
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
            component: () => import('@/views/pages/auth/Register.vue')
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
    ]
});

// Navigation guard to check authentication before each navigation
router.beforeEach(async (to, from, next) => {
    if (to.name !== 'login' && !store.getters.login) {
        next({ name: 'login' });
    } else {
        console.log('have login');
        next(); // Continue with the navigation
    }
});

export default router;
