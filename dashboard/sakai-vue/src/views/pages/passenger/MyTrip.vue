<script setup>
import { ref, onMounted } from 'vue';
import { TripService } from '@/service/TripService';
import { LocationService } from '@/service/LocationService';
const tripService = new TripService();
const locationService = new LocationService();
const dataviewValue = ref([]);

onMounted(async () => {
  let currentPage = 1; // 設置當前頁數
  const pageSize = 10; // 設置每頁顯示的行程數量

  try {
    while (true) {
      if(currentPage == 1){
        const response = await tripService.getMyTrips();
        console.log(response);
        const { items, total } = response; // 假設後端返回的數據中有 items 和 totalPages 屬性
        dataviewValue.value = items;
        // 如果當前頁已達到總頁數，則跳出循環
        if (currentPage*pageSize >= total) {
            console.log("break");
            break;
        }

        currentPage++; // 頁數加一，繼續獲取下一頁的數據
      }
      else{
        const response = await tripService.getThePageTrip(currentPage);
        //console.log(response);
        const { items, total } = response; // 假設後端返回的數據中有 items 和 totalPages 屬性
        console.log("1");
        console.log(dataviewValue.value);
        console.log("1");
        
        console.log(items);
        for(var i = 0;i < items.length;i++){
            dataviewValue.value.push(items[i]);
        }
        // 如果當前頁已達到總頁數，則跳出循環
        if (currentPage*pageSize >= total) {
            console.log("break");
            break;
        }

        currentPage++; // 頁數加一，繼續獲取下一頁的數據
      }
    }
  } catch (error) {
    console.error('Error fetching trips:', error);
  }
});


const datetime24h = ref(null);
const StartStaion = ref(null);
const EndStaion = ref(null);
const layout = ref('list');
const sortKey = ref(null);
const sortOrder = ref(null);
const sortField = ref(null);

// const test = (data) => {
    
//     console.log(locationService.getLocation(data));
        
// };

</script>

<template>
    <div class="grid">
        <div class="col-12">
            <h3>My Trip</h3>
            <div class="card">
                <DataView :value="dataviewValue" :layout="layout" :paginator="true" :rows="10">
                   
                    <template #list="slotProps">
                        <div class="col-12">
                            <div class="flex flex-column md:flex-row align-items-center p-3 w-full">
                                <div class="flex-1 text-center md:text-left">
                                    <div class="font-bold text-2xl">From:{{ slotProps.data.SourceID }}</div>
                                    <div class="font-bold text-2xl">To:{{ slotProps.data.DestinationID }}</div>
                                    <div class="mb-3">Driver Name:{{ slotProps.data.DriverName }}</div>
                                    <div class="mb-3">Start Time:{{ slotProps.data.StartTime }}</div>
                                    <!-- <Rating :modelValue="slotProps.data.rating" :readonly="true" :cancel="false" class="mb-2"></Rating> -->
                                </div>
                                <div class="flex flex-row md:flex-column justify-content-between w-full md:w-auto align-items-center md:align-items-end mt-5 md:mt-0">
                                    <!-- <span class="text-2xl font-semibold mb-2 align-self-center md:align-self-end">${{ slotProps.data.price }}</span> -->
                                    <!-- <Button label="Apply" :disabled="slotProps.data.inventoryStatus === 'OUTOFSTOCK'" class="mb-2" onclick="location.href='/#/TripDetail/'+ {{ slotProps.data.id }}"></Button> -->
                                    <router-link :to="'/OnTrip/' + slotProps.data.TripID">
                                        <Button label="Start Trip" :disabled="slotProps.data.inventoryStatus === 'OUTOFSTOCK'" class="mb-2"></Button>
                                    </router-link>
                                </div>
                            </div>
                        </div>
                    </template>
                </DataView>
            </div>
        </div>

        
    </div>
</template>
