import {writable} from 'svelte/store';
import {browser} from "$app/environment"

export const user = writable({});

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

export const settings = writable(stored);

if (browser) {
    settings.subscribe((value) => {
        localStorage.setItem("settings", JSON.stringify(value))
    })
}
