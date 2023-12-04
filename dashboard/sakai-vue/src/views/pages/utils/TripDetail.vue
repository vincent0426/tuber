<!-- The page for trip detail. Notice Customer and Driver may see different layout of this page. -->
<script setup>
import { ref, onMounted,onBeforeMount } from 'vue';
import ProductService from '@/service/ProductService';
import { FilterMatchMode, FilterOperator } from 'primevue/api';
import CustomerService from '@/service/CustomerService';
import { TripService } from '@/service/TripService';

const tripService = new TripService();
const tripID = window.location.href.split('/').filter(segment => segment.trim() !== '')[4];
const StartStaion = ref(null);
const EndStaion = ref(null);
const source = ref(null);
const destination = ref(null);
const driverName = ref(null);
const driver_image_url = ref(null);
const driver_plate = ref(null);
const start_time = ref(null);
var tripData;

onMounted(() => {
    tripService.getTrip(tripID).then((data) => {
        source.value = data.source_name;
        destination.value = data.destination_name;
        driverName.value = data.driver_name;
        driver_image_url.value = data.driver_image_url;
        driver_plate.value = data.driver_plate;
        start_time.value = data.start_time;
        console.log(data);
    });
});

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
              console.log(this);
              var place = this.getPlace();
              createMarker(place);
              map.setCenter(place.geometry.location);
              map.setZoom(17.5);
            });
            autocompletes.push(autocomplete);
          }
        });

        directionsRenderer.setMap(map);

        const onChangeHandler = function () {
          deleteMarkers();
          AddStop();
          calculateAndDisplayRoute(directionsService, directionsRenderer);
        };
        const AddStop = function () {
          const start = document.getElementById("Start").value;
          const end = document.getElementById("End").value;
          stops.push({
            location:start,
            stopover: true,
          });
          stops.push({
            location:end,
            stopover: true,
          });
          
        }
        document.getElementById("checkPath").addEventListener("click", onChangeHandler);
        tripService.getTrip(tripID).then((data) => {
            tripData = data;
            calculateAndDisplayRoute(directionsService, directionsRenderer);
        });
        const onJoinTrip = function (){
            tripService.joinTrip(tripID).then((data) => {
                console.log(data);
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
    function deleteMarkers(){
      for (let i = 0; i < markers.length; i++) {
        markers[i].setMap(null);
      }
      markers = [];
    }
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
                    <div class="font-bold text-2xl">Start Time:{{ start_time }}</div>
                </div>
            </div>
        </div>
        <!-- <div class="card col-12">
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
                    :globalFilterFields="['name', 'status']"
                >
                    
                    <template #empty> No customers found. </template>
                    <template #loading> Loading customers data. Please wait. </template>
                    <Column field="name" header="Name" style="min-width: 12rem">
                        <template #body="{ data }">
                            {{ data.name }}
                        </template>
                        <template #filter="{ filterModel }">
                            <InputText type="text" v-model="filterModel.value" class="p-column-filter" placeholder="Search by name" />
                        </template>
                    </Column>
                    <Column field="status" header="Status" :filterMenuStyle="{ width: '14rem' }" style="min-width: 12rem">
                        <template #body="{ data }">
                            <span :class="'customer-badge status-' + data.status">{{ data.status }}</span>
                        </template>
                        <template #filter="{ filterModel }">
                            <Dropdown v-model="filterModel.value" :options="statuses" placeholder="Any" class="p-column-filter" :showClear="true">
                                <template #value="slotProps">
                                    <span :class="'customer-badge status-' + slotProps.value" v-if="slotProps.value">{{ slotProps.value }}</span>
                                    <span v-else>{{ slotProps.placeholder }}</span>
                                </template>
                                <template #option="slotProps">
                                    <span :class="'customer-badge status-' + slotProps.option">{{ slotProps.option }}</span>
                                </template>
                            </Dropdown>
                        </template>
                    </Column>
                    <Column field="PickUpAt" header="Pick Up At" :showFilterMatchModes="false" style="min-width: 12rem">
                        <template #body="{ data }">
                            <ProgressBar :value="data.activity" :showValue="false" style="height: 0.5rem"></ProgressBar>
                        </template>
                        <template #filter="{ filterModel }">
                            <Slider v-model="filterModel.value" :range="true" class="m-3"></Slider>
                            <div class="flex align-items-center justify-content-between px-2">
                                <span>{{ filterModel.value ? filterModel.value[0] : 0 }}</span>
                                <span>{{ filterModel.value ? filterModel.value[1] : 100 }}</span>
                            </div>
                        </template>
                    </Column>
                    <Column field="DropAt" header="Drop At" dataType="boolean" bodyClass="text-center" style="min-width: 8rem">
                        <template #body="{ data }">
                            <i class="pi" :class="{ 'text-green-500 pi-check-circle': data.verified, 'text-pink-500 pi-times-circle': !data.verified }"></i>
                        </template>
                        <template #filter="{ filterModel }">
                            <TriStateCheckbox v-model="filterModel.value" />
                        </template>
                    </Column>
                </DataTable>
            </div> -->
        <div class="col-12 card">
            <div class="grid grid-nogutter">
                <div class="col-4 text-left">
                    <h5>Start</h5>
                    <InputText class="search-location" placeholder="Search" id="Start" v-model="StartStaion"/>
                </div>
                <div class="col-4 text-middle">
                    <h5>End</h5>
                    <InputText class="search-location" placeholder="Search" id="End" v-model="EndStaion"/>
                    <br><br>
                </div>
                <div class="col-4 text-right">
                    <Button label="CheckPath" class="mr-2 mb-2" id="checkPath"></Button>
                    <Button label="Save" class="mr-2 mb-2" id="Save"></Button>
                </div>
                
            </div>
        </div>
        
    </div>
</template>
