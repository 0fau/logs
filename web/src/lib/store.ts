import {writable} from 'svelte/store';
import {browser} from "$app/environment"
import {getContext, setContext} from "svelte";

let stored;
if (browser) {
    if (localStorage.getItem("settings")) {
        stored = JSON.parse(localStorage.settings);
    } else {
        stored = {
            logs: {
                announcement: true,
                scope: "Arkesia"
            }
        };
    }
}

export function getUser() {
    let store = getContext("user");
    if (!store) {
        store = writable({});
        setContext("user", store);
    }
    return store;
}

export function getSettings() {
    let settings = getContext("settings");
    if (!settings) {
        settings = writable(stored);
        setContext("settings", settings);

        if (browser) {
            settings.subscribe((value) => {
                localStorage.setItem("settings", JSON.stringify(value))
            })
        }
    }
    return settings;
}