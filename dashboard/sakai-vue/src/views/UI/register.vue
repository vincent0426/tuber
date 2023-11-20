<template>
    <p>Upload your Driving License</p>
    <div class="upload-container">
        <div class="drop-zone" @click="openFilePicker('fileInput1')" @dragover.prevent @drop="onFileDrop('file1', $event)">
          <p> {{ file1.p }}</p>
          <input type="file" @change="onFileChange('file1', $event)" ref="fileInput1" style="display: none;" />
          <div v-if="file1.imageUrl" class="image-preview">
            <img :src="file1.imageUrl" alt="Uploaded Image" />
          </div>
        </div>
        
        <div class="drop-zone" @click="openFilePicker('fileInput2')" @dragover.prevent @drop="onFileDrop('file2', $event)">
          <p>{{  file2.p }}</p>
          <input type="file" @change="onFileChange('file2', $event)" ref="fileInput2" style="display: none;" />
          <div v-if="file2.imageUrl" class="image-preview">
            <img :src="file2.imageUrl" alt="Uploaded Image" />
          </div>
        </div>
    </div>
    <p>Enter your license plate number</p>
    <div class="upload-container">
    <input type="text" v-model="plateNumber" @input="handleInput" class="full-width-input"/>
    </div>
    <p>Upload your car with license plate number</p>
    <div class="drop-zone" @click="openFilePicker('fileInput3')" @dragover.prevent @drop="onFileDrop('file3', $event)">
        <p>{{ file3.p }}</p>
        <input type="file" @change="onFileChange('file3', $event)" ref="fileInput3" style="display: none;" />
        <div v-if="file3.imageUrl" class="image-preview">
            <img :src="file3.imageUrl" alt="Uploaded Image" />
        </div>
    </div>
    <div>
        <button class="arrow-button black-button"><span>Next</span> <span class="arrow">→</span></button>
        <button class="arrow-button white-button">Skip (As a passenger) <span class="arrow">→</span></button>

    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        file1: { imageUrl: null , p:"Upload Driving License (front)"},
        file2: { imageUrl: null , p:"Upload Driving License (back)"},
        plateNumber: null,
        file3: { imageUrl: null , p:"Upload car with plate number"},
      };
    },
    methods: {
      onFileChange(fileKey, e) {
        const file = e.target.files[0];
        this.loadImage(fileKey, file);
      },
      onFileDrop(fileKey, e) {
        const file = e.dataTransfer.files[0];
        this.loadImage(fileKey, file);
      },
      loadImage(fileKey, file) {
        this[fileKey].imageUrl = URL.createObjectURL(file);
        this[fileKey].p = "";
      },
      openFilePicker(inputRef) {
        this.$refs[inputRef].click();
      },
      handleInput(){
        // TODO
      },
    }
  };
  </script>
  
  <style>
  .upload-container {
    text-align: center;
    padding: 5px;
  }
  .full-width-input{
    width:98%;
    height: 32px;
    box-sizing: border-box;
  }
  
  .drop-zone {
    border: 2px dashed #ccc;
    padding: 40px;
    cursor: pointer;
    margin: 10px;
    margin-bottom: 20px;
  }
  
  .drop-zone p {
    margin: 10px 0;
    font-size: 16px;
    color: #555;
  }
  
  .image-preview img {
    max-width: 100%;
    height: auto;
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
    margin-bottom: 14;
    }
    .black-button{
        background-color: black;
        color: white;
    }
    .white-button{
        background-color: #EEE;
        color: black;
    }
    .arrow-button .arrow {
        position: absolute;
        right: 10px;
    }
  </style>
  