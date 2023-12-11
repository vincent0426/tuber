import Rating from 'primevue/rating';
<script setup></script>

<script>
import { TripService, LocationService } from '@/service';
const tripService = new TripService();
const locationService = new LocationService();
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
        generateRandomRevenue() {
            const comments = ['31', '52', '10', '12', '14', '15', '20', '24', '31', '55', '84', '85'];
            return comments[Math.floor(Math.random() * comments.length)];
        },
        generateRandomData(response) {
            response.items.forEach((item) => {
                item.received_revenue = this.generateRandomRevenue();
                // item.source_name = this.getLocationName(item.MySourceID);
                // item.destination_name = this.getLocationName(item.MyDestinationID);
            });
        },
        async fetchHistory() {
            try {
                const response = await tripService.getHistory({ trip_status: 'finished', is_driver: false });
                this.generateRandomData(response);
                console.log(response);
                this.rideHistory = response.items;
            } catch (e) {
                console.error('Error fetching History:', error);
            }
        }
    }
};
</script>
<!-- <script>
export default {
    data() {
        return {
            rideHistory: [
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    source_name: '桃園',
                    destination_name: '新竹',
                    received_revenue: '20',
                    customers: [
                        {
                            customer_img: '../../../assets/images/Sandy.jpg',
                            customer_name: 'Sandy',
                            customer_source_name: '北車',
                            customer_destination_name: '公館',
                            collaboration_time: 2
                        },
                        {
                            customer_img: '../../../assets/images/Patrick.svg',
                            customer_name: 'Patrick',
                            customer_source_name: '太平洋',
                            customer_destination_name: '大西洋',
                            collaboration_time: 8
                        }
                    ],
                    start_time: '2023-01-01T08:00:00Z'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    source_name: '台北',
                    destination_name: '墾丁',
                    received_revenue: '200',
                    customers: [
                        {
                            customer_img: '../../../assets/images/Sandy.jpg',
                            customer_name: 'Sandy',
                            customer_source_name: '北車',
                            customer_destination_name: '墾丁',
                            collaboration_time: 2
                        },
                        {
                            customer_img: '../../../assets/images/Patrick.jpg',
                            customer_name: 'Patrick',
                            customer_source_name: '太平洋',
                            customer_destination_name: '大西洋',
                            collaboration_time: 8
                        }
                    ],
                    start_time: '2023-01-01T08:00:00Z'
                }
            ]
        };
    }
};
</script> -->
<template>
    <div>
        <h3 style="text-align: center">Driver History</h3>
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
                <template #subtitle class="custom-content">
                    {{ ride.SourceName }} -> {{ ride.DestinationName }}
                    <div style="font-weight: bold; font-size: 20px; margin-bottom: -25px">Customers:</div>
                </template>
                <template #content class="custom-content">
                    <div style="background: rgba(128, 128, 128, 0.05); border-radius: 3px">
                        <div v-for="(customer, cIndex) in ride.customers" :key="cIndex">
                            <div class="driver-info">
                                <div class="avatar-container">
                                    <img alt="customer avatar" src="../../../assets/images/Sandy.jpg" class="avatar" />
                                </div>
                                <div>
                                    <p style="font-weight: bold; font-size: 15px; margin-bottom: -2px">{{ customer.customer_name }}</p>
                                    <div class="location-content">
                                        <div style="display: flex; flex-direction: row; align-items: center">
                                            <div class="location-marker"></div>
                                            <p>{{ customer.customer_source_name }}</p>
                                        </div>
                                        <div style="margin-left: 13px; margin-right: 13px; text-align: center">到</div>
                                        <div style="display: flex; flex-direction: row; align-items: center">
                                            <div class="location-marker"></div>
                                            <p>{{ customer.customer_destination_name }}</p>
                                        </div>
                                    </div>
                                    <p style="font-size: 12px">Collaboration Time: {{ customer.collaboration_time }} hours</p>
                                </div>
                            </div>
                            <div class="horizontal-divider" :style="{ innerHeight: cIndex === ride.customers.length - 1 ? '0' : '10px' }"></div>
                        </div>
                    </div>
                </template>
                <template #footer>
                    <div style="height: 50px">
                        <p class="flex justify-content-center">
                            {{ ride.start_time }}
                        </p>
                        <div class="flex justify-content-center">
                            <p style="font-weight: bold; font-size: 16px" class="m-0">Total revenue: {{ ride.received_revenue }}</p>
                        </div>
                    </div>
                </template>
            </Card>
        </div>
    </div>
</template>
<style scoped>
.location-content {
    display: flex;
    align-items: center;
}
.horizontal-divider {
    width: 100%;
    height: 1px;
    background-color: rgba(56, 75, 57, 0.8);
}
.location-marker {
    width: 10px;
    height: 10px;
    background-color: #4caf50; /* Green color */
    border-radius: 50%; /* Make it a circle */
    margin-right: 10px; /* Adjust spacing between the marker and location name */
    box-shadow: 0 0 10px rgba(76, 175, 80, 0.8); /* Glowing effect */
}
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
    justify-content: flex-start;
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
    margin-right: 20px;
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
