<template>
    <p class="subtitle">Enter your mobile number</p>
    
    <div class="select-box">
        <div class="selected-option">
            <div>
                <span class="iconify" data-icon="flag:gb-4x3"></span>
                <strong>+1</strong>
            </div>
            <input type="tel" name="tel" placeholder=" Number"  v-model="phone">
        </div>
        <div class="options">
            <input type="text" class="search-box" placeholder="Search Country Name">
            <ol>

            </ol>
        </div>
    </div>


    <button class="arrow-button black-button" @click="goNext"><span>Next</span> <span class="arrow">→</span></button>
    <p class="p-message">By continuing you may receive an SMS for verification. Message and data rates may apply</p>
  
    <div class="divider">
      <span class="divider-text">or</span>
    </div>
        
    <!--
    <button type="button" class="login-with-google-btn" @click="callback" >
        Sign in with Google
    </button>
    -->
    <GoogleLogin :callback="callback" prompt  />


  </template>

  <script>
import { googleSdkLoaded , decodeCredential} from "vue3-google-login";
import axios from "axios";
  

  export default {
    data() {
      return {
            phone: null,
            userDetails: null,
            callback: (response) => {
                console.log(response)
                console.log(decodeCredential(response.credential))
            }
        }
    },
    methods: {
        login() {
            googleSdkLoaded(google => {
                google.accounts.oauth2
                .initCodeClient({
                    client_id:
                    "127171133807-4qm2bj37o4tkk7h868j5kmptgq4l878e.apps.googleusercontent.com",
                    scope: "email profile openid",
                    redirect_uri: "http://localhost:4000/auth/callback",
                    callback: "callback",
                })
                .requestCode();
            });
        },
        async sendCodeToBackend(code) {
            try {
                const headers = {
                Authorization: code
                };
                const response = await axios.post("http://localhost:4000/auth", null, { headers });
                const userDetails = response.data;
                console.log("User Details:", userDetails);
                this.userDetails = userDetails;

                // Redirect to the homepage ("/")
                //this.$router.push("/rex");
            } catch (error) {
                console.error("Failed to send authorization code:", error);
            }
            }
    },
    name: 'GoogleSignIn',
    mounted() {
        gapi.load('auth2', function() {
        gapi.auth2.init({
            client_id: '127171133807-4qm2bj37o4tkk7h868j5kmptgq4l878e.apps.googleusercontent.com',
        });
        });
  }
  };
  </script>
  
  <style>

  p{
    text-align: left;
  }
  .subtitle{
    font-size: 16px;
    color: black;
    font-weight: bold;
  }
  .full-width-input{
    width:98%;
    height: 32px;
    box-sizing: border-box;
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
    width: 100%; /* Adjust the width as needed */
    margin: 10px;
    margin-top: 20px;
    margin-bottom: 20px;
    }
    .black-button{
        background-color: black;
        color: white;
    }

    .arrow-button .arrow {
        position: absolute;
        right: 10px;
    }
    .p-message{
        color:#757575;
        font-size:14px;
    }
    
    .select-box {
    position: relative;
}

.select-box input {
    background-color: #eee;
    width: 100%;
    padding: 1rem .6rem;
    font-size: 1.1rem;
    
    border: .1rem solid transparent;
    outline: none;
}

input[type="tel"] {
    border-radius: 0 .5rem .5rem 0;
}

.select-box input:focus {
    border: .1rem solid var(--primary);
}
.selected-option input{
    width: 80%;
}
.selected-option {
    border-radius: .5rem;
    overflow: hidden;

    display: flex;
    justify-content: space-between;
    align-items: center;
}

.selected-option div{
    position: relative;

    width: 6rem;
    padding: 0 2.8rem 0 .5rem;
    text-align: center;
    cursor: pointer;
}

.selected-option div::after{
    position: absolute;
    content: "";
    right: .8rem;
    top: 50%;
    transform: translateY(-50%) rotate(45deg);
    
    width: .8rem;
    height: .8rem;
    border-right: .12rem solid var(--primary);
    border-bottom: .12rem solid var(--primary);

    transition: .2s;
}

.selected-option div.active::after{
    transform: translateY(-50%) rotate(225deg);
}

.select-box .options {
    position: absolute;
    top: 4rem;
    
    width: 100%;
    background-color: #fff;
    border-radius: .5rem;

    display: none;
}

.select-box .options.active {
    display: block;
}

.select-box .options::before {
    position: absolute;
    content: "";
    left: 1rem;
    top: -1.2rem;

    width: 0;
    height: 0;
    border: .6rem solid transparent;
    border-bottom-color: var(--primary);
}

input.search-box {
    background-color: var(--primary);
    color: #fff;
    border-radius: .5rem .5rem 0 0;
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
    border-radius: .4rem;
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
    border-bottom: .1rem solid #eee;
}

.select-box ol li:hover {
    background-color: lightcyan;
}

.select-box ol li .country-name {
    margin-left: .4rem;
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
  transition: background-color .3s, box-shadow .3s;
    width: 98%;
    height: 57px;
  padding: 12px 16px 12px 42px;
  border: none;
  border-radius: 15px;

  box-shadow: 0 -1px 0 rgba(0, 0, 0, .04), 0 1px 1px rgba(0, 0, 0, .25);
  
    margin-top: 20px;

  color: black;
  font-size: 14px;
  font-weight: 500;
  font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen,Ubuntu,Cantarell,"Fira Sans","Droid Sans","Helvetica Neue",sans-serif;
  
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