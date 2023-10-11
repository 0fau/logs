import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({params, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS

    if (params.user.charAt(0) !== "@") {
        return {}
    }

    const resp = await fetch(
        url + "/api/users/" + params.user.slice(1),
    )
    if (!resp.ok) {
        return {}
    }
    const user = await resp.json()

    const recent = await fetch(
        url + "/api/logs/recent?user=" + user.id,
    );
    if (!recent.ok) {
        return {
            user: user,
        }
    }

    return {
        user: user,
        recent: await recent.json(),
    };
};