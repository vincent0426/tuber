<template>
<div class="header">
    <span class="left-arrow" @click="goBack">←</span> <!-- Left arrow; replace with an icon if necessary -->
    <p class="title">Edit Profile</p> <!-- Title -->
</div>

<div class="center-container">
    <div class="photo-uploader">
      <div class="photo-container">
        <img src="../../image/Ellipse2.png" alt="Profile" class="profile-photo">
        <img for="photo-input" class="camera-icon" src="../../image/camera.png"/>
        <input type="file" id="photo-input" @change="onFileChange" accept="image/*" hidden>
      </div>
    </div>
  </div>

    <div class="field">
  <p>Name</p>
  <input class="full-width-input"/>
    </div>

    <div class="field">
  <p>Email</p>
  <input class="full-width-input"/>
    </div>

    <div class="field">
  <p>Password</p>
  <input class="full-width-input" type="password"/>
    </div>

    <div class="field">
  <p>Data of Birth</p>
  <input class="full-width-input"/>
    </div>

    <div class="text-area">
    <p>Personal profile</p>
  <textarea class="full-width-input full-width-area"/>
    </div>
</template>

<script>
  export default {

    data(){
        return{

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
  position: absolute;
  width: 100%;
  height: 100%;
  object-fit: cover;
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
}
div.text-area{
    margin-top: 55px;
}
</style>