import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({request, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    let header;
    let token = cookies.get("session")
    if (token) {
        header = {cookie: "session=" + token}
    } else {
        return {
            me: {}
        }
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

    let [/*stats,*/ me] = await Promise.all(fetches)

    let data = {
        me: me
    };

    if (env.LBF_DEV_FRONTEND_ONLY) {
        data.point = true;
    }

    return data;
};