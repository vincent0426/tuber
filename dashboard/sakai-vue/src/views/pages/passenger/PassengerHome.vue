<script setup>
import { ref, onMounted } from 'vue';
import { TripService } from '@/service/TripService';
const tripService = new TripService();
const dataviewValue = ref([]);

const layout = ref('list');
const sortKey = ref(null);
const sortOrder = ref(null);
const sortField = ref(null);
const sortOptions = ref([
    { label: 'Price High to Low', value: '!price' },
    { label: 'Price Low to High', value: 'price' }
]);


onMounted(async () => {
  let currentPage = 1; // 設置當前頁數
  const pageSize = 10; // 設置每頁顯示的行程數量

  try {
    while (true) {
      if(currentPage == 1){
        const response = await tripService.getMyTrips(false);
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




function DateConvert(dateString) {
    const date = new Date(dateString);
    // 取得日期和時間的部分
    const year = date.getFullYear(); // 年份
    const month = `0${date.getMonth() + 1}`.slice(-2); // 月份（補0）
    const day = `0${date.getDate()}`.slice(-2); // 日（補0）
    const hours = `0${date.getHours()}`.slice(-2); // 小時（補0）
    const minutes = `0${date.getMinutes()}`.slice(-2); // 分鐘（補0）
    // 格式化成"YYYY-MM-DDTHH:MM:SSZ"的形式
    const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:00`;

    return formattedDate;
}


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
            <h3>My Trip</h3>
            <div class="card">
                <DataView :value="dataviewValue" :layout="layout" :paginator="true" :rows="10">
                   
                    <template #list="slotProps">
                        <div class="card">
                            <!-- <div class="border-round border-1 surface-border p-4"> -->
                                <div class="flex mb-3">
                                <img :src="slotProps.data.driver_image_url" :alt="slotProps.data.driver_name" class="w-3 shadow-2 my-3 mx-0" style="margin: auto;"/>
                                <div style="margin: auto;">
                                    <div class="mb-1">Driver:  {{ slotProps.data.driver_name }}</div>
                                    <div class="mb-2">Time:<br>{{ DateConvert(slotProps.data.start_time) }}</div>
                                </div>
                                </div>
                                <div class="flex-1 text-left md:text-left" width="100%" height="150%">
                                    <Tag class="mr-2" value="From" :rounded="true"></Tag>
                                    <div class="font-bold text-xl">{{ slotProps.data.source_name }}</div>
                                    <Tag class="mr-2" value="To" :rounded="true"></Tag>
                                    <div class="font-bold text-xl">{{ slotProps.data.destination_name }}</div>
                                   
                                </div>
                                <div class="flex justify-content-between mt-3">
                                    <router-link :to="'/TripDetail/' + slotProps.data.id">
                                        <Button label="Apply" :disabled="slotProps.data.inventoryStatus === 'OUTOFSTOCK'" class="mb-2"></Button>
                                    </router-link>
                                </div>
                            <!-- </div> -->
                        </div>
                    </template>
                </DataView>
            </div>
        </div>
    </div>
</template>
