<script lang="ts">
    import type {PageData} from './$types';
    import {formatDistance} from 'date-fns';
    import numeral from 'numeral';
    import IconArrow from '~icons/bxs/up-arrow'

    export let data: PageData;
    let toggled = {};
    let details = {};
    let view = {};
    let focused = {};
    let sort = {};

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
                let entities = await resp.json()
                details[enc] = new Map<string, object>();
                entities.forEach((entity) => {
                    if (entity.enttype !== "PLAYER") {
                        return
                    }

                    entity.skills = entity.skills.filter((skill) => skill.totalDamage > 0)
                    entity.skills.sort((a, b) => b.totalDamage - a.totalDamage)
                    details[enc].set(entity.name, entity)
                })
                sort[enc] = [...details[enc].keys()]
                sort[enc].sort((a, b) => details[enc].get(b).damage - details[enc].get(a).damage)
            }

            toggled[enc] = !toggled[enc]
        }
    }

    function focus(enc: number, selected: string) {
        return () => focused[enc] = selected
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
        <br />
        <p class="font-semibold text-pink-600">THE UI IS TEMP DURING DEV</p>
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
                            <span class="font-semibold ml-2">Focus: </span>
                            {#if !focused[encounter.id]}
                                <span class="text-yellow-600 font-semibold">Party</span>
                            {:else}
                                <span class="text-green-600 font-medium">{focused[encounter.id]}</span>
                                <button class="font-thin" on:click={focus(encounter.id, "")}>(back)</button>
                            {/if}
                            <br/>
                            {#if !view[encounter.id]}
                                {#if !focused[encounter.id]}
                                    <table class="table-auto inline-block">
                                        <thead>
                                        <tr>
                                            <th>Name</th>
                                            <th>Class</th>
                                            <th>DPS</th>
                                            <th>Damage</th>
                                        </tr>
                                        </thead>
                                        {#each sort[encounter.id] as player}
                                            <tr>
                                                <td class="text-green-600 font-medium"><button on:click={focus(encounter.id, player)}>{player}</button></td>
                                                <td>{details[encounter.id].get(player).class}</td>
                                                <td>{formatDamage(details[encounter.id].get(player).dps)}</td>
                                                <td>{formatDamage(details[encounter.id].get(player).damage)}</td>
                                            </tr>
                                        {/each}
                                    </table>
                                {:else}
                                    <table class="table-auto inline-block">
                                        <thead>
                                        <tr>
                                            <th>Name</th>
                                            <th>DPS</th>
                                            <th>Damage</th>
                                            <th>Casts</th>
                                            <th>Hits</th>
                                        </tr>
                                        </thead>
                                        <tr>
                                            <td class="text-green-600 font-medium">{focused[encounter.id]}</td>
                                            <td>{formatDamage(details[encounter.id].get(focused[encounter.id]).dps)}</td>
                                            <td>{formatDamage(details[encounter.id].get(focused[encounter.id]).damage)}</td>
                                        </tr>
                                        {#each details[encounter.id].get(focused[encounter.id]).skills as skill}
                                            <tr>
                                                <td>{skill.name}</td>
                                                <td>{formatDamage(skill.dps)}</td>
                                                <td>{formatDamage(skill.totalDamage)}</td>
                                                <td>{skill.casts}</td>
                                                <td>{skill.hits}</td>
                                            </tr>
                                        {/each}
                                    </table>
                                {/if}
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