import type {LayoutServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: LayoutServerLoad = async ({request, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    const sessions = cookies.get("sessions")
    if (!sessions) {
        return {
            me: {},
        }
    }

    const me = await fetch(
        url + "/api/users/@me",
        {headers: {cookie: "sessions=" + cookies.get("sessions")}},
    )

    return {
        me: await me.json(),
    };
};