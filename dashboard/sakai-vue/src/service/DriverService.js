import request from './wrapper';

export class DriverService {
    getDriver(id) {
        return request({
            method: 'get',
            url: `/drivers/${id}`,
            data: {}
        });
    }
    getFavorite() {
        return request({
            method: 'get',
            url: '/favorite-drivers',
            data: {}
        });
    }
    postFavorite(id) {
        return request({
            method: 'post',
            url: `/favorite-drivers/${id}`,
            data: {}
        });
    }
    postBecomeDriver(formData) {
        return request({
            method: 'post',
            url: '/drivers',
            data: formData
        });
    }
    // getFavorite(id_token) {
    //     return request({
    //         method: 'post',
    //         url: '/auth/login',
    //         data: {},
    //         headers: {
    //             id_token: id_token
    //         }
    //     });
    // }
}
