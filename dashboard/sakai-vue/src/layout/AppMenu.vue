<script setup>
import { ref } from 'vue';

import AppMenuItem from './AppMenuItem.vue';
import store from '../store/index.js';

const model_passenger = ref([
    {
        label: 'Uber',
        items: [
            // { label: 'Dashboard', icon: 'pi pi-fw pi-home', to: '/' },
            // { label: 'logo', icon: 'pi pi-fw pi-home', to: '/logo' },
            // { label: 'Register', icon: 'pi pi-fw pi-home', to: '/auth/register' },
            // { label: 'login', icon: 'pi pi-fw pi-home', to: '/login' },
            // { label: 'set profile', icon: 'pi pi-fw pi-home', to: '/profile' },
            { label: 'Home', icon: 'pi pi-fw pi-home', to: '/passenger/home' },
            { label: 'Search', icon: 'pi pi-fw pi-search', to: '/passenger/search' },
            { label: 'My trip', icon: 'pi pi-fw pi-car', to: '/passenger/mytrip' },
            { label: 'History trip', icon: 'pi pi-fw pi-book', to: '/passenger/history' },
            { label: 'Setting', icon: 'pi pi-fw pi-cog', to: '/setting' },
            { label: 'Favorite Driver', icon: 'pi pi-fw pi-heart', to: '/passenger/favorite' },
            { label: 'Chat', icon: 'pi pi-fw pi-comment', to: '/chatList' }
        ]
    },
    {
        label: 'Pages',
        icon: 'pi pi-fw pi-briefcase',
        to: '/pages',
        items: [
            {
                label: 'Auth',
                icon: 'pi pi-fw pi-user',
                items: [
                    {
                        label: 'Logout',
                        icon: 'pi pi-fw pi-sign-out',
                        to: '/login'
                    },
                    {
                        label: 'Become a Driver',
                        icon: 'pi pi-fw pi-lock',
                        to: '/become-driver'
                    },
                    {
                        label: 'Switch to Driver',
                        icon: 'pi pi-fw pi-arrows-h',
                        to: '/driver/home'
                    }
                ]
            }
        ]
    }
]);

const model_driver = ref([
    {
        label: 'Uber',
        items: [
            { label: 'Home', icon: 'pi pi-fw pi-home', to: '/driver/home' },
            { label: 'Create A Trip', icon: 'pi pi-fw pi-plus', to: '/driver/create' },
            { label: 'History', icon: 'pi pi-fw pi-book', to: '/driver/history' },
            { label: 'Setting', icon: 'pi pi-fw pi-cog', to: '/setting' },
            { label: 'Chat', icon: 'pi pi-fw pi-comment', to: '/chatList' }
        ]
    },
    {
        label: 'Pages',
        icon: 'pi pi-fw pi-briefcase',
        to: '/pages',
        items: [
            {
                label: 'Auth',
                icon: 'pi pi-fw pi-user',
                items: [
                    {
                        label: 'Logout',
                        icon: 'pi pi-fw pi-sign-out',
                        to: '/login'
                    },
                    // {
                    //     label: 'Access Denied',
                    //     icon: 'pi pi-fw pi-lock',
                    //     to: '/auth/access'
                    // },
                    {
                        label: 'Switch to Passenger',
                        icon: 'pi pi-fw pi-arrows-h',
                        to: '/passenger/home'
                    }
                ]
            }
        ]
    }
]);
</script>

<template>
    <ul class="layout-menu">
        <template v-if="store.getters.role === 'passenger'" v-for="(item, i) in model_passenger" :key="item">
            <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
            <li v-if="item.separator" class="menu-separator"></li>
        </template>

        <template v-if="store.getters.role === 'driver'" v-for="(item, i) in model_driver" :key="item">
            <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
            <li v-if="item.separator" class="menu-separator"></li>
        </template>
    </ul>
</template>

<style lang="scss" scoped></style>
