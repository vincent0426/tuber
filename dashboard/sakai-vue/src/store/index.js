import { createStore } from "vuex";
import { LoginService } from "@/service";

/**
 * Centralized state management for the application.
 */

const loginService = new LoginService();


export default createStore({
  state: {
    user: {},
    role: 'customer',
    login: false,
  },
  getters: {
    user: (state) => state.user,
    role: (state) => state.role,
    login: (state) => state.login,
  },
  mutations: {
    setRole(state, role) {
      state.role = role;
    },
    setUser(state, user) {
      state.user = user;
    },
    setLogin(state, login) {
      state.login = login;
    },
  },
  actions: {
    setLogin({ commit }, login) {
      commit("setLogin", login);
    },
    setUser({ commit }, user) {
      commit("setUser", user);
    },
    setRole({ commit }, role) {
      commit("setRole", role);
    },

    async checkLogin({ dispatch }) {
      try {
        const user = await loginService.checkLogin();
        dispatch("setUser", user);
        dispatch("setLogin", true);
      } catch (e) {
        dispatch("setLogin", false);
        dispatch("setUser", null);
        throw "Not logged in";
      }
    },

    async login({ dispatch }, { username, password }) {
      try {
        const user = await loginService.postLogin(username, password);
        dispatch("setUser", user);
        dispatch("setRole", 'customer');
        dispatch("setLogin", true);
      } catch (e) {
        dispatch("setLogin", false);
        throw "Login failed";
      }
    },

    async logout({ dispatch }) {
      try {
        await loginService.delLogin();
        dispatch("setUser", null);
        dispatch("setLogin", false);
      } catch (e) {
        dispatch("setUser", null),
        dispatch("setLogin", false);
      }
    },
  },
});

