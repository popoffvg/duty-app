// import { createApp } from 'vue'
// import App from './App.vue'
// import router from './router'
//
// createApp(App).use(router).mount('#app')

import * as Vue from 'vue' // in Vue 3
import App from './App.vue'
import BootstrapVue3 from 'bootstrap-vue-3'
//import IconsPlugin from 'bootstrap-vue-3'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueBasicAlert from 'vue-basic-alert'



// Import Bootstrap an BootstrapVue CSS files (order is important)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue-3/dist/bootstrap-vue-3.css'

axios.defaults.baseURL = '//localhost:8080/';
const app = Vue.createApp(App)

app.use(router)
app.use(VueAxios, axios)
// Make BootstrapVue available throughout your project
app.use(BootstrapVue3)
// Optionally install the BootstrapVue icon components plugin
//app.use(IconsPlugin)
app.use(VueBasicAlert)



app.mount('#app')

require('./assets/css/main.css');
