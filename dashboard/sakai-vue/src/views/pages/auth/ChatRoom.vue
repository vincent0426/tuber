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
        <div>Trip ID: {{ tripId }}</div>
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
            <button @click="sendMessage">Send</button>
        </div>
    </div>
</template>

<style scoped>
.message-container {
    max-height: 300px;
    overflow-y: auto;
    border: 1px solid #ccc;
    padding: 10px;
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
    margin-top: 10px;
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
