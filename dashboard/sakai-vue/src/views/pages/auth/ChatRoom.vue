<template>
    <div>
        <div>Trip ID: {{ $route.params.tripId }}</div>
        <h3 style="text-align: center">Chat Room</h3>
        <div class="message-container">
            <div v-for="(message, index) in messages" :key="index" class="message">
                <strong>{{ message.sender }}</strong
                >: {{ message.text }}
            </div>
        </div>
        <div class="input-container">
            <input v-model="newMessage" @keyup.enter="sendMessage" placeholder="Type your message..." />
            <button @click="sendMessage">Send</button>
        </div>
    </div>
</template>
<script setup>
import { useStore } from 'vuex';
const store = useStore();

// 使用 getter 取得 user
const user = store.getters.user;
console.log(user);
</script>
<script>
export default {
    name: 'TripDetails',
    data() {
        return {
            tripId: null,
            socket: null,
            latestMessage: ''
        };
    },
    created() {
        // Set the trip ID from the route parameter
        this.tripId = this.$route.params.tripId;

        // Initialize the WebSocket connection
        this.initializeWebSocket();
    },
    methods: {
        initializeWebSocket() {
            // Replace 'ws://example.com/path' with your WebSocket server URL
            this.socket = new WebSocket(`ws://localhost:3000/v1/chat/ws?room=${this.tripId}`);

            this.socket.addEventListener('open', (event) => {
                console.log('WebSocket is open now.');
                // You can send a message to the server if needed
                // this.socket.send('Hello Server!');
            });

            this.socket.addEventListener('message', (event) => {
                console.log('Message from server:', event.data);
                this.latestMessage = event.data;
            });

            this.socket.addEventListener('error', (event) => {
                console.error('WebSocket error observed:', event);
            });

            this.socket.addEventListener('close', (event) => {
                console.log('WebSocket is closed now.');
            });
        }
    },
    beforeDestroy() {
        // Close the WebSocket connection when the component is destroyed
        if (this.socket) {
            this.socket.close();
        }
    }
};
</script>

<style scoped>
.message-container {
    max-height: 300px;
    overflow-y: auto;
    border: 1px solid #ccc;
    padding: 10px;
}

.message {
    margin-bottom: 5px;
}

.input-container {
    margin-top: 10px;
}

input {
    padding: 5px;
    margin-right: 5px;
}

button {
    padding: 5px;
}
</style>
