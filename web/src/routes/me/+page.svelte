<script lang="ts">
    import type {PageData} from './$types';
    import {formatDistance} from 'date-fns';
    import numeral from 'numeral';
    import IconArrow from '~icons/bxs/up-arrow'

    export let data: PageData;
    let toggled = {};
    let details = {};
    let view = {};

    function formatDate(date: number): string {
        return formatDistance(new Date(date), new Date(), {addSuffix: true})
    }

    function formatDuration(duration: number): string {
        let date = new Date(duration);
        return date.getMinutes() + 'm' + date.getSeconds() + 's';
    }

    function formatDamage(damage: number): string {
        return numeral(damage).format('0.0a')
    }

    function toggleEncounter(enc: number): Promise<Response> {
        return async () => {
            if (details[enc] == undefined) {
                let resp = await fetch("/api/logs/" + enc)
                details[enc] = await resp.json()
                details[enc] = details[enc].filter((ent) => ent.enttype == "PLAYER")
                details[enc].sort((a, b) => b.damage - a.damage)
            }

            toggled[enc] = !toggled[enc]
        }
    }

    function inspect(enc: number, selected: string) {
        return () => view[enc] = selected
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
        {#if data.recent.length !== 0}
            <div class="m-5">
                <p>[<span class="font-semibold">Recent Logs</span>]</p>
                {#each data.recent as encounter, i (encounter.id)}
                    <div>
                        <p>
                            <button on:click={toggleEncounter(encounter.id)}>
                                <IconArrow width=1em height=0.9em
                                           style="margin-top:-0.1em; display: inline; transform: rotate({toggled[encounter.id] ? 180 : 90}deg)"/>
                                [#{encounter.id}]
                            </button>
                            <span class="underline">Raid: {encounter.raid}</span> {formatDate(encounter.date)}
                        </p>
                        <p>{formatDamage(encounter.damage)} damage dealt in {formatDuration(encounter.duration)}</p>
                        {#if toggled[encounter.id]}
                            <span class="font-semibold">Inspect:</span>
                            <button on:click={inspect(encounter.id, "")}
                                    class="text-red-600"
                                    class:underline={!view[encounter.id]}
                                    class:font-semibold={!view[encounter.id]}>
                                DMG
                            </button>
                            <span class="font-semibold">|</span>
                            <button on:click={inspect(encounter.id, "buff")}
                                    class="text-purple-600"
                                    class:underline={view[encounter.id] === "buff"}
                                    class:font-semibold={view[encounter.id] === "buff"}>
                                BUFF
                            </button>
                            <br/>
                            {#if !view[encounter.id]}
                                <table class="table-auto inline-block">
                                    <thead>
                                    <tr>
                                        <th>Name</th>
                                        <th>Class</th>
                                        <th>DPS</th>
                                        <th>Damage</th>
                                    </tr>
                                    </thead>
                                    {#each details[encounter.id] as player}
                                        <tr>
                                            <td>{player.name}</td>
                                            <td>{player.class}</td>
                                            <td>{formatDamage(player.dps)}</td>
                                            <td>{formatDamage(player.damage)}</td>
                                        </tr>
                                    {/each}
                                </table>
                            {:else if view[encounter.id] === "buff"}
                                hello
                            {/if}
                            <div id="encounter-{encounter.id}"></div>
                        {/if}
                    </div>
                    {#if i !== data.recent.length - 1}
                        <p>---</p>
                    {/if}
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