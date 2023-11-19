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
                    component: () => import('@/views/UI/logo.vue')
                },
                {
                    path: '/register',
                    name: 'register',
                    component: () => import('@/views/UI/register.vue')
                },
                {
                    path: '/login',
                    name: 'UIlogin',
                    component: () => import('@/views/UI/login.vue')
                },
                {
                    path: '/profile',
                    name: 'profile',
                    component: () => import('@/views/UI/profile.vue')
                },
                {
                    path: '/uikit/formlayout',
                    name: 'formlayout',
                    component: () => import('@/views/uikit/FormLayout.vue')
                },
                {
                    path: '/passenger/search', // /search?q=yourSearchQuery
                    name: 'PassengerSearch',
                    component: () => import('@/views/pages/Empty.vue')
                },
                {
                    path: '/passenger/setting',
                    name: 'CutomerSetting',
                    component: () => import('@/views/pages/passenger/PassengerSetting.vue')
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
        },
        {
            path: '/auth/chatList',
            name: 'ChatList',
            component: () => import('@/views/pages/auth/ChatList.vue')
        }
    ]
});

// Navigation guard to check authentication before each navigation
router.beforeEach(async (to, from, next) => {
    if (to.name !== 'login' && !store.getters.login) {
        try {
            console.log('check login');
            await store.dispatch('checkLogin');
            next();
        } catch {
            next({ name: 'login' });
        }
    } else {
        console.log('have login');
        next(); // Continue with the navigation
    }
});

export default router;
