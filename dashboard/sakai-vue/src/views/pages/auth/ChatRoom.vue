<script setup>
import { useStore } from 'vuex';
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const store = useStore();
const user = store.getters.user;
const route = useRoute();
const tripId = ref(route.params.tripId);
const socket = ref(null);
const messages = ref([]);
const newMessage = ref('');

const sendMessage = () => {
    if (newMessage.value.trim() !== '') {
        socket.value.send(JSON.stringify({ text: newMessage.value, sender: user.name }));
        newMessage.value = '';
    }
};

const processMessage = (rawMessage) => {
    const parsedMessage = JSON.parse(rawMessage);
    const messageContent = JSON.parse(parsedMessage.Message);
    messages.value.push({
        UserID: parsedMessage.UserID,
        Username: parsedMessage.Username,
        ImageURL: parsedMessage.ImageURL,
        MessageText: messageContent.text
    });
};

const initializeWebSocket = () => {
    socket.value = new WebSocket(`ws://localhost:3000/v1/chat/ws?room=${tripId.value}&user_id=${user.id}`);

    socket.value.addEventListener('open', (event) => {
        console.log('WebSocket is open now.', event);
    });

    socket.value.addEventListener('message', (event) => {
        console.log('Message from server:', event.data);
        processMessage(event.data);
    });

    socket.value.addEventListener('error', (event) => {
        console.error('WebSocket error observed:', event);
    });

    socket.value.addEventListener('close', (event) => {
        console.log('WebSocket is closed now.', event);
    });
};

onMounted(initializeWebSocket);

onBeforeUnmount(() => {
    if (socket.value) {
        socket.value.close();
    }
});
</script>

<template>
    <div>
        <!-- <div>Trip ID: {{ tripId }}</div> -->
        <h3 style="text-align: center">Chat Room</h3>
        <div class="message-container">
            <div v-for="(message, index) in messages" :key="index" class="message">
                <img :src="message.ImageURL" class="message-avatar" />
                <div class="message-content">
                    <strong>{{ message.Username }}</strong
                    >: {{ message.MessageText }}
                </div>
            </div>
        </div>
        <div class="input-container">
            <input v-model="newMessage" @keyup.enter="sendMessage" placeholder="Type your message..." />
            <!-- <button @click="sendMessage">Send</button> -->
            <div class="send-button-container">
                <button @click="sendMessage" class="pi pi-send" style="color: bisque"></button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.send-button-container {
    border-radius: 20%;
    background-color: rgba(0, 0, 0, 0.7);
    width: 60px;
    height: 30px;
    display: flex;
    justify-content: center;
    align-items: center;
}
.send-button-container button {
    padding: 0;
    margin: 0;
    margin-right: 8px;
    font-size: 15px;
    border: none;
    background: none;
    cursor: pointer;
}
.message-container {
    overflow-y: auto;
    border: 1px solid #cecbcb;
    padding: 10px;
    width: 380px; /* Set the maximum width of the container */
    height: 660px;
    margin: 0 auto; /* Center the container horizontally */
    overflow-y: auto; /* Hide content if it exceeds the height of the container */
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
}
.message {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
}

.message-avatar {
    width: 40px;
    height: 40px;
    border-radius: 20px;
    margin-right: 10px;
}

.message-content {
    flex: 1;
}

.input-container {
    position: fixed;
    bottom: 35px;
    left: 0;
    right: 0;
    margin: 0 auto;
    display: flex;
    flex-direction: row;
    background-color: white; /* 背景色 */
    padding: 15px;
    width: 360px;
    box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
    border-radius: 5px;
}

input {
    padding: 5px;
    margin-right: 5px;
    width: 80%; /* Adjust as needed */
}

button {
    padding: 5px;
    width: 15%; /* Adjust as needed */
}
</style>
