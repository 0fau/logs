import type {PageServerLoad} from './$types';
import {env} from '$env/dynamic/private';

export const load: PageServerLoad = async ({params, cookies}) => {
    let url = "http"
    if (env.LBF_API_SERVER_ADDRESS_SECURE == "true") {
        url += "s"
    }
    url += "://" + env.LBF_API_SERVER_ADDRESS


    const enc = await fetch(
        url + "/api/log/" + params.log + "/short",
    );
    if (!enc.ok) {
        return {};
    }

    return {
        encounter: await enc.json(),
    };
};