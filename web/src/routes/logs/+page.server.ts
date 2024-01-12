import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({request, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    const token = cookies.get("sessions")
    let header;
    if (token) {
        header = {cookie: "sessions=" + token}
    }

    const fetches = []
    for (const path of [
        /*"/api/logs/stats",*/ "/api/users/@me"
    ]) {
        fetches.push(fetch(url + path, {
            headers: header,
        }).then(resp => {
            return resp.ok ? resp.json() : {}
        }))
    }

    const [/*stats,*/ me] = await Promise.all(fetches)

    return {
        me: me,
        /*stats: stats,*/
    };
};