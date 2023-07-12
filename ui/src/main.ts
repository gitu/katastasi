import './assets/main.css'

import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from './router'
import moment from 'moment';

const app = createApp(App)

app.config.globalProperties.$filters = {
    timeAgo(date: moment.MomentInput) {
        return moment(date).fromNow()
    },
}
app.use(createPinia())
app.use(router)
app.mount('#app')
