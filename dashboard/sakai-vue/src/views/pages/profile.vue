<script setup>
import store from '@/store';
import { ref } from 'vue';

</script>

<template>

{{ getUserName }}

<div class="center-container">
    <div class="photo-uploader">
      <div class="photo-container">
        <img :src="imgURL" alt="Profile" class="profile-photo">
        <img for="photo-input" class="camera-icon" src="../../image/camera.png"/>
        <input type="file" id="photo-input" @change="onFileChange" accept="image/*" hidden>
      </div>
    </div>
  </div>
    <div class="field">
        <p>Name</p>
        <p class="content">{{ user_name }}</p>
    </div>

    <div class="field">
        <p>Email</p>
        <p class="content">{{ user_email }}</p>
    </div>

    <div class="field switch-field">
        <label class="switch">
        <input type="checkbox">
        <span class="slider round"></span>
        </label>
        <span class="content">Recieving E-mail Alert</span>
    </div>

   
</template>

<script>
  export default {

    data(){
        return{
            user: ref(null),
            user_id: ref(""),
            user_name: ref(""),
            user_email: ref(""),
            imgURL: ref("../../image/Ellipse2.png"),
        }
    },
    methods: {
        goBack(){
            console.log("Clicked")
        },

        onFileChange(event) {
            const file = event.target.files[0];
            this.photoUrl = URL.createObjectURL(file);
        },
    },
    computed:{
        getUserName(){
            this.user =  this.$store.getters.user
            console.log(this.user)

            this.user_id = this.user.id;
            this.user_email = this.user.email;
            this.user_name = this.user.name;
            this.imgURL = this.user.imageURL;
        }
    }
  };
</script>

<style lang="scss" scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}
p{
    margin: 0;
    font-weight: bold;
}
p.content{
    font-size: 20px;
    font-weight:500;
}
span.content{
    font-size: 22px;
    text-align: center;
    margin-left: 10px;
}
.left-arrow {
  /* Additional styling for arrows */
  cursor: pointer; /* If they are clickable */
  font-size: 24px;
}

.title {
  text-align: center;
  flex-grow: 1; /* Ensures title remains centered */
  font-size: 24px;
  color: black;
}

.center-container {
  display: flex;
  justify-content: center;
  align-items: center;
}

.photo-uploader {
  position: relative;
  width: 150px; /* Adjust size as needed */
  height: 150px;
}

.photo-container {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background-color: #f0f0f0; /* Default background color */
}

.profile-photo {
  position: relative;
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.camera-icon {
  position: absolute;
  right: -10px; /* Adjust position as needed */
  bottom: -10px;
  cursor: pointer;
  background-color: #fff;
  border-radius: 50%;
  padding: 8px;
}

.full-width-input{
    width:98%;
    height: 44px;
    box-sizing: border-box;
  }

.full-width-area{
    height: 110px;
}

div.field{
    margin-top:10px;
    margin-bottom: 10px;
    align-items: center;
}
div.switch-field{
    display: flex;
}
div.text-area{
    margin-top: 55px;
}

.switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

.switch input { 
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  -webkit-transition: .4s;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  -webkit-transition: .4s;
  transition: .4s;
}

input:checked + .slider {
  background-color: #2196F3;
}

input:focus + .slider {
  box-shadow: 0 0 1px #2196F3;
}

input:checked + .slider:before {
  -webkit-transform: translateX(26px);
  -ms-transform: translateX(26px);
  transform: translateX(26px);
}

/* Rounded sliders */
.slider.round {
  border-radius: 34px;
}

.slider.round:before {
  border-radius: 50%;
}

</style>