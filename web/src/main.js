import { createApp } from 'vue';
import { createWebHashHistory, createRouter } from 'vue-router'
import Antd from 'ant-design-vue';
import App from './App';
import 'ant-design-vue/dist/reset.css';

import LoginView from './components/LoginView.vue';
import MainPage from './MainPage.vue'

const routes = [
    { name: '/', path: '/', component: LoginView },
    { name: 'login', path: '/login', component: LoginView },
    { name: 'main', path: '/main', component: MainPage },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
    base: process.env.VUE_APP_BASE_URL,
})
router.beforeEach((to, from, next) => {
    // 假设我们要检查名为'user'的cookie
    const hasCookie = document.cookie.split(';').some(cookie => cookie.trim().startsWith('token='));
    console.log(hasCookie,to.path);
    if (to.path !== '/login' && !hasCookie) {
        next('/login');
    } else {
        next();
    }
});

const app = createApp(App);

app.use(Antd).use(router).mount('#app');