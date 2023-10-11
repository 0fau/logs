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
        headers.cookies = "sessions=" + cookies.get("sessions")
    }

    const [enc, details] = await Promise.all([
        fetch(
            url + "/api/logs/" + params.log,
            {headers: headers},
        ),
        fetch(
            url + "/api/logs/" + params.log + "/details",
            {headers: headers},
        )
    ]);

    if (!enc.ok || !details.ok) {
        return {};
    }

    return {
        encounter: await enc.json(),
        details: await details.json(),
    };
};