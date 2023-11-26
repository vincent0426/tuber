<script setup>
import { ref } from 'vue';

import AppMenuItem from './AppMenuItem.vue';

const model = ref([
    {
        label: 'Uber',
        items: [
            { label: 'Home', icon: 'pi pi-fw pi-home', to: '/customer/home' },
            { label: 'Search', icon: 'pi pi-fw pi-search', to: '/customer/search'},
            { label: 'My trip', icon: 'pi pi-fw pi-car', to: '/customer/mytrip'},
            { label: 'History', icon: 'pi pi-fw pi-book', to: '/customer/history'},
            { label: 'Setting', icon: 'pi pi-fw pi-cog', to: '/customer/setting'},
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
                        to: '/auth/login'
                    },
                    {
                        label: 'Access Denied',
                        icon: 'pi pi-fw pi-lock',
                        to: '/auth/access'
                    },
                    {
                        label: 'Switch to Driver',
                        icon: 'pi pi-fw pi-arrows-h',
                        to: '/driver/home'
                    }
                ]
            },
        ]
    },
]);
</script>

<template>
    <ul class="layout-menu">
        <template v-for="(item, i) in model" :key="item">
            <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
            <li v-if="item.separator" class="menu-separator"></li>
        </template>
    </ul>
</template>

<style lang="scss" scoped></style>
