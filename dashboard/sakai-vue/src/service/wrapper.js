import axios from 'axios';

const client = axios.create({
    baseURL: import.meta.env.VITE_BASE_URL,
    withCredentials: true
});

/**
 * axios api wrapper, for success and error handler
 * @param {*} options - options passed to axios
 */
const request = async function (options) {
    if (process.env.NODE_ENV !== 'production') {
        console.debug('Request Option', options);
    }
    const onSuccess = function (response) {
        if (process.env.NODE_ENV !== 'production') {
            console.debug('Request Successful!', response);
        }
        if (options.responseType === 'blob') {
            return response;
        }
        return response.data;
    };

    const onError = function (error) {
        console.error('Request Failed:', error.config);

        if (error.response) {
            // Request was made but server responded with something
            // other than 2xx
            console.error('Status:', error.response.status);
            console.error('Data:', error.response.data);
            console.error('Headers:', error.response.headers);
        } else {
            // Something else happened while setting up the request
            // triggered the error
            console.error('Error Message:', error.message);
        }

        return Promise.reject(error.response?.data || error.message);
    };

    return client(options).then(onSuccess).catch(onError);
};

export default request;
