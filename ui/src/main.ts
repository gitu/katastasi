import './assets/main.css'

import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from './router'

import CheckCircleOutline from 'vue3-material-design-icons-ts/dist/CheckCircleOutline.vue';
import InformationOutline from "vue3-material-design-icons-ts/dist/InformationOutline.vue";
import AlertOutline from "vue3-material-design-icons-ts/dist/AlertOutline.vue";
import AlertOctagram from "vue3-material-design-icons-ts/dist/AlertOctagramOutline.vue";
import HelpCircleOutline from "vue3-material-design-icons-ts/dist/HelpCircleOutline.vue";
import ArrowExpand from "vue3-material-design-icons-ts/dist/ArrowExpand.vue";


const app = createApp(App)

app.component('CheckCircleOutline', CheckCircleOutline)
app.component('InformationOutline', InformationOutline)
app.component('AlertOutline', AlertOutline)
app.component('AlertOctagram', AlertOctagram)
app.component('HelpCircleOutline', HelpCircleOutline)
app.component('ArrowExpand', ArrowExpand)


app.use(createPinia())
app.use(router)
app.mount('#app')
