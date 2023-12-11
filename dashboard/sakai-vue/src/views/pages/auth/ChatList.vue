import Rating from 'primevue/rating';
<script setup>
import { ref } from 'vue';
const value = ref(null);
</script>

<script>
import { TripService, LocationService } from '@/service';
const tripService = new TripService();
const locationService = new LocationService();
export default {
    data() {
        return {
            chatList: []
        };
    },
    mounted() {
        this.fetchChatList();
    },
    methods: {
        generateRandomData(response) {
            response.items.forEach((item) => {
                console.log(item);
                console.log('TripId', item.TripID);
                item.trip_user_list = this.getTripPassengers(item.TripID);
                item.source_name = this.getLocationName(item.MySourceID);
                item.destination_name = this.getLocationName(item.MyDestinationID);
            });
        },
        async getTripPassengers(tripId) {
            const resp = tripService.getPassenger(tripId);
            console.log('trip_passengers:', resp);
        },
        async getLocationName(id) {
            const resp = await locationService.getLocationName(id);
            console.log('locNameGet', resp);
        },
        async fetchChatList() {
            try {
                const response = await tripService.getHistory({ trip_status: 'finished', is_driver: false });
                console.log('before', response);
                this.generateRandomData(response);
                console.log('after', response);
                this.chatList = response.items;
            } catch (e) {
                console.error('Error fetching History:', error);
            }
        }
    }
};
</script>
<template>
    <div>
        <h3 style="text-align: center">Chatroom List</h3>
        <div class="ride-container">
            <div v-for="(ride, index) in chatList" :key="index" class="ride-card">
                <div class="custom-content">
                    <div class="avatar-container">
                        <img alt="driver avatar" src="../../../assets/images/Patrick.svg" class="avatar" />
                    </div>
                    <div style="width: 270px; height: min-content; margin-right: 10px">
                        <div style="m-0; word-wrap: break-word; font-weight: 600;">Driver: {{ ride.DriverName }}</div>
                        <div style="font-weight: 600; padding-top: 10px">Passengers: {{ ride.trip_user_list }}</div>
                    </div>
                    <router-link :to="{ name: 'ChatRoom', params: { tripId: ride.TripID } }">
                        <div style="align-self: center; border-radius: 50%; background-color: rgba(0, 0, 0, 0.7); width: 55px; height: 55px; display: flex; align-items: center; justify-content: center">
                            <i class="pi pi-pencil" style="color: bisque"></i>
                        </div>
                    </router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.custom-content {
    background: rgba(128, 128, 128, 0.05);
    display: flex;
    flex-direction: row;
}
.ride-container {
    max-width: 380px; /* Set the maximum width of the container */
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
    width: 70px; /* 設定avatar容器的寬度 */
    height: 70px; /* 設定avatar容器的高度 */
}

.avatar {
    width: 100%;
    height: 100%;
    object-fit: cover; /* 保持頭像比例 */
    border-radius: 50%; /* 圓形頭像 */
    border: 2px solid #070707; /* 添加2px的邊框 */
}
</style>
