<script setup>
import { ref, onMounted } from 'vue';
import { TripService } from '@/service/TripService';


const tripService = new TripService();
const dataviewValue = ref([]);

onMounted(async () => {
  let currentPage = 1; // 設置當前頁數
  const pageSize = 10; // 設置每頁顯示的行程數量

  try {
    while (true) {
      if(currentPage == 1){
        const response = await tripService.getMyTrips(true);
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


const layout = ref('list');

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
</script>

<template>
    <div className="card mt-4">
        <h3>Drive Home</h3>
        <h3>with extra Revenue</h3>
        <div class="col-12 md:col-6">
            <Button label="Create a trip">
                <router-link to="/driver/create" class="text-white">Create a trip</router-link>
            </Button>
        </div>
    </div>
    <div class="grid">
        <div class="col-12">
            <h3>My Trip</h3>
            <div class="card">
                <DataView :value="dataviewValue" :layout="layout" :paginator="true" :rows="10">
                   
                    <template #list="slotProps">
                        <div class="col-12">
                            <div class="flex flex-column md:flex-row align-items-center p-3 w-full">
                                <div class="flex-1 md:text-left">
                                    <div class="font-bold text-2xl">From:{{ slotProps.data.SourceName }}</div>
                                    <div class="font-bold text-2xl">To:{{ slotProps.data.DestinationName }}</div>
                                    <div class="mb-3">Start Time:{{ DateConvert(slotProps.data.StartTime) }}</div>
                                </div>
                                <div class="text-center">
                                    <router-link :to="'/OnTrip/' + slotProps.data.TripID">
                                        <Button label="Start Trip" class="mb-2"></Button>
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
