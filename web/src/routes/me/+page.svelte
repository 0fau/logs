<script lang="ts">
    import type {PageData} from './$types';
    import {formatRelative, formatDistance} from 'date-fns';
    import numeral from 'numeral';

    export let data: PageData;

    function formatDate(date: number) : string {
        return formatDistance(new Date(date), new Date(), { addSuffix: true })
    }

    function formatDuration(duration: number) : string {
        let date = new Date(duration);
        return date.getMinutes() + 'm' + date.getSeconds() + 's';
    }

    function formatDamage(damage: number) : string {
        return numeral(damage).format('0.0a')
    }
</script>

<div class="container flex h-screen justify-center items-center">
    <div class="columns-3xs text-center ...">
        {#if data.me.username}
            <p><span class="font-semibold">User</span> {data.me.username}</p>
            <form action="/logout" method="post">
                <button class="underline">Logout</button>
            </form>
        {:else}
            <form action="/oauth2" method="post">
                <button class="underline">Login</button>
            </form>
        {/if}
        {#if data.recent}
            <div class="m-5">
                <p>[<span class="font-semibold">Recent Logs</span>]</p>
                {#each data.recent as encounter}
                    <div>
                        <p><span class="underline">Raid: {encounter.raid}</span> {formatDate(encounter.date)}</p>
                        <p>{formatDamage(encounter.damage)} damage dealt in {formatDuration(encounter.duration)}</p>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<style lang="postcss">
    :global(html) {
        background-color: theme(colors.gray.100);
    }
</style>