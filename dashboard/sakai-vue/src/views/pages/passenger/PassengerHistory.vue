import Rating from 'primevue/rating';
<script setup>
import { ref } from 'vue';

const value = ref(null);
</script>

<script>
import { TripService, LocationService, DriverService } from '@/service';
const tripService = new TripService();
const locationService = new LocationService();
const driverService = new DriverService();
export default {
    data() {
        return {
            rideHistory: []
        };
    },
    mounted() {
        this.fetchHistory();
    },
    methods: {
        generateRandomCost() {
            const comments = ['3', '5', '10', '12', '14', '15', '20', '24', '31', '55', '84', '85'];
            return comments[Math.floor(Math.random() * comments.length)];
        },
        async generateRandomData(response) {
            const favoriteDriversResponse = await driverService.getFavorite();
            const favoriteDrivers = favoriteDriversResponse.items.map((driver) => driver.driver_id);
            // 使用 Promise.all 等待所有的非同步操作完成
            await Promise.all(
                response.items.map(async (item) => {
                    item.Cost = this.generateRandomCost();
                    item.driver_plate = await this.getDriverPlate(item.DriverID);
                    if (favoriteDrivers.includes(item.DriverID)) {
                        //確認此人是否已經在最愛司機列表
                        item.isFavorite = true;
                    } else {
                        item.isFavorite = false;
                    }
                })
            );
        },
        async getDriverPlate(id) {
            const resp = await driverService.getDriver(id);
            console.log('getDriverPlateResponse:', resp.plate);
            return resp.plate;
        },
        async fetchHistory() {
            try {
                const response = await tripService.getHistory({ trip_status: 'finished', is_driver: false });
                await this.generateRandomData(response);
                // await this.checkInFavorite(response);
                console.log(response);
                this.rideHistory = response.items;
            } catch (e) {
                console.error('Error fetching History:', error);
            }
        },
        formatTime(isoDateString) {
            const date = new Date(isoDateString);
            const options = {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            };

            return date.toLocaleTimeString([], options);
        },
        async addToFavorites(ride) {
            try {
                //加入driver
                const response = await driverService.postFavorite(ride.DriverID);
                console.log(response);
                ride.isFavorite = true;
            } catch (e) {
                console.error('Error posting fav driver:', error);
            }
        }
    }
};
</script>
<template>
    <div>
        <h3 style="text-align: center">History</h3>
        <div class="ride-container">
            <Card v-for="(ride, index) in rideHistory" :key="index" class="ride-card">
                <template #header>
                    <div class="driver-info">
                        <div class="avatar-container">
                            <img alt="driver avatar" src="../../../assets/images/Patrick.svg" class="avatar" />
                        </div>
                        <div class="driver-text">
                            <p style="font-weight: bold; font-size: 16px">{{ ride.DriverName }}</p>
                            <p style="font-size: 13px; font-style: italic">{{ ride.driver_plate }}</p>
                        </div>
                        <img alt="user header" src="../../../assets/images/modelS.jpg" style="width: 130px; height: 100px; object-fit: cover; object-position: center" />
                    </div>
                </template>
                <Divider />
                <template #subtitle class="custom-content">{{ ride.SourceName }} -> {{ ride.DestinationName }}</template>
                <template #content class="custom-content">
                    <div style="background: rgba(128, 128, 128, 0.05); border-radius: 3px; padding: 10px">
                        <p class="m-0">Comment: {{ ride.Comment }}</p>
                        <p class="m-0">
                            {{ formatTime(ride.StartTime) }}
                        </p>
                        <p class="m-0">Ride cost: {{ ride.Cost }}</p>
                    </div>
                </template>
                <template #footer>
                    <!-- <Rating v-model="value" readonly :cancel="false" /> -->
                    <div class="flex justify-content-center">
                        <Rating :modelValue="ride.Rating" :stars="5" :cancel="false" />
                    </div>
                    <div class="flex justify-content-center mt-2">
                        <Button v-if="!ride.isFavorite" @click="addToFavorites(ride)">Add to Favorite</Button>
                        <Button v-else disabled>Added to Favorite</Button>
                    </div>
                    <!-- https://primevue.org/rating/ -->
                </template>
            </Card>
        </div>
    </div>
</template>
<style scoped>
.custom-content {
    background: rgba(128, 128, 128, 0.05);
    padding: 10px;
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
// export default { // data() { // return { // rideHistory: [ // { // driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e', // driver_name: 'John Doe2', // driver_brand: 'Toyota', // driver_model: 'Camry', // driver_plate: 'ABC123', // source_name:
'台北', // destination_name: '新竹', // cost: '5', // rating: 1, // start_time: '2023-01-01T08:00:00Z' // }, // { // driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e', // driver_name: 'A', // driver_brand: 'Toyota', // driver_model: 'Camry', //
driver_plate: 'ABC123', // source_name: 'Location 1', // destination_name: 'Location 2', // cost: '5', // rating: 5, // start_time: '2023-01-01T08:00:00Z' // }, // { // driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e', // driver_name: '我瘋子', //
driver_brand: 'Toyota', // driver_model: 'Camry', // driver_plate: 'ABC123', // source_name: '北車', // destination_name: '行天宮', // cost: '5', // rating: 4, // start_time: '2023-01-01T08:00:00Z' // }, // { // driver_id:
'fe22e8fa-04d2-49b5-8bac-1535153b687e', // driver_name: '爆肝人', // driver_brand: 'Toyota', // driver_model: 'Camry', // driver_plate: 'ABC123', // source_name: 'Location 1', // destination_name: 'Location 2', // cost: '5', // rating: 1, //
start_time: '2023-01-01T08:00:00Z' // } // ] // }; // } // };
