import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({params, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    const headers: HeadersInit = {}
    if (cookies.get("sessions")) {
        headers.cookie = "sessions=" + cookies.get("sessions")
    }

    const enc = await fetch(
        url + "/api/log/" + params.log,
        {headers: headers},
    );
    if (!enc.ok) {
        return {};
    }

    return {
        encounter: await enc.json(),
    };
};