<script setup>
import { ref } from 'vue';

import AppMenuItem from './AppMenuItem.vue';

const model = ref([
    {
        label: 'Uber',
        items: [
            { label: 'Home', icon: 'pi pi-fw pi-home', to: '/driver/home' },
            { label: 'Create', icon: 'pi pi-fw pi-plus', to: '/driver/create'},
            { label: 'History', icon: 'pi pi-fw pi-book', to: '/driver/history'},
            { label: 'Setting', icon: 'pi pi-fw pi-cog', to: '/driver/setting'},
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
                        label: 'Switch to Customer',
                        icon: 'pi pi-fw pi-arrows-h',
                        to: '/customer/home'
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
