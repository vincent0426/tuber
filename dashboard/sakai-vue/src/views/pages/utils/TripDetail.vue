<!-- The page for trip detail. Notice Customer and Driver may see different layout of this page. -->
<script setup>
import { ref, onMounted,onBeforeMount } from 'vue';

import { TripService } from '@/service/TripService';
import Dropdown from 'primevue/dropdown';

const tripService = new TripService();
const tripID = window.location.href.split('/TripDetail/').filter(segment => segment.trim() !== '')[1];
const StartStaion = ref(null);
const EndStaion = ref(null);
const source = ref(null);
const destination = ref(null);
const driverName = ref(null);
const driver_image_url = ref(null);
const driver_plate = ref(null);
const start_time = ref(null);
const customer1 = ref(null);
const mid_start = ref(null);
const mid_end = ref(null);
const passenger_limit = ref(null);

var tripData;

onMounted(() => {
    tripService.getTrip(tripID).then((data) => {
        //console.log(data);
        source.value = data.source_name;
        destination.value = data.destination_name;
        driverName.value = data.driver_name;
        driver_image_url.value = data.driver_image_url;
        driver_plate.value = data.driver_plate;
        start_time.value = data.start_time;
        mid_start.value = data.mid;
        mid_start.value.push({
            Lat: data.destination_latitude,
            Lon: data.destination_longitude,
            Name: data.destination_name,
            ID: data.destination_id
        });
        mid_start.value.push({
            Lat: data.source_latitude,
            Lon: data.source_longitude,
            Name: data.source_name,
            ID: data.source_id
        });
        mid_end.value = mid_start.value;
        console.log(mid_end);
    });
    tripService.getPassengers(tripID).then((data) => {
        customer1.value = data.passenger_details;
        //console.log(data);
    });
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
    const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:00 `;
    //console.log(formattedDate)
    return formattedDate;
}
(g=>{var h,a,k,p="The Google Maps JavaScript API",c="google",l="importLibrary",q="__ib__",m=document,b=window;b=b[c]||(b[c]={});var d=b.maps||(b.maps={}),r=new Set,e=new URLSearchParams,u=()=>h||(h=new Promise(async(f,n)=>{await (a=m.createElement("script"));e.set("libraries",[...r]+"");for(k in g)e.set(k.replace(/[A-Z]/g,t=>"_"+t[0].toLowerCase()),g[k]);e.set("callback",c+".maps."+q);a.src=`https://maps.${c}apis.com/maps/api/js?`+e;d[q]=f;a.onerror=()=>h=n(Error(p+" could not load."));a.nonce=m.querySelector("script[nonce]")?.nonce||"";m.head.append(a)}));d[l]?console.warn(p+" only loads once. Ignoring:",g):d[l]=(f,...n)=>r.add(f)&&u().then(()=>d[l](f,...n))})({
      key: "AIzaSyCWk9OsA3BidynIgg5_ybz2dWVIBkuWpxE",
      v: "weekly",
      // Use the 'v' parameter to indicate the version to use (weekly, beta, alpha, etc.).
      // Add other bootstrap parameters as needed, using camel case.
    });
    
    let map;
    let markerBegin;
    let markerEnd;
    let currentPosition;
    let infowindow;
    let stops = [];
    let markers = [];
    var searchInputs = document.getElementsByClassName("search-location");
    var searchpts;

    async function initMap() {
        
        // Request libraries when needed, not in the script tag.
        const { Map } = await google.maps.importLibrary("maps");
        const { Geometry } = await google.maps.importLibrary("geometry");
        const { Place } = await google.maps.importLibrary("places");
        // Short namespaces can be used.
        map = new Map(document.getElementById("map"), {
            center: { lat: 23.556, lng: 121.0122 },
            zoom: 7,
        });
        
        const directionsService = new google.maps.DirectionsService();
        const directionsRenderer = new google.maps.DirectionsRenderer();
        infowindow = new google.maps.InfoWindow();
        navigator.geolocation.getCurrentPosition(function(position){
          currentPosition = {
            lat: position.coords.latitude,
            lng: position.coords.longitude
          };
          var autocompletes = [];
          var options = {
              bounds: {
                east: currentPosition.lng + 0.001 ,
                west: currentPosition.lng - 0.001,
                south: currentPosition.lat - 0.001,
                north: currentPosition.lat + 0.001
              },
              strictBounds:false,
              types: ['establishment'],
              componentRestrictions: { country: "tw" },
              fields: ["formatted_address", "geometry", "place_id", "name"],
          };
          for (var i = 0; i < searchInputs.length; i++) {
            var autocomplete = new google.maps.places.Autocomplete(searchInputs[i], options);
            autocomplete.addListener('place_changed', function(){
              //console.log(this);
              var place = this.getPlace();
              createMarker(place);
              map.setCenter(place.geometry.location);
              map.setZoom(17.5);
            });
            autocompletes.push(autocomplete);
          }
        });
        directionsRenderer.setMap(map);
        tripService.getTrip(tripID).then((data) => {
            tripData = data;
            let midStops = data.mid;
            for(var i = 0;i < midStops.length;i++){
                stops.push({
                    location:midStops[i].Name,
                    stopover: true,
                });
            }
            calculateAndDisplayRoute(directionsService, directionsRenderer);
        });
        const onJoinTrip = function (){
            let sourceID = StartStaion.value;
            let destinationID = EndStaion.value;
            console.log(tripID,sourceID,destinationID);
            tripService.joinTrip(tripID,sourceID,destinationID)
                .then(response => {
                    // 處理成功回傳的資料
                    alert("Success");
                    console.log(response);
                })
                .catch(error => {
                    // 處理錯誤
                    console.log(error);
                    alert(error.error);
                });
        };
        document.getElementById("Save").addEventListener("click",onJoinTrip);
    }

    // Set markers at the location of each place result
    function createMarker(place) {
      if (!place.geometry || !place.geometry.location) return;

      const marker = new google.maps.Marker({
        map: map,
        position: place.geometry.location,
      });

      google.maps.event.addListener(marker, "click", () => {
        infowindow.setContent(`<div><textarea id="stop">${place.name}</textarea></div>`);
        infowindow.open(map,marker);
        //console.log(place);
      });
      markers.push(marker);
    }
    const getStartStaionValue = () => {
        console.log("Selected Start Station:", StartStaion.value);
        
    };
    // print route on map & mark recommand stop points
    function calculateAndDisplayRoute(directionsService, directionsRenderer) {
      console.log(stops);
      directionsService
        .route({
          origin: {
            query: tripData.source_name,
          },
          destination: {
            query: tripData.destination_name,
          },
          waypoints:stops,
          travelMode: google.maps.TravelMode.DRIVING,
        })
        .then((response) => {
          directionsRenderer.setDirections(response);
        })
        .catch((e) => console.log("Directions request failed due to " + e));
    }
    initMap();
</script>


<template>
    <div class="grid">
        <div class="col-12">
            <h3>Trip Info</h3>
            <div class="card">
                <div class="grid grid-nogutter">
                    <div class="col-4 text-left">
                        <img :src="driver_image_url" :alt="driverName" class="w-3 shadow-2 my-3 mx-0" />
                    </div>
                    <div class="col-18 text-middle">
                        <div class="font-bold text-2xl">Driver:{{ driverName }}</div>
                        <div class="font-bold text-2xl">License plate number:{{ driver_plate }}</div>
                    </div>
                
                </div>

                
                <br>
            </div>
        </div>
        <div class="col-12">
            <div class="card">
                <h4>Trip</h4>
                <div class="flex justify-content-left">
                    <div id="map" style="width: 100%; height: 50vh"></div>
                </div>
                <div class="col-16 text-left">
                    <div class="font-bold text-2xl">From:{{ source }}</div>
                    <div class="font-bold text-2xl">To:{{ destination }}</div>
                    <div class="font-bold text-2xl">Start Time:{{ DateConvert(start_time) }}</div>
                    
                </div>
            </div>
        </div>
        <div class="card col-12">
                <h5>Other Passengers</h5>
                <DataTable
                    :value="customer1"
                    :paginator="true"
                    class="p-datatable-gridlines"
                    :rows="10"
                    dataKey="id"
                    :rowHover="true"
                    v-model:filters="filters1"
                    filterDisplay="menu"
                    :loading="loading1"
                    :filters="filters1"
                    responsiveLayout="scroll"
                    :globalFilterFields="['passenger_id', 'status']"
                >
                    
                    <template #empty> No customers found. </template>
                    <template #loading> Loading customers data. Please wait. </template>
                    <Column field="passenger_id" header="ID" style="min-width: 12rem">
                        <template #body="{ data }">
                            {{ data.passenger_name }}
                        </template>
                    </Column>
                    <Column field="source_name" header="Source" style="min-width: 12rem">
                        <template #body="{ data }">
                            {{ data.source_name }}
                        </template>
                    </Column>
                    <Column field="destination_name" header="Destination" style="min-width: 12rem">
                        <template #body="{ data }">
                            {{ data.destination_name }}
                        </template>
                    </Column>
                    
                </DataTable>
            </div> 
        <div class="col-12 card">
            <div class="grid grid-nogutter">
                
                <div class="col-12">
                    <h5>Start</h5>
                    <Dropdown id="StartStop" v-model="StartStaion" :options="mid_start" optionValue="ID" optionLabel="Name" placeholder="Select a Start" @change="getStartStaionValue" class="w-full md:w-14rem">
                    </Dropdown>
                </div>
                <div class="col-12">
                    <h5>End</h5>
                    <Dropdown id="EndStop" v-model="EndStaion" :options="mid_end" optionValue="ID" optionLabel="Name" placeholder="Select a End" class="w-full md:w-14rem">
                    </Dropdown>
                </div>
                <br><br>
                <div class="col-4">
                    <Button label="Save" class="mr-2 mb-2" id="Save"></Button>
                </div>
                
            </div>
        </div>
        
    </div>
</template>
