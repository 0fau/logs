import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({params, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    const headers: HeadersInit = {}
    if (cookies.get("session")) {
        headers.cookie = "session=" + cookies.get("session")
    }

    const fetches = []
    for (const path of [
        "/api/log/" + params.log, "/api/users/@me"
    ]) {
        fetches.push(fetch(url + path, {
            headers: headers,
        }).then(resp => {
            return resp.ok ? resp.json() : null
        }))
    }

    const [encounter, me] = await Promise.all(fetches)

    return {
        encounter: encounter,
        me: me,
    };
};