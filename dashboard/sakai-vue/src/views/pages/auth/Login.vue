<script setup>
import { onMounted, ref } from 'vue';
import AppConfig from '@/layout/AppConfig.vue';
import store from '@/store';
import router from '@/router';

const loading = ref(false);
const username = ref('');
const password = ref('');
const checked = ref(false);
const googleLoginBtnRef = ref(null);

onMounted(async () => {
    const gClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
    try {
        window.google.accounts.id.initialize({
            client_id: gClientId,
            callback: handleCredentialResponse,
            auto_select: true
        });

        window.google.accounts.id.renderButton(googleLoginBtnRef.value, { text: 'signin_with', size: 'large', width: '300', theme: 'outline', logo_alignment: 'left' });
    } catch (error) {
        console.error('Google login initialization failed:', error);
    }
});

const handleCredentialResponse = async (response) => {
    loading.value = true;
    try {
        await store.dispatch('login', { id_token: response.credential });
        router.push({ name: 'PassengerHome' });
    } catch (error) {
        // Handle login error
        console.error('Login failed:', error);
        alert('Login failed. Please check your username and password and try again.');
    } finally {
        loading.value = false;
    }
};
</script>

<template>
    <div class="flex align-items-center justify-content-center mt-8">
        <div class="flex flex-column align-items-center justify-content-center">
            <div style="border-radius: 56px; padding: 0.3rem; background: linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)">
                <div class="w-full surface-card py-6 px-5 sm:px-8" style="border-radius: 53px">
                    <div class="text-center mb-8">
                        <div class="text-900 text-5xl font-bold mb-4">TUber</div>
                    </div>

                    <div>
                        <label for="username1" class="block text-900 text-xl font-medium mb-2">Username</label>
                        <InputText id="username1" type="text" placeholder="Email address" class="w-full md:w-30rem mb-5" style="padding: 1rem" v-model="username" />

                        <label for="password1" class="block text-900 font-medium text-xl mb-2">Password</label>
                        <Password id="password1" v-model="password" placeholder="Password" :toggleMask="true" class="w-full mb-3" inputClass="w-full" :inputStyle="{ padding: '1rem' }"></Password>

                        <div class="flex align-items-center justify-content-between mb-5 gap-5">
                            <div class="flex align-items-center">
                                <Checkbox v-model="checked" id="rememberme1" binary class="mr-2"></Checkbox>
                                <label for="rememberme1">Remember me</label>
                            </div>
                            <a class="font-medium no-underline ml-2 text-right cursor-pointer" style="color: var(--primary-color)">Forgot password?</a>
                        </div>
                        <div class="flex gap-5">
                            <div ref="googleLoginBtnRef"></div>
                        </div>
                    </div>
                </div>
                <input type="tel" name="tel" placeholder=" Number" v-model="phone" />
            </div>
            <div class="options">
                <input type="text" class="search-box" placeholder="Search Country Name" />
                <ol></ol>
            </div>
        </div>

        <button class="arrow-button black-button" @click="goNext"><span>Next</span> <span class="arrow">→</span></button>
        <p class="p-message">By continuing you may receive an SMS for verification. Message and data rates may apply</p>

        <div class="divider">
            <span class="divider-text">or</span>
        </div>
        <div>
            <div class="flex gap-5">
                <div ref="googleLoginBtnRef"></div>
            </div>
        </div>
    </div>
</template>

<style>
p {
    text-align: left;
}
.subtitle {
    font-size: 16px;
    color: black;
    font-weight: bold;
}
.full-width-input {
    width: 98%;
    height: 32px;
    box-sizing: border-box;
}

.container {
    padding-left: 20px;
}
.arrow-button {
    padding: 10px 20px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 5px;
    font-size: 24px;
    cursor: pointer;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative; /* Needed for absolute positioning of the arrow */
    width: 90%; /* Adjust the width as needed */
    margin: 10px;
    margin-top: 20px;
    margin-bottom: 20px;
}
.black-button {
    background-color: black;
    color: white;
}

.arrow-button .arrow {
    position: absolute;
    right: 10px;
}
.p-message {
    color: #757575;
    font-size: 14px;
}

.select-box {
    position: relative;
}

.select-box input {
    background-color: #eee;
    width: 100%;
    padding: 1rem 0.6rem;
    font-size: 1.1rem;

    border: 0.1rem solid transparent;
    outline: none;
}
.gap-5 {
    text-align: center;
    display: flex;
}
input[type='tel'] {
    border-radius: 0 0.5rem 0.5rem 0;
}

.select-box input:focus {
    border: 0.1rem solid var(--primary);
}
.selected-option input {
    width: 75%;
}
.selected-option {
    border-radius: 0.5rem;
    overflow: hidden;

    display: flex;
    justify-content: space-between;
    align-items: center;
}

.selected-option div {
    position: relative;

    width: 6rem;
    padding: 0 2.8rem 0 0.5rem;
    text-align: center;
    cursor: pointer;
}

.selected-option div::after {
    position: absolute;
    content: '';
    right: 0.8rem;
    top: 50%;
    transform: translateY(-50%) rotate(45deg);

    width: 0.8rem;
    height: 0.8rem;
    border-right: 0.12rem solid var(--primary);
    border-bottom: 0.12rem solid var(--primary);

    transition: 0.2s;
}

.selected-option div.active::after {
    transform: translateY(-50%) rotate(225deg);
}

.select-box .options {
    position: absolute;
    top: 4rem;

    width: 100%;
    background-color: #fff;
    border-radius: 0.5rem;

    display: none;
}

.select-box .options.active {
    display: block;
}

.select-box .options::before {
    position: absolute;
    content: '';
    left: 1rem;
    top: -1.2rem;

    width: 0;
    height: 0;
    border: 0.6rem solid transparent;
    border-bottom-color: var(--primary);
}

input.search-box {
    background-color: var(--primary);
    color: #fff;
    border-radius: 0.5rem 0.5rem 0 0;
    padding: 1.4rem 1rem;
}

.select-box ol {
    list-style: none;
    max-height: 23rem;
    overflow: overlay;
}

.select-box ol::-webkit-scrollbar {
    width: 0.6rem;
}

.select-box ol::-webkit-scrollbar-thumb {
    width: 0.4rem;
    height: 3rem;
    background-color: #ccc;
    border-radius: 0.4rem;
}

.select-box ol li {
    padding: 1rem;
    display: flex;
    justify-content: space-between;
    cursor: pointer;
}

.select-box ol li.hide {
    display: none;
}

.select-box ol li:not(:last-child) {
    border-bottom: 0.1rem solid #eee;
}

.select-box ol li:hover {
    background-color: lightcyan;
}

.select-box ol li .country-name {
    margin-left: 0.4rem;
}
.divider {
    display: flex;
    align-items: center;
    text-align: center;
}
.divider::before,
.divider::after {
    content: '';
    flex: 1;
    border-bottom: 1px solid #ccc;
}

.divider-text {
    padding: 0 10px;
    color: #666; /* Text color */
    font-size: 16px;
}

.login-with-google-btn {
    transition: background-color 0.3s, box-shadow 0.3s;
    width: 98%;
    height: 57px;
    padding: 12px 16px 12px 42px;
    border: none;
    border-radius: 15px;

    box-shadow: 0 -1px 0 rgba(0, 0, 0, 0.04), 0 1px 1px rgba(0, 0, 0, 0.25);

    margin-top: 20px;

    color: black;
    font-size: 14px;
    font-weight: 500;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;

    background-image: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTgiIGhlaWdodD0iMTgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGcgZmlsbD0ibm9uZSIgZmlsbC1ydWxlPSJldmVub2RkIj48cGF0aCBkPSJNMTcuNiA5LjJsLS4xLTEuOEg5djMuNGg0LjhDMTMuNiAxMiAxMyAxMyAxMiAxMy42djIuMmgzYTguOCA4LjggMCAwIDAgMi42LTYuNnoiIGZpbGw9IiM0Mjg1RjQiIGZpbGwtcnVsZT0ibm9uemVybyIvPjxwYXRoIGQ9Ik05IDE4YzIuNCAwIDQuNS0uOCA2LTIuMmwtMy0yLjJhNS40IDUuNCAwIDAgMS04LTIuOUgxVjEzYTkgOSAwIDAgMCA4IDV6IiBmaWxsPSIjMzRBODUzIiBmaWxsLXJ1bGU9Im5vbnplcm8iLz48cGF0aCBkPSJNNCAxMC43YTUuNCA1LjQgMCAwIDEgMC0zLjRWNUgxYTkgOSAwIDAgMCAwIDhsMy0yLjN6IiBmaWxsPSIjRkJCQzA1IiBmaWxsLXJ1bGU9Im5vbnplcm8iLz48cGF0aCBkPSJNOSAzLjZjMS4zIDAgMi41LjQgMy40IDEuM0wxNSAyLjNBOSA5IDAgMCAwIDEgNWwzIDIuNGE1LjQgNS40IDAgMCAxIDUtMy43eiIgZmlsbD0iI0VBNDMzNSIgZmlsbC1ydWxlPSJub256ZXJvIi8+PHBhdGggZD0iTTAgMGgxOHYxOEgweiIvPjwvZz48L3N2Zz4=);
    background-size: 24px 24px;
    background-color: white;
    background-repeat: no-repeat;
    background-position: 16px 17px;
}

body {
    text-align: center;
    padding-top: 2rem;
}
</style>
