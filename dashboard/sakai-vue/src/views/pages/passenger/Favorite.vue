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
    width: 50px;
    height: 50px;
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
import { DriverService } from '@/service';
const driverService = new DriverService();
export default {
    data() {
        return {
            FavoriteDriver: []
        };
    },
    mounted() {
        this.fetchFavorite();
    },
    methods: {
        generateRandomComment() {
            const comments = ['good driver, YEAH', 'A good driver, we often share the same route.', 'cool driver, drives fast'];
            return comments[Math.floor(Math.random() * comments.length)];
        },
        generateRandomRating() {
            const ratings = [2.4, 4.8, 5.0];
            return ratings[Math.floor(Math.random() * ratings.length)];
        },
        generateRandomCollborationTime() {
            const col_times = [2, 5, 14, 31, 50];
            return col_times[Math.floor(Math.random() * col_times.length)];
        },
        generateRandomData() {
            // 为每个驾驶员生成随机评论和评分
            this.responseData.items.forEach((item) => {
                console.log(item);
                item.comment = this.generateRandomComment();
                item.driver_rating = this.generateRandomRating();
                item.collaboration_time = this.generateRandomCollborationTime();
            });
        },
        async fetchFavorite() {
            try {
                const response = await driverService.getFavorite();
                console.log(response);
                this.responseData = response;
                // 生成随机评论和评分
                this.generateRandomData();
                console.log(response);
                this.FavoriteDriver = response.items;
            } catch (e) {
                console.error('Error fetching fav drivers:', error);
            }
        }
    }
};
</script>

<template>
    <div>
        <h3 style="text-align: center">Favorite Driver</h3>
        <div class="ride-container">
            <div v-for="(ride, index) in FavoriteDriver" :key="index" class="ride-card">
                <div class="driver-info">
                    <div class="avatar-container">
                        <img alt="driver avatar" :src="ride.driver_image_url" class="avatar" />
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
