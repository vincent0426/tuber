<template>
    <div>
        <h3 style="text-align: center">Favorite Driver</h3>
        <div class="ride-container">
            <div v-for="(ride, index) in AppFavoriteDriver" :key="index" class="ride-card">
                <div class="driver-info">
                    <div class="avatar-container">
                        <img alt="driver avatar" src="../../../assets/images/Patrick.svg" class="avatar" />
                    </div>
                    <div class="driver-text">
                        <div class="name-rating">
                            <div style="font-weight: bold; font-size: 16px">{{ ride.driver_name }}</div>
                            <div class="rating-container">
                                <div>{{ ride.driver_rating }}</div>
                                <i class="pi pi-star-fill" style="color: black"></i>
                            </div>
                        </div>

                        <p style="font-size: 13px; font-style: italic; padding-top: 4px">{{ ride.driver_plate }}</p>
                    </div>
                    <img alt="user header" src="../../../assets/images/modelS.jpg" style="width: 130px; height: 100px; object-fit: cover; object-position: center" />
                </div>
                <div class="custom-content">
                    <div style="width: 270px; height: min-content; margin-right: 10px">
                        <div style="m-0; word-wrap: break-word;">Comment: {{ ride.comment }}</div>
                        <div style="font-weight: 600; padding-top: 10px">Collaboration time: {{ ride.collaboration_time }}</div>
                    </div>

                    <div style="align-self: center; border-radius: 50%; background-color: rgba(0, 0, 0, 0.7); width: 55px; height: 55px; display: flex; align-items: center; justify-content: center">
                        <i class="pi pi-pencil" style="color: bisque"></i>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.name-rating {
    display: flex;
    align-items: center;
    justify-content: space-between;
    text-align: center;
    margin: 0;
}
.rating-container {
    display: flex;
    align-items: center;
    gap: 5px; /* Adjust the gap between rating and star icon */
    margin: 0;
}
.custom-content {
    background: rgba(128, 128, 128, 0.05);
    padding: 10px;
    margin: 0 10px;
    min-height: 80px;
    border-radius: 3px;
    display: flex;
    align-items: center;
    flex-direction: row;
}
.ride-container {
    max-width: 380px; /* Set the maximum width of the container */
    max-height: 755px;
    min-height: fit-content;
    margin: 0 auto; /* Center the container horizontally */
    overflow-y: auto; /* Hide content if it exceeds the height of the container */
    overflow-x: hidden;
    display: flex;
    flex-wrap: wrap;
    gap: 15px; /* 設定Card之間的空隙 */
}

.ride-card {
    border: 0, 0, 0, 1px solid #ccc;
    padding: 10px;
    margin-bottom: 10px;
    flex: 1 0 25em;
    height: min-content;
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

<script setup>
import { ref } from 'vue';

import 'primeicons/primeicons.css';

const value = ref(null);
</script>
<script>
export default {
    data() {
        return {
            AppFavoriteDriver: [
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    driver_rating: 4.7,
                    collaboration_time: 2,
                    comment: 'A good driver, we often share the same route.'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: 'John Doe',
                    driver_brand: 'Tesla',
                    driver_model: 'Camry',
                    driver_plate: 'UDC-256',
                    driver_rating: 3.3,
                    collaboration_time: 5,
                    comment: 'A crazzzzzzzxychxycxychizhvhishvifhivhsivhfihvifsadfasdfafdsfsfafdafdsasdfadfafdasfasdfasdfsadfsdafdhvhsdivhdfihvidhadfasdffafsafdafasfafsfdasidfhvishfhishivhdfi.'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    driver_rating: 4.7,
                    collaboration_time: 2,
                    comment: 'A good driver, we often share the same route.'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    driver_rating: 4.7,
                    collaboration_time: 2,
                    comment: 'A good driver, we often share the same route.'
                },
                {
                    driver_id: 'fe22e8fa-04d2-49b5-8bac-1535153b687e',
                    driver_name: '爆肝人',
                    driver_brand: 'Toyota',
                    driver_model: 'Camry',
                    driver_plate: 'ABC123',
                    driver_rating: 4.7,
                    collaboration_time: 2,
                    comment: 'A good driver, we often share the same route.'
                }
            ]
        };
    }
};
</script>

// { "driver_id": "fe22e8fa-04d2-49b5-8bac-1535153b687e", "driver_name": "John Doe", "driver_brand": "Toyota", "driver_model": "Camry", "driver_plate": "ABC123", "source_name": "Location 1", "destination_name": "Location 2", "cost":"5", "start_time":
// "2023-01-01T08:00:00Z", }

<!-- type AppFavoriteDriver struct {
ID              string json:"id"
DriverID        string json:"driver_id"
DriverName      string json:"driver_name"
DriverImageURL  string json:"driver_image_url"
DriverBrand     string json:"driver_brand"
DriverModel     string json:"driver_model"
DriverColor     string json:"driver_color"
DriverPlate     string json:"driver_plate"
DriverCreatedAt string json:"driver_created_at"
}  -->
<!-- <template>
  <div>
      <h1>Favorite List</h1>
      <div v-for="trip in trips" :key="trip.id">
          <h2>{{ trip.driver_name }}</h2>
      </div>
  </div>
</template>

<script>
import axios from 'axios';
export default {
  data() {
      return {
          trips: []
      };
  },
  mounted() {
      this.fetchTrips();
  },
  methods: {
      async fetchTrips() {
          try {
              const response = await axios.get('http://localhost:3000/v1/trips');
              this.trips = response.data.items;
          } catch (error) {
              console.error('Error fetching trips:', error);
          }
      }
  }
};
</script> -->
