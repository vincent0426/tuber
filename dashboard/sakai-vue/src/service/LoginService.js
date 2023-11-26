import request from "./wrapper";

export class LoginService {
    postLogin(username, password) {
        return request({
            method: "post",
            url: "demo/data/login.json",
            data: {
                username,
                password,
            },
        });
    }

    checkLogin() {
        return request({
            method: "get",
            url: "demo/data/whoami.json",
        });
    }

    delLogin() {
        return request({
            method: "get", 
            url: "demo/data/logout.json",
        });
    }
}