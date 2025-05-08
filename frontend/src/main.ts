import './style.css'
import App from './App.svelte'
import {createDevicesStore} from "./lib/stores/devices";

createDevicesStore();

const app = new App({
    target: document.getElementById('app')
})

export default app
