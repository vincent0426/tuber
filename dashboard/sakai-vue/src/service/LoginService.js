import request from './wrapper';

export class LoginService {
    postLogin(id_token) {
        return request({
            method: 'post',
            url: '/auth/login',
            data: {},
            headers: {
                'id_token': id_token
            }
        });
    }

    delLogin() {
        return request({
            method: 'get',
            url: 'demo/data/logout.json'
        });
    }
}
