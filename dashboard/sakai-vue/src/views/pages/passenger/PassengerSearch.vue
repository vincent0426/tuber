<script setup>
import { ref, onMounted } from 'vue';
import { TripService } from '@/service/TripService';

const tripService = new TripService();
const dataviewValue = ref([]);
// onMounted(() => {
//     tripService.getAllTrip().then((data) => {
//         //console.log(data); // 在控制台輸出數據的形式
//         dataviewValue.value = data.items; // 將數據賦值給 dataviewValue
        
//     });
// });
onMounted(async () => {
  let currentPage = 1; // 設置當前頁數
  const pageSize = 10; // 設置每頁顯示的行程數量

  try {
    while (true) {
      if(currentPage == 1){
        const response = await tripService.getAllTrip();
        //console.log(response);
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

// (g=>{var h,a,k,p="The Google Maps JavaScript API",c="google",l="importLibrary",q="__ib__",m=document,b=window;b=b[c]||(b[c]={});var d=b.maps||(b.maps={}),r=new Set,e=new URLSearchParams,u=()=>h||(h=new Promise(async(f,n)=>{await (a=m.createElement("script"));e.set("libraries",[...r]+"");for(k in g)e.set(k.replace(/[A-Z]/g,t=>"_"+t[0].toLowerCase()),g[k]);e.set("callback",c+".maps."+q);a.src=`https://maps.${c}apis.com/maps/api/js?`+e;d[q]=f;a.onerror=()=>h=n(Error(p+" could not load."));a.nonce=m.querySelector("script[nonce]")?.nonce||"";m.head.append(a)}));d[l]?console.warn(p+" only loads once. Ignoring:",g):d[l]=(f,...n)=>r.add(f)&&u().then(()=>d[l](f,...n))})({
//       key: "AIzaSyCWk9OsA3BidynIgg5_ybz2dWVIBkuWpxE",
//       v: "weekly",
//       // Use the 'v' parameter to indicate the version to use (weekly, beta, alpha, etc.).
//       // Add other bootstrap parameters as needed, using camel case.
//     });
    
//     let currentPosition;
//     var searchInputs = document.getElementsByClassName("search-location");

// async function init() {
//         // Request libraries when needed, not in the script tag.
//         const { Map } = await google.maps.importLibrary("maps");
//         const { Geometry } = await google.maps.importLibrary("geometry");
//         const { Place } = await google.maps.importLibrary("places");
       
        
//         var autocompletes = [];
//         navigator.geolocation.getCurrentPosition(function(position){
//           currentPosition = {
//             lat: position.coords.latitude,
//             lng: position.coords.longitude
//           };

//           var autocompletes = [];
//           var options = {
//               bounds: {
//                 east: currentPosition.lng + 0.001 ,
//                 west: currentPosition.lng - 0.001,
//                 south: currentPosition.lat - 0.001,
//                 north: currentPosition.lat + 0.001
//               },
//               strictBounds:false,
//               types: ['establishment'],
//               componentRestrictions: { country: "tw" },
//               fields: ["formatted_address", "geometry", "icon", "name","address_components"],
//           };
//           for (var i = 0; i < searchInputs.length; i++) {
//             var autocomplete = new google.maps.places.Autocomplete(searchInputs[i], options);
//             autocompletes.push(autocomplete);
//           }
//         });
       
//     }
// init();
console.log(dataviewValue);
</script>

<template>
    <div class="grid">
        <div class="col-12">
            <h3>Search Trip</h3>
            <div class="card">
                <DataView :value="dataviewValue" :layout="layout" :paginator="true" :rows="10">
                    <!-- <template #header>
                        <div class="grid grid-nogutter">
                            <div class="col-6 text-left">
                                <h5>Start</h5>
                                <InputText class="search-location" placeholder="Search" id="Start" type="text" v-model="floatValue" />
                            </div>
                            <div class="col-6 text-middle">
                                <h5>End</h5>
                                <InputText class="search-location" placeholder="Search" id="End" type="text" v-model="floatValue" />
                                <br><br>
                            </div>
                            <div class="col-6 text-left">
                                <Calendar id="departuretime" :showIcon="true" :showButtonBar="true" v-model="datetime24h" showTime hourFormat="24"></Calendar>
                            </div>
                            <div class="col-6 text-right">
                                <Button label="Save" class="mr-2 mb-2"></Button>
                            </div>
                            
                        </div>
                    </template> -->
                    <template #list="slotProps">
                        <div class="col-12">
                            <div class="flex flex-column md:flex-row align-items-center p-3 w-full">
                                <img :src="slotProps.data.driver_image_url" :alt="slotProps.data.driver_name" class="w-3 shadow-2 my-3 mx-0" />
                                <div class="flex-1 text-center md:text-left">
                                    <div class="font-bold text-2xl">From:{{ slotProps.data.source_name }}</div>
                                    <div class="font-bold text-2xl">To:{{ slotProps.data.destination_name }}</div>
                                    <div class="mb-3">Start Time:{{ slotProps.data.start_time }}</div>
                                    <!-- <Rating :modelValue="slotProps.data.rating" :readonly="true" :cancel="false" class="mb-2"></Rating> -->
                                </div>
                                <div class="flex flex-row md:flex-column justify-content-between w-full md:w-auto align-items-center md:align-items-end mt-5 md:mt-0">
                                    <!-- <span class="text-2xl font-semibold mb-2 align-self-center md:align-self-end">${{ slotProps.data.price }}</span> -->
                                    <!-- <Button label="Apply" :disabled="slotProps.data.inventoryStatus === 'OUTOFSTOCK'" class="mb-2" onclick="location.href='/#/TripDetail/'+ {{ slotProps.data.id }}"></Button> -->
                                    <router-link :to="'/TripDetail/' + slotProps.data.id">
                                        <Button label="Apply" :disabled="slotProps.data.inventoryStatus === 'OUTOFSTOCK'" class="mb-2"></Button>
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
