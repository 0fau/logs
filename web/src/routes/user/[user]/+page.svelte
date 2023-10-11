<script lang="ts">
    import type {PageData} from './$types';
    import {formatDamage, formatDate, formatDuration} from '$lib/components/meter/print'

    export let data: PageData;
    console.log(data)

    let more = (data.recent ?? []).length == 5

    async function expand() {
        if (data.recent.length === 0) {
            return
        }

        let past = data.recent[data.recent.length - 1].date

        let url = location.protocol + '//' + location.host;

        const recent = await (await fetch(
            url + "/api/logs/recent?user" + data.user.id + "&past=" + past
        )).json()

        if (recent.length < 5) {
            more = false
        }

        data.recent = data.recent.concat(recent)
    }
</script>

<svelte:head>
    {#if data.user && data.user.username}
        <title>@{data.user.username}</title>
    {:else}
        <title>HUH</title>
    {/if}
</svelte:head>

<div class="mx-3 my-14 mx-auto text-center max-w-fit">
    {#if data.user && data.user.username}
        <p class="font-semibold text-pink-600">@{data.user.username}</p>
        <br/>
        <p>[<span class="font-semibold">Recent Logs</span>]</p>
        {#each data.recent as log}
            <a href="/log/{log.id}">
                [#{log.id}]
                <span class="underline">Raid: {log.raid}</span> {formatDate(log.date)}
            </a>
            <p>{formatDamage(log.damage)} damage dealt in {formatDuration(log.duration)}</p>
            <p>---</p>
        {/each}
        {#if more}
            <button class="underline" on:click={expand}>show more</button>
        {:else}
            end
        {/if}
    {:else}
        uh oh user not found :p
    {/if}
</div>