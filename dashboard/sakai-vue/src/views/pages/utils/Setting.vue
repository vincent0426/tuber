<script setup>
import { ref } from 'vue';

const menu = ref(null);
const contextMenu = ref(null);

const menuitems = ref([
    {
        label: 'Account',
        items: [
            {
                label: 'Notifications',
                icon: 'pi pi-fw pi-bell'
            },
            {
                label: 'Edit profile',
                icon: 'pi pi-fw pi-user-edit'
                // to: '',
            },
            {
                label: 'Privacy',
                icon: 'pi pi-fw pi-lock'
            }
        ]
    },
    {
        label: 'Support & About',
        items: [
            {
                label: 'Help & Support',
                icon: 'pi pi-fw pi-question'
            },
            {
                label: 'Terms and Policies',
                icon: 'pi pi-fw pi-exclamation-circle'
            }
        ]
    },
    {
        label: 'Actions',
        items: [
            {
                label: 'Report a problem',
                icon: 'pi pi-fw pi-flag'
            },
            {
                label: 'Register new account',
                icon: 'pi pi-fw pi-user-plus',
                to: '/auth/register'
            },
            {
                label: 'Log out',
                icon: 'pi pi-fw pi-sign-out',
                to: '/login'
            }
        ]
    }
]);

const toggleMenu = (event) => {
    menu.value.toggle(event);
};

const onContextRightClick = (event) => {
    contextMenu.value.show(event);
};
</script>

<template>
    <!-- <div class="grid p-fluid"> -->
    <div class="col-12 md:col-8 lg:col-6 xl:col-4">
        <div class="card" style="height: 100%; scale: 1.15">
            <h5>Settings</h5>
            <Menu :model="menuitems" style="width: 100%; height: 100%"></Menu>
        </div>
    </div>
    <!-- </div> -->
</template>
