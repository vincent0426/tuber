import { createStore } from 'vuex';
import { LoginService } from '@/service';

/**
 * Centralized state management for the application.
 */

const loginService = new LoginService();

// 尝试从 localStorage 中恢复状态
const initialState = () => {
    const savedState = localStorage.getItem('vuex-state');
    const savedStateExpiry = localStorage.getItem('vuex-state-expiry');
    const now = new Date();
    if (savedState && (now.getTime() < savedStateExpiry)){
        return JSON.parse(savedState);
    }
    return {
        user: {},
        role: '',
        login: false,
        expiry: 0
    };
};

export default createStore({
    state: initialState(),
    getters: {
        user: (state) => state.user,
        role: (state) => state.role,
        login: (state) => state.login,
        expiry: (state) => state.expiry
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
        saveState(state) {
            const now = new Date();
            localStorage.setItem('vuex-state', JSON.stringify(state));
            localStorage.setItem('vuex-state-expiry', now.getTime() + 3600000);
        }
    },
    actions: {
        setLogin({ commit }, login) {
            commit('setLogin', login);
        },
        setUser({ commit }, user) {
            commit('setUser', user);
        },
        setRole({ commit }, role) {
            commit('setRole', role);
        },
        saveState({ commit }, state) {
            commit('saveState', state);
        },

        async login({ dispatch }, {id_token}) {
            try {
                const user = await loginService.postLogin(id_token);
                const now = new Date();
                console.log(user);
                dispatch('setUser', user);
                dispatch('setRole', "passenger");
                dispatch('setLogin', true);
                dispatch('saveState',this.state);
            } catch (e) {
                dispatch('setLogin', false);
                throw 'Login failed';
            }
        },

        async logout({ dispatch }) {
            try {
                await loginService.delLogin();
                dispatch('setUser', null);
                dispatch('setLogin', false);
                dispatch('saveState',this.state);
            } catch (e) {
                dispatch('setUser', null), dispatch('setLogin', false);
            }
        },

        async swapRole({ dispatch }) {
            try {
                if (this.getters.role === 'passenger') {
                    dispatch('setRole', 'driver');
                } else {
                    dispatch('setRole', 'passenger');
                }
                dispatch('saveState',this.state);
            } catch (e) {
                throw 'Role Change Error';
            }
        }
    }
});
