import LoginView from './pages/LoginPage.vue';
import MainPage from './pages/MainPage.vue'
import UserPage from './pages/UserPage.vue'
import WAFPage from './pages/WAFPage.vue'

export const routes = [
    { name: '/', path: '/', component: LoginView },
    { name: 'login', path: '/login', component: LoginView },
    { name: 'main', path: '/main', component: MainPage },
    { name: 'user', path: '/user', component: UserPage },
    { name: 'waf', path: '/waf', component: WAFPage },
]