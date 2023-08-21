import type {PageServerLoad} from './$types';

export const load: PageServerLoad = async ({cookies}) => {
    const me = await fetch(
        "http://localhost:3000/api/users/@me",
        {headers: { cookie: "sessions=" + cookies.get("sessions")}},
    )
    return await me.json();
};