<script>
import request from '../../../service/wrapper';

export default {
    data() {
        return {
            formData: {
                license: '',
                brand: '',
                model: '',
                color: '',
                plate: ''
            }
        };
    },
    methods: {
        async postBecomeDriver() {
            try {
                const response = await request({
                    method: 'post',
                    url: 'http://localhost:3000/v1/drivers', // Ensure this is the correct URL
                    data: this.formData
                });
                console.log('Success:', response.data);
                // redirect to driver home page
                this.$router.push('/driver/home');
                // Handle success (e.g., show a success message, redirect, etc.)
            } catch (error) {
                console.error('Error:', error.response ? error.response.data : error.message);
                // Handle error (e.g., show an error message)
            }
        },
        submitForm() {
            this.postBecomeDriver();
        }
    }
};
</script>

<template>
    <div class="grid p-fluid">
        <div class="col-12">
            <h3>Become a Driver</h3>
            <div class="card">
                <h5>Licence Number</h5>
                <InputText v-model="formData.license" placeholder="123456ABC" id="License" type="text" />
                <h5>Car Brand</h5>
                <InputText v-model="formData.brand" placeholder="Benz" id="Brand" type="text" />
                <h5>Car Model</h5>
                <InputText v-model="formData.model" placeholder="C300" id="Model" type="text" />
                <h5>Car Color</h5>
                <InputText v-model="formData.color" placeholder="White" id="Color" type="text" />
                <h5>License Plate Number</h5>
                <InputText v-model="formData.plate" placeholder="XYZ-1234" id="Plate" type="text" />
                <Button label="Submit" class="mr-2 mb-2 mt-4" id="Submit" @click="submitForm"></Button>
            </div>
        </div>
    </div>
</template>
