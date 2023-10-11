import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({request, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    const recent = await fetch(
        url + "/api/logs/recent",
    )

    const sessions = cookies.get("sessions")
    if (!sessions) {
        return {
            me: {},
            recent: await recent.json(),
        }
    }

    const me = await fetch(
        url + "/api/users/@me",
        {headers: {cookie: "sessions=" + cookies.get("sessions")}},
    )

    return {
        me: await me.json(),
        recent: await recent.json(),
    };
};