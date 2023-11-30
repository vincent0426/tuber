<script setup>
import { ref } from 'vue';

const menu = ref(null);
const contextMenu = ref(null);

const menuitems = ref([
    {
        label: 'Customers',
        items: [
            {
                label: 'New',
                icon: 'pi pi-fw pi-plus'
            },
            {
                label: 'Edit',
                icon: 'pi pi-fw pi-user-edit'
            }
        ]
    },
    {
        label: 'Orders',
        items: [
            {
                label: 'View',
                icon: 'pi pi-fw pi-list'
            },
            {
                label: 'Search',
                icon: 'pi pi-fw pi-search'
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
    <div class="grid p-fluid">
        <div class="col-12 md:col-4">
            <div class="card">
                <h5>Settings</h5>
                <Menu :model="menuitems"></Menu>
            </div>
        </div>
    </div>
</template>