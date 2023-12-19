<script setup>
import { TripService } from '@/service/TripService';

import { useStore } from 'vuex';
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const store = useStore();
const user = store.getters.user;
const route = useRoute();
const socket = ref(null);
const messages = ref([]);
const newMessage = ref('');

const tripService = new TripService();

const tripID = window.location.href.split('/OnTrip/').filter(segment => segment.trim() !== '')[1];
const source = ref(null);
const destination = ref(null);
const driverName = ref(null);
const driver_image_url = ref(null);
const driver_plate = ref(null);
const start_time = ref(null);
const passengerLimit = ref(null);
const mid = ref(null);
const customer1 =ref(null);

const myData = localStorage.getItem('vuex-state');
const parsedData = JSON.parse(myData);
const role = parsedData.role;
var tripData;
var location;
var isOnTrip;
const sendMessage = (currentPosition) => {
    socket.value.send(JSON.stringify({ "latitute": currentPosition.lat, "longitude": currentPosition.lng }));
    console.log(currentPosition);
};

const processMessage = (rawMessage) => {
    const parsedMessage = JSON.parse(rawMessage);
    
    currentPosition = {
        lat: parsedMessage.latitute,
        lng: parsedMessage.longitude
    };
    if(markerNow !== ''){
        markerNow.setMap(null);
    }
    markerNow = new google.maps.Marker({
        map: map,
        position: {lat: currentPosition.lat, lng: currentPosition.lng},
    });
    map.setCenter(currentPosition);
    map.setZoom(20);
};
function initializeWebSocket(){
    if(role == 'driver'){
        socket.value = new WebSocket(`ws://localhost:3003/v1/ws/driver?trip_id=${tripID}`);
        console.log("driver");
    }
    else{
        socket.value = new WebSocket(`ws://localhost:3003/v1/ws/passenger?trip_id=${tripID}`);
        console.log("pass");
    }
    

    socket.value.addEventListener('open', (event) => {
        console.log('WebSocket is open now.', event);
    });

    socket.value.addEventListener('message', (event) => {
        console.log('Message from server:', event.data);
        if(role == 'passenger'){
            processMessage(event.data);
        }
    });

    socket.value.addEventListener('error', (event) => {
        console.error('WebSocket error observed:', event);
    });

    socket.value.addEventListener('close', (event) => {
        console.log('WebSocket is closed now.', event);
    });
}

onMounted(() => {
    initializeWebSocket();
    tripService.getTrip(tripID).then((data) => {
        source.value = data.source_name;
        destination.value = data.destination_name;
        driverName.value = data.driver_name;
        driver_image_url.value = data.driver_image_url;
        driver_plate.value = data.driver_plate;
        start_time.value = data.start_time;
        mid.value = data.mid;
        passengerLimit.value = data.passenger_limit;
        console.log(data);
    });
    tripService.getPassengers(tripID).then((data) => {
        let arr = data.passenger_details;
        let arrlen = arr.length;
        for(var i = 0;i < arr.length;i++){
            if(arr[i].passenger_status == 'pending'){
                delete arr[i];
                arrlen -= 1;
            }
        }
        console.log(arr);
        arr.length = arrlen;
        customer1.value = arr;
    });
});
onBeforeUnmount(() => {
    if (socket.value) {
        socket.value.close();
    }
});

(g=>{var h,a,k,p="The Google Maps JavaScript API",c="google",l="importLibrary",q="__ib__",m=document,b=window;b=b[c]||(b[c]={});var d=b.maps||(b.maps={}),r=new Set,e=new URLSearchParams,u=()=>h||(h=new Promise(async(f,n)=>{await (a=m.createElement("script"));e.set("libraries",[...r]+"");for(k in g)e.set(k.replace(/[A-Z]/g,t=>"_"+t[0].toLowerCase()),g[k]);e.set("callback",c+".maps."+q);a.src=`https://maps.${c}apis.com/maps/api/js?`+e;d[q]=f;a.onerror=()=>h=n(Error(p+" could not load."));a.nonce=m.querySelector("script[nonce]")?.nonce||"";m.head.append(a)}));d[l]?console.warn(p+" only loads once. Ignoring:",g):d[l]=(f,...n)=>r.add(f)&&u().then(()=>d[l](f,...n))})({
      key: "AIzaSyCWk9OsA3BidynIgg5_ybz2dWVIBkuWpxE",
      v: "weekly",
      // Use the 'v' parameter to indicate the version to use (weekly, beta, alpha, etc.).
      // Add other bootstrap parameters as needed, using camel case.
    });
    
    let map;
    let currentPosition;
    let infowindow;
    var searchInputs = document.getElementsByClassName("search-location");
    let stops = [];
    let markers = [];
    let markerNow;

    async function initMap() {
        isOnTrip = true;
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
        navigator.geolocation.getCurrentPosition(function(position){
          currentPosition = {
            lat: position.coords.latitude,
            lng: position.coords.longitude
          };
        //   map.setCenter(currentPosition);
        //   map.setZoom(17.5);
          markerNow = new google.maps.Marker({
                map: map,
                position: {lat: currentPosition.lat, lng: currentPosition.lng},
          });
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
              fields: ["formatted_address", "geometry", "icon", "name","address_components"],
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
            refreshCurrentPlace();
        });
        document.getElementById('endTrip').addEventListener('click', EndTrip);
    }
   
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
    function refreshCurrentPlace(){
        console.log(isOnTrip);
        if(role == 'driver' && isOnTrip){
            navigator.geolocation.getCurrentPosition(function(position){
                currentPosition = {
                    lat: position.coords.latitude,
                    lng: position.coords.longitude
                };
                if(markerNow !== ''){
                    markerNow.setMap(null);
                    console.log("reset");
                }
                console.log("MakerNow:",markerNow)
                markerNow = new google.maps.Marker({
                    map: map,
                    position: {lat: currentPosition.lat, lng: currentPosition.lng},
                });
                map.setCenter(currentPosition);
                map.setZoom(20);
                sendMessage(currentPosition);
            });
        }
    }
    const EndTrip = function(){
        tripService.endTrip(tripID,Number(passengerLimit)).then((data) => {
            console.log(data);
            isOnTrip = false;
            alert("End");
        });
    }
    initMap();
    setInterval(function() {refreshCurrentPlace()}, 15000);

</script>
<template>
    <div class="grid">
        <div class="col-12 card">
                <h4>Trip</h4>
                <router-link :to="'/'+ role +'/home'">
                    <Button label="Home" class="mb-2"></Button>
                </router-link>
                <div class="flex justify-content-left">
                    <div id="map" style="width: 100%; height: 70vh"></div>
                </div>
                <div class="col-16 text-left">
                    <br>
                    <div class="font-bold text-2xl">From:{{ source }}</div>
                    <br>
                    <div class="font-bold text-2xl">To:{{ destination }}</div>
                    <br>
                    <div class="font-bold text-2xl">Start Time:{{ start_time }}</div>
                    <br>
                    <Button class="font-bold text-cenetr" id="endTrip" label="End" v-if="role=='driver'"/>
                    
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
                    <Column field="passenger_id" header="ID" style="min-width: 6rem">
                        <template #body="{ data }">
                            {{ data.passenger_name }}
                        </template>
                    </Column>
                    <Column field="source_name" header="Source" style="min-width: 6rem">
                        <template #body="{ data }">
                            {{ data.source_name }}
                        </template>
                    </Column>
                    <Column field="destination_name" header="Destination" style="min-width: 6rem">
                        <template #body="{ data }">
                            {{ data.destination_name }}
                        </template>
                    </Column>
                    <Column field="status" header="Status" style="min-width: 6rem" v-if="isDriver">
                        <template #body="{ data }">
                            {{ data.passenger_status }}
                            
                        </template>
                        
                    </Column>
                    <Column field="Add" header="Add" style="min-width: 6rem" v-if="isDriver">
                        <template #body="{ data }">
                            <Button label="+" class="mr-2 mb-2"></Button>
                        </template>
                    </Column>
                    
                </DataTable>
            </div> 
    </div>
</template>