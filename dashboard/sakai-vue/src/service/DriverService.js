import request from './wrapper';

export class DriverService {
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
