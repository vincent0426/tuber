<script setup>
import { ref, onMounted } from 'vue';
import { TripService, DriverService } from '@/service';
const tripService = new TripService();
const driverService = new DriverService();
const chatList = ref([]);

onMounted(async () => {
    try {
        const parsedData = JSON.parse(localStorage.getItem('vuex-state'));
        const role = parsedData.role;
        const isDriver = role === 'passenger' ? false : true;
        console.log('isDriver?', isDriver);
        const now_trip = await tripService.getHistory({ trip_status: 'in_trip', is_driver: isDriver });
        console.log('now_trip:', now_trip);
        const not_start_trip = await tripService.getHistory({ trip_status: 'not_start', is_driver: isDriver });
        console.log('not_started_trip:', not_start_trip);
        const finished_trip = await tripService.getHistory({ trip_status: 'finished', is_driver: isDriver });
        console.log('finished_trip:', finished_trip);
        const allTripsWithColor = processTripList(now_trip, not_start_trip, finished_trip);
        chatList.value = await iterateTripGetData(allTripsWithColor);
        console.log('chatList\n', chatList.value);
    } catch (e) {
        console.error('Error fetching History:', e);
    }
});

const processTripList = (now_trip, not_start_trip, finished_trip) => {
    let now_trip_result = now_trip.items.map((item) => ({ tripId: item.TripID, color: '#666633' }));
    let not_start_trip_result = not_start_trip.items.map((item) => ({ tripId: item.TripID, color: '#CCCCCC' }));
    let finished_trip_result = finished_trip.items.map((item) => ({ tripId: item.TripID, color: '#333333' }));
    return [...now_trip_result, ...not_start_trip_result, ...finished_trip_result];
};

const iterateTripGetData = async (trips) => {
    await Promise.all(
        trips.map(async (item) => {
            item.trip_user_list = await getTripPassengers(item.tripId);
            item.trip_driver_name = item.trip_user_list.driver_details.driver_name;
            const passengerNames = item.trip_user_list.passenger_details
                .map((passenger) => passenger.passenger_name)
                .filter((passengerName) => passengerName !== item.trip_driver_name)
                .join(', ');
            item.trip_passengers_name = passengerNames;
        })
    );
    return trips;
};

const getTripPassengers = async (tripId) => {
    try {
        const resp = await tripService.getPassengers(tripId);
        return resp;
        // console.log('trip_passengers:', resp);
    } catch (error) {
        console.error('Error fetching trip passengers:', error);
    }
};
const getBackgroundColor = (color) => {
    return `linear-gradient(to bottom left, ${color}, #FFFFFF 30%, #FFFFFF)`;
};
</script>
<template>
    <div style="padding-top: 6px">
        <h3 style="text-align: center">Chatroom List</h3>
        <div class="ride-container">
            <div v-for="(ride, index) in chatList" :key="index" class="ride-card" :style="{ background: getBackgroundColor(ride.color) }">
                <div class="custom-content">
                    <div class="left-part">
                        <div class="left-upper">
                            <div class="avatar-container">
                                <img alt="driver avatar" :src="ride.trip_user_list.driver_details.driver_image_url" class="avatar" />
                            </div>
                            <div style="display:flex flex-direction:column">
                                <div class="driver-name">Driver: {{ ride.trip_driver_name }}</div>
                                <div class="locations">{{ ride.trip_user_list.source_name }} -> {{ ride.trip_user_list.destination_name }}</div>
                            </div>
                        </div>
                        <div style="display: flex; justify-content: flex-start; width: 270px; height: min-content; margin-right: 10px; margin-bottom: 10px">
                            <div style="font-weight: 600; padding-top: 10px; font-style: italic">Passengers: {{ ride.trip_passengers_name }}</div>
                        </div>
                    </div>

                    <div class="right-part">
                        <router-link :to="{ name: 'ChatRoom', params: { tripId: ride.tripId } }">
                            <div style="align-items: center; border-radius: 50%; background-color: rgba(0, 0, 0, 0.7); width: 40px; height: 40px; display: flex; justify-content: center">
                                <i class="pi pi-comments" style="color: bisque"></i>
                            </div>
                        </router-link>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.locations {
    font-weight: 200;
    font-size: small;
    margin-left: 20px;
    font-style: italic;
}
.custom-content {
    background: rgba(128, 128, 128, 0.05);
    display: flex;
    flex-direction: row;
    height: 100%;
    width: 100%;
}
.ride-container {
    width: 380px; /* Set the maximum width of the container */
    max-height: 755px;
    margin: 0 auto; /* Center the container horizontally */
    overflow-y: auto; /* Hide content if it exceeds the height of the container */
    display: flex;
    flex-wrap: wrap;
    gap: 15px; /* 設定Card之間的空隙 */
}

.ride-card {
    border: 0, 0, 0, 1px solid #ccc;
    padding: 10px;
    margin-bottom: 10px;
    flex: 1 0 25em;
    width: 100%;
}

.driver-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100px;
    padding-left: 10px;
    padding-right: 20px;
}

.driver-text {
    flex-grow: 1;
    padding-left: 10px;
    padding-right: 10px; /* 添加右側空隙 */
}
.avatar-container {
    width: 40px; /* 設定avatar容器的寬度 */
    height: 40px; /* 設定avatar容器的高度 */
    margin-right: 10px;
}
.driver-name {
    text-align: center; /* 将这个样式应用到 driver-name */
    font-weight: 600;
    font-size: larger;
    /* align-items: center; */
    justify-content: flex-start;
    align-content: flex-start;
    flex-wrap: wrap;
    margin-top: 10px;
    padding-left: -30px;
    align-self: flex-start;
}

.avatar {
    width: 50px;
    height: 50px;
    object-fit: cover; /* 保持頭像比例 */
    border-radius: 50%; /* 圓形頭像 */
    border: 2px solid #070707; /* 添加2px的邊框 */
}
.left-part {
    display: flex;
    flex-direction: column;
}
.right-part {
    width: 40px;
    height: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    /* justify-content: center; */
    align-content: center;
    /* flex-wrap: wrap; */
    right: 0px;
}
.left-upper {
    display: flex;
    flex-direction: row;
    align-items: center;
    /* justify-content: center; */
    align-content: center;
    /* flex-wrap: wrap; */
}
</style>
