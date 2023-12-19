<script setup>
import { ref, onMounted } from 'vue';
import { TripService } from '@/service/TripService.js';
const autoValue = ref(null);
const selectedAutoValue = ref(null);
const autoFilteredValue = ref([]);
const datetime24h1 = ref(null);
const datetime24h2 = ref(null);
const licenseplatenumberValue = ref(null);
const passengernumberValue = ref(null);

((g) => {
    var h,
        a,
        k,
        p = 'The Google Maps JavaScript API',
        c = 'google',
        l = 'importLibrary',
        q = '__ib__',
        m = document,
        b = window;
    b = b[c] || (b[c] = {});
    var d = b.maps || (b.maps = {}),
        r = new Set(),
        e = new URLSearchParams(),
        u = () =>
            h ||
            (h = new Promise(async (f, n) => {
                await (a = m.createElement('script'));
                e.set('libraries', [...r] + '');
                for (k in g)
                    e.set(
                        k.replace(/[A-Z]/g, (t) => '_' + t[0].toLowerCase()),
                        g[k]
                    );
                e.set('callback', c + '.maps.' + q);
                a.src = `https://maps.${c}apis.com/maps/api/js?` + e;
                d[q] = f;
                a.onerror = () => (h = n(Error(p + ' could not load.')));
                a.nonce = m.querySelector('script[nonce]')?.nonce || '';
                m.head.append(a);
            }));
    d[l] ? console.warn(p + ' only loads once. Ignoring:', g) : (d[l] = (f, ...n) => r.add(f) && u().then(() => d[l](f, ...n)));
})({
    key: 'AIzaSyAe-UQfeWgzNct7R2lKiZ1iAmjFBhi2qPA',
    v: 'weekly'
    // Use the 'v' parameter to indicate the version to use (weekly, beta, alpha, etc.).
    // Add other bootstrap parameters as needed, using camel case.
});

var map;
var markerBegin;
var markerEnd;
var currentPosition;
var infowindow;
var stops;
var markers;
var searchInputs = document.getElementsByClassName('search-location');
var searchpts;
var source;
var destination;
var tempStop;
var mid = [];
async function initMap() {
    // Request libraries when needed, not in the script tag.
    const { Map } = await google.maps.importLibrary('maps');
    const { Geometry } = await google.maps.importLibrary('geometry');
    const { Place } = await google.maps.importLibrary('places');
    // Short namespaces can be used.
    map = new Map(document.getElementById('map'), {
        center: { lat: 23.556, lng: 121.0122 },
        zoom: 7
    });

    const directionsService = new google.maps.DirectionsService();
    const directionsRenderer = new google.maps.DirectionsRenderer();
    infowindow = new google.maps.InfoWindow();
    stops = [];
    markers = [];
    navigator.geolocation.getCurrentPosition(function (position) {
        currentPosition = {
            lat: position.coords.latitude,
            lng: position.coords.longitude
        };
        map.setCenter(currentPosition);
        map.setZoom(17.5);

        var autocompletes = [];
        var options = {
            bounds: {
                east: currentPosition.lng + 0.001,
                west: currentPosition.lng - 0.001,
                south: currentPosition.lat - 0.001,
                north: currentPosition.lat + 0.001
            },
            strictBounds: false,
            types: ['establishment'],
            componentRestrictions: { country: 'tw' },
            fields: ['formatted_address', 'geometry', 'place_id', 'name']
        };
        for (var i = 0; i < searchInputs.length; i++) {
            var autocomplete = new google.maps.places.Autocomplete(searchInputs[i], options);
            if (i == 0) {
                autocomplete.addListener('place_changed', function () {
                    var place = this.getPlace();
                    createMarker(place);
                    map.setCenter(place.geometry.location);
                    map.setZoom(17.5);
                    source = {
                        lat: place.geometry.location.lat(),
                        lon: place.geometry.location.lng(),
                        name: document.getElementById('Start').value,
                        place_id: place.place_id
                    };
                });
            } else if (i == 1) {
                autocomplete.addListener('place_changed', function () {
                    var place = this.getPlace();
                    createMarker(place);
                    map.setCenter(place.geometry.location);
                    map.setZoom(17.5);
                    destination = {
                        lat: place.geometry.location.lat(),
                        lon: place.geometry.location.lng(),
                        name: document.getElementById('End').value,
                        place_id: place.place_id
                    };
                });
            } else {
                autocomplete.addListener('place_changed', function () {
                    var place = this.getPlace();
                    createMarker(place);
                    map.setCenter(place.geometry.location);
                    map.setZoom(17.5);
                    tempStop = {
                        lat: place.geometry.location.lat(),
                        lon: place.geometry.location.lng(),
                        name: document.getElementById('Stop').value,
                        place_id: place.place_id
                    };
                    console.log(this);
                    console.log(place);
                });
            }

            autocompletes.push(autocomplete);
        }
    });

    directionsRenderer.setMap(map);
    const onChangeHandler = function () {
        deleteMarkers();
        calculateAndDisplayRoute(directionsService, directionsRenderer, stops);
    };
    const AddStop = function () {
        console.log('AddStop');
        var StopContent = document.getElementById('Stops').value;
        console.log('StopContent', StopContent);
        const stopName = document.getElementById('Stop').value;
        console.log('stopName', stopName);
        StopContent = StopContent + stopName + '\r\n';
        document.getElementById('Stops').value = StopContent;
        stops.push({
            location: stopName,
            stopover: true
        });
        document.getElementById('Stop').value = '';
        mid.push(tempStop);
        console.log(mid);
        
        tempStop = {};
    };
    const SaveTrip = function () {
        // 在其他地方使用你的 Service
        const tripService = new TripService();
        const newTripData = tripData();
        console.log(newTripData);
        tripService
            .createTrip(newTripData)
            .then((response) => {
                // 處理成功回傳的資料
                console.log(response);
                alert("Success");
                //window.location.href = '/driver/trip';
            })
            .catch((error) => {
                // 處理錯誤
                alert(error);
            });
    };
    document.getElementById('checkPath').addEventListener('click', onChangeHandler);
    document.getElementById('addStop').addEventListener('click', AddStop);
    document.getElementById('Save').addEventListener('click', SaveTrip);
}
// decode route's Poltline
function decodePolyline(encoded) {
    if (!encoded) {
        return [];
    }
    var poly = [];
    var index = 0,
        len = encoded.length;
    var lat = 0,
        lng = 0;

    while (index < len) {
        var b,
            shift = 0,
            result = 0;

        do {
            b = encoded.charCodeAt(index++) - 63;
            result = result | ((b & 0x1f) << shift);
            shift += 5;
        } while (b >= 0x20);

        var dlat = (result & 1) != 0 ? ~(result >> 1) : result >> 1;
        lat += dlat;

        shift = 0;
        result = 0;

        do {
            b = encoded.charCodeAt(index++) - 63;
            result = result | ((b & 0x1f) << shift);
            shift += 5;
        } while (b >= 0x20);

        var dlng = (result & 1) != 0 ? ~(result >> 1) : result >> 1;
        lng += dlng;

        var p = {
            latitude: lat / 1e5,
            longitude: lng / 1e5
        };
        poly.push(p);
    }
    return poly;
}

// Perform a Places Nearby Search Request
function getNearbyPlaces(position) {
    let request = {
        location: new google.maps.LatLng(position.latitude, position.longitude),
        radius: '10', //search radius(m)
        types: ['point_of_interest'], // set nearby search Objs' type
        fields: ['formatted_address', 'name', 'geometry']
    };
    let service = new google.maps.places.PlacesService(map);
    service.nearbySearch(request, nearbyCallback);
}

// Handle the results (up to 20) of the Nearby Search
function nearbyCallback(results, status) {
    if (status == google.maps.places.PlacesServiceStatus.OK) {
        for (var i = 0; i < results.length; i++) {
            createMarker(results[i]);
        }
    }
}

// Set markers at the location of each place result
function createMarker(place) {
    if (!place.geometry || !place.geometry.location) return;

    const marker = new google.maps.Marker({
        map: map,
        position: place.geometry.location
    });

    google.maps.event.addListener(marker, 'click', () => {
        infowindow.setContent(`<div><textarea id="stop">${place.name}</textarea></div>`);
        infowindow.open(map, marker);
    });
    markers.push(marker);
}
function deleteMarkers() {
    for (let i = 0; i < markers.length; i++) {
        markers[i].setMap(null);
    }
    markers = [];
}
// print route on map & mark recommand stop points
function calculateAndDisplayRoute(directionsService, directionsRenderer) {
    directionsService
        .route({
            origin: {
                query: document.getElementById('Start').value
            },
            destination: {
                query: document.getElementById('End').value
            },
            waypoints: stops,
            travelMode: google.maps.TravelMode.DRIVING
        })
        .then((response) => {
            directionsRenderer.setDirections(response);
            console.log(stops);
            if (stops.length == 0) {
                searchpts = decodePolyline(response.routes[0].overview_polyline);
                for (var i = 0; i < searchpts.length; i++) {
                    getNearbyPlaces(searchpts[i]);
                }
            }
        })
        .catch((e) => console.log('Directions request failed due to ' + e));
}
function DateConvert(dateString) {
    const date = new Date(dateString);
    // 取得日期和時間的部分
    const year = date.getFullYear(); // 年份
    const month = `0${date.getMonth() + 1}`.slice(-2); // 月份（補0）
    const day = `0${date.getDate()}`.slice(-2); // 日（補0）
    const hours = `0${date.getHours()-8}`.slice(-2); // 小時（補0）
    const minutes = `0${date.getMinutes()}`.slice(-2); // 分鐘（補0）
    // 格式化成"YYYY-MM-DDTHH:MM:SSZ"的形式
    const formattedDate = `${year}-${month}-${day}T${hours}:${minutes}:00Z`;
    console.log(hours)
    return formattedDate;
}
function tripData() {
    const nowTripData = {
        destination: destination,
        mid: mid,
        passenger_limit: Number(document.getElementById('passengernumber').value),
        source: source,
        start_time: DateConvert(document.getElementById('departuretime').querySelector('input').value)
    };
    //console.log(document.getElementById("departuretime").value);
    return nowTripData;
}
initMap();
</script>

<template>
    <div class="grid p-fluid">
        <div class="col-12">
            <h3>Create A Trip</h3>
            <div class="card">
                <!-- <h5>License Plate Number</h5>
                <span class="p-float-label">
                    <InputNumber id="licenseplatenumber" type="text" v-model="licenseplatenumberValue" />
                    <label for="licenseplatenumber">License Plate Number</label>
                </span> -->
                <h5>Passenger Number</h5>
                <InputText id="passengernumber" placeholder="e.g. 4(包含司機自己)" v-model="passengernumberValue" />
                <h5>From</h5>
                <InputText class="search-location" placeholder="Search" id="Start" type="text" />
                <h5>To</h5>
                <InputText class="search-location" placeholder="Search" id="End" type="text" />

                <h5>Middle Station</h5>
                <div class="flex flex-row">
                    <input id="Stop" type="text" class="search-location" placeholder="Search" />
                    <Button id="addStop" label="+" class="save-stop-btn"></Button>
                </div>
                <textarea id="Stops" class="stop-list-textarea" rows="5" cols="30" readonly></textarea>

                <Button label="CheckPath" class="mr-2 mb-2" id="checkPath"></Button>

                <br /><br />
                <div id="map" style="width: 100%; height: 30vh"></div>
                <br />
                <h5>Departure Time</h5>
                <Calendar id="departuretime" :showIcon="true" :showButtonBar="true" v-model="datetime24h1" showTime hourFormat="24"></Calendar>
                <br />
                <Button label="Save" class="mr-2 mb-2" id="Save"></Button>
            </div>
        </div>
    </div>
</template>

<style>
.stop-container {
    max-width: 400px;
    margin: auto;
}

.search-location {
    width: 100%;
    height: 40px;
    padding: 10px;
    margin-bottom: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
}

.stop-list {
    list-style-type: none;
    padding: 0;
    margin-bottom: 10px;
    border: 1px solid #ddd;
    background-color: #f9f9f9;
}

.stop-list li {
    padding: 8px;
    border-bottom: 1px solid #eee;
}

.stop-list li:last-child {
    border-bottom: none;
}

.save-stop-btn {
    margin-left: 10px;
    width: 40px;
    height: 40px;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

.stop-list-textarea {
    width: 100%; /* Full width */
    padding: 8px; /* Padding for some space inside */
    margin-bottom: 10px; /* Margin at the bottom */
    border: 1px solid #ccc; /* A light border similar to other inputs */
    border-radius: 4px; /* Rounded corners */
    background-color: #f9f9f9; /* Light grey background */
    font-family: Arial, sans-serif; /* Font style */
    font-size: 1rem; /* Readable font size */
    line-height: 1.5; /* Spacing between lines */
    color: #333; /* Font color */
    resize: vertical; /* Allow only vertical resizing */
}

.stop-list-textarea[readonly] {
    background-color: #e9ecef; /* Slightly different background for readonly */
    cursor: default; /* Default cursor for readonly */
}
</style>
