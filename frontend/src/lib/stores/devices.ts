import {writable} from 'svelte/store';
import {EventsOn} from '../../../wailsjs/runtime';
import type {ha} from '../../../wailsjs/go/models';
import {GetDevices} from '../../../wailsjs/go/main/App';
import {tryCatch} from "../try-catch";

export const devices = writable<ha.Device[]>([]);
export const loading = writable<boolean>(true);
export const error = writable<string | null>(null);

function log(message?: any, ...optionalParams: any[]) {
    console.log("[devices] " + message, ...optionalParams);
}

export function createDevicesStore() {
    loadDevices();

    EventsOn('devices:list', (devicesList: ha.Device[]) => {
        log("received devices:list event:", devicesList);

        devices.set(devicesList);
        error.set(null);
        loading.set(false);
    });

    EventsOn('devices:error', (errorMessage: string) => {
        log("received devices:error event:", errorMessage);

        error.set(errorMessage);
        loading.set(false);
    });
}

export async function loadDevices() {
    loading.set(true);

    log("loading devices ...");

    const {data, err} = await tryCatch(GetDevices())
    if (err != null) {
        log("failed to load devices:", err);

        error.set(err.message);
        loading.set(false);
        return;
    }

    log("loaded devices:", data);

    devices.set(data);
    loading.set(false);
}