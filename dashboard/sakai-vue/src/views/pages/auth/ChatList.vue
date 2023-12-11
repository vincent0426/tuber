<script setup>
import { ref, onMounted } from 'vue';
import { TripService, DriverService } from '@/service';
const tripService = new TripService();
const driverService = new DriverService();
const chatList = ref([]);
// const tripIdList = ref([]);

onMounted(async () => {
    try {
        const parsedData = JSON.parse(localStorage.getItem('vuex-state'));
        const role = parsedData.role;
        const isDriver = role === 'passenger' ? false : true;
        console.log('isDriver?', isDriver);
        const now_trip = await tripService.getHistory({ trip_status: 'in_trip', is_driver: isDriver });
        const not_start_trip = await tripService.getHistory({ trip_status: 'not_start', is_driver: isDriver });
        const finished_trip = await tripService.getHistory({ trip_status: 'finished', is_driver: isDriver });
        const response = finished_trip;
        const resultTrips = ref([]);
        //console.log('my all trip history', response);
        // tripIdList.value = now_trip.items.map((item) => item.TripID);
        // console.log('my all tripID', tripIdList[0]);
        generateRandomData(response);
        chatList.value = response.items;
    } catch (e) {
        console.error('Error fetching History:', e);
    }
});

// // 处理每个 tripIdList
// const processTripList = (tripList, color) => {
//     tripList.forEach((tripId) => {
//         const driverName = getDriverName(tripId);
//         const passengerNames = getPassengerNames(tripId);

//         // 将结果添加到 resultTrips
//         resultTrips.value.push({
//             tripId,
//             color,
//             driverName,
//             passengerNames
//         });
//     });
// };

// // 模拟异步获取 driver 名称的函数
// const getDriverName = async (tripId) => {
//     // 模拟异步操作，实际应用中需要替换为真实的异步操作
//     return new Promise((resolve) => {
//         setTimeout(() => {
//             resolve(`Driver-${tripId}`);
//         }, 500);
//     });
// };

// // 模拟异步获取 passenger 名称的函数
// const getPassengerNames = async (tripId) => {
//     // 模拟异步操作，实际应用中需要替换为真实的异步操作
//     return new Promise((resolve) => {
//         setTimeout(() => {
//             const count = Math.floor(Math.random() * 3) + 1; // 随机生成 1 到 3 个乘客
//             const passengerNames = Array.from({ length: count }, (_, index) => `Passenger-${index + 1}`);
//             resolve(passengerNames);
//         }, 500);
//     });
// };

// processTripList(redTripIdList.value, 'red');
// processTripList(yellowTripIdList.value, 'yellow');
// processTripList(blueTripIdList.value, 'blue');

const generateRandomData = (response) => {
    response.items.forEach((item) => {
        item.trip_user_list = getTripPassengers(item.TripID);
    });
};

const getTripPassengers = (tripId) => {
    const resp = tripService.getPassenger(tripId);
    console.log('trip_passengers:', resp);
};
</script>
<template>
    <div>
        <h3 style="text-align: center">Chatroom List</h3>
        <div class="ride-container">
            <div v-for="(ride, index) in chatList" :key="index" class="ride-card">
                <div class="custom-content">
                    <div class="left-part">
                        <div class="left-upper">
                            <div class="avatar-container">
                                <img alt="driver avatar" src="../../../assets/images/Patrick.svg" class="avatar" />
                            </div>
                            <div class="driver-name">Driver: {{ ride.DriverName }}</div>
                        </div>
                        <div style="width: 270px; height: min-content; margin-right: 10px; margin-bottom: 10px">
                            <div style="font-weight: 600; padding-top: 10px; font-style: italic">Passengers: {{ ride.trip_user_list }}</div>
                        </div>
                    </div>

                    <div class="right-part">
                        <router-link :to="{ name: 'ChatRoom', params: { tripId: ride.TripID } }">
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
    align-items: center;
    justify-content: center;
    align-content: center;
    flex-wrap: wrap;
}

.avatar {
    width: 100%;
    height: 100%;
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
