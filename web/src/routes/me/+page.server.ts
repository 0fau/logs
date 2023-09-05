import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    const sessions = cookies.get("sessions")
    if (!sessions) {
        return {}
    }

    const me = await fetch(
        url + "/api/users/@me",
        {headers: {cookie: "sessions=" + cookies.get("sessions")}},
    )
    return await me.json();
};