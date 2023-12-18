<script setup>
import { ref, onMounted } from 'vue';
import ProductService from '@/service/ProductService';

const dataviewValue = ref(null);
const layout = ref('list');
const sortKey = ref(null);
const sortOrder = ref(null);
const sortField = ref(null);
const sortOptions = ref([
    { label: 'Price High to Low', value: '!price' },
    { label: 'Price Low to High', value: 'price' }
]);

const productService = new ProductService();

onMounted(() => {
    productService.getProductsSmall().then((data) => (dataviewValue.value = data));
});


const onSortChange = (event) => {
    const value = event.value.value;
    const sortValue = event.value;

    if (value.indexOf('!') === 0) {
        sortOrder.value = -1;
        sortField.value = value.substring(1, value.length);
        sortKey.value = sortValue;
    } else {
        sortOrder.value = 1;
        sortField.value = value;
        sortKey.value = sortValue;
    }
};
</script>

<template>
    <div className="card">
        <h3>Take a Ride</h3>
        <h3> with others</h3>
      <div class="col-12 md:col-6">
        <div class="p-inputgroup">
            <!-- This is for Search Page. -->
            <router-link to="/passenger/search">
                <Button label="Search"> </Button>
            </router-link>
            <InputText placeholder="Where to?" />
        </div>
      </div>        
    </div>
    <div class="grid">
        <div class="col-12">
            <div class="card">
                <router-link to="/passenger/mytrip">
                    <h5>Your Rides</h5>
                </router-link>
                 <DataView :value="dataviewValue" :layout="layout" :paginator="true" :rows="6" :sortOrder="sortOrder" :sortField="sortField">
                 <!--   <template #header>
                        <div class="grid grid-nogutter">
                            <div class="col-6 text-left">
                                <Dropdown v-model="sortKey" :options="sortOptions" optionLabel="label" placeholder="Sort By Price" @change="onSortChange($event)" />
                            </div>
                            <div class="col-6 text-right">
                                <DataViewLayoutOptions v-model="layout" />
                            </div> 
                        </div>
                    </template>
                    <template #list="slotProps">
                        <div class="col-12">
                            <div class="flex flex-column md:flex-row align-items-center p-3 w-full">
                                img :src="'demo/images/product/' + slotProps.data.image" :alt="slotProps.data.name" class="my-4 md:my-0 w-9 md:w-10rem shadow-2 mr-5" />
                                <div class="flex-1 text-center md:text-left">
                                    <div class="font-bold text-2xl">{{ slotProps.data.name }}</div>
                                    <div class="mb-3">{{ slotProps.data.description }}</div>
                                    <Rating :modelValue="slotProps.data.rating" :readonly="true" :cancel="false" class="mb-2"></Rating>
                                    <div class="flex align-items-center">
                                        <i class="pi pi-tag mr-2"></i>
                                        <span class="font-semibold">{{ slotProps.data.category }}</span>
                                    </div>
                                </div>
                                <div class="flex flex-row md:flex-column justify-content-between w-full md:w-auto align-items-center md:align-items-end mt-5 md:mt-0">
                                    <span class="text-l font-semibold mb-2 align-self-center md:align-self-end">${{ slotProps.data.price }}</span>
                                     <Button icon="pi pi-shopping-cart" label="Add to Cart" :disabled="slotProps.data.inventoryStatus === 'OUTOFSTOCK'" class="mb-2"></Button>
                                    <span :class="'product-badge status-' + slotProps.data.inventoryStatus.toLowerCase()">{{ slotProps.data.inventoryStatus }}</span>
                                </div>
                            </div>
                        </div>
                    </template> -->
                </DataView> 
            </div>
        </div>
      </div>
</template>
