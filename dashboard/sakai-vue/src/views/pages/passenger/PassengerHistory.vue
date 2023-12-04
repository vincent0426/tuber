import Rating from 'primevue/rating';
<script setup>
import { ref } from 'vue';

const value = ref(null);
</script>

<script>
export default {
    data() {
        return {
            rideHistory: [
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: 'John Doe2',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    source_name: '台北',
                    destination_name: '新竹',
                    cost: '5',
                    rating: 1,
                    start_time: '2023-01-01T08:00:00Z'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: 'A',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    source_name: 'Location 1',
                    destination_name: 'Location 2',
                    cost: '5',
                    rating: 5,
                    start_time: '2023-01-01T08:00:00Z'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '我瘋子',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    source_name: '北車',
                    destination_name: '行天宮',
                    cost: '5',
                    rating: 4,
                    start_time: '2023-01-01T08:00:00Z'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    source_name: 'Location 1',
                    destination_name: 'Location 2',
                    cost: '5',
                    rating: 1,
                    start_time: '2023-01-01T08:00:00Z'
                }
            ]
        };
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
                            <p style="font-weight: bold; font-size: 16px">{{ ride.driver_name }}</p>
                            <p style="font-size: 13px; font-style: italic">{{ ride.driver_plate }}</p>
                        </div>
                        <img alt="user header" src="../../../assets/images/modelS.jpg" style="width: 130px; height: 100px; object-fit: cover; object-position: center" />
                    </div>
                </template>
                <Divider />
                <template #subtitle class="custom-content">{{ ride.source_name }} -> {{ ride.destination_name }}</template>
                <template #content class="custom-content">
                    <div style="background: rgba(128, 128, 128, 0.05); border-radius: 3px">
                        <p class="m-0">
                            {{ ride.start_time }}
                        </p>
                        <p class="m-0">花費: {{ ride.cost }}</p>
                    </div>
                </template>
                <template #footer>
                    <!-- <Rating v-model="value" readonly :cancel="false" /> -->
                    <div class="flex justify-content-center">
                        <Rating :modelValue="ride.rating" :stars="5" :cancel="false" />
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

// { "driver_id": "fe22e8fa-04d2-49b5-8bac-1535153b687e", "driver_name": "John Doe", "driver_brand": "Toyota", "driver_model": "Camry", "driver_plate": "ABC123", "source_name": "Location 1", "destination_name": "Location 2", "cost":"5", "start_time":
// "2023-01-01T08:00:00Z", }
