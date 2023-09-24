<script lang="ts">
    import type {PageData} from './$types';
    import {formatDistance} from 'date-fns';
    import numeral from 'numeral';
    import IconArrow from '~icons/bxs/up-arrow'

    export let data: PageData;
    console.log(data);
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

    function formatPercent(percent: number): string {
        return numeral(percent * 100).format('0.0')
    }

    function toggleEncounter(enc: number): Promise<Response> {
        return async () => {
            if (details[enc] == undefined) {
                let resp = await (await fetch("/api/logs/" + enc)).json()
                details[enc] = {
                    players: new Map<string, object>(),
                    buffs: resp.buffs,
                    debuffs: resp.debuffs,
                    hpLog: resp.hpLog,
                    partyInfo: resp.partyInfo,
                };

                let partyLookup = new Map<string, string>();
                if (details[enc].partyInfo) {
                    let parties = Object.keys(details[enc].partyInfo);

                    for (let i = 0; i < parties.length; i++) {
                        let party = details[enc].partyInfo[parties[i]];
                        for (let j = 0; j < party.length; j++) {
                            partyLookup.set(party[j], parties[i]);
                        }
                    }
                }
                console.log(partyLookup)

                resp.entities.forEach((entity) => {
                    if (entity.enttype !== "PLAYER") {
                        return
                    }

                    if (details[enc].partyInfo) {
                        entity.party = partyLookup.get(entity.name)
                    }

                    // entity.skills = entity.skills.filter((skill) => skill.totalDamage > 0)
                    entity.skills.sort((a, b) => b.totalDamage - a.totalDamage)
                    entity.opener = getOpener(entity)
                    details[enc].players.set(entity.name, entity)
                })
                sort[enc] = [...details[enc].players.keys()]
                sort[enc].sort((a, b) => details[enc].players.get(b).damage - details[enc].players.get(a).damage)

                groupSynergies(enc)
            }

            toggled[enc] = !toggled[enc]
        }
    }

    function getPartyColor(enc, name) {
        let player = details[enc].players.get(name);
        switch (player.party) {
            case "0":
                return "text-fuchsia-700"
            case "1":
                return "text-teal-700"
            default:
                return "text-amber-700"
        }
    }

    function getOpener(player) {
        let opener = [];
        for (let i = 0; i < player.skills.length; i++) {
            let skill = player.skills[i];
            if (skill.id === 0
                || skill.name.includes("(Summon)")
                || skill.name === "Stand Up"
                || skill.name === "Weapon Attack"
                || skill.name.includes("Basic Attack")) {
                continue;
            }
            for (let j = 0; j < skill.castLog.length; j++) {
                if (skill.castLog[j] < 20000) {
                    opener.push({
                        name: player.skills[i].name,
                        time: skill.castLog[j],
                    })
                }
            }
        }
        opener.sort((a, b) => a.time - b.time)
        return opener;
    }

    function groupSynergiesAdd(synergies, buffId, buff) {
        if (!["classskill", "identity", "ability"].includes(buff.buffCategory) || buff.target !== "PARTY") {
            return
        }

        if (![1, 2, 4, 128].includes(buff.buffType)) {
            return
        }

        let key = buff.uniqueGroup ? buff.uniqueGroup : buff.source.skill?.classId + "_" + buff.source.skill?.name;
        if (!synergies.has(key)) {
            synergies.set(key, new Map<string, object>());
        }
        synergies.get(key).set(buffId, buff);
    }

    function groupSynergies(enc: number) {
        let encounter = details[enc];
        let buffs = Object.keys(encounter.buffs);
        let debuffs = Object.keys(encounter.debuffs)

        let synergies = new Map<string, Map<string, object>>();
        for (let i = 0; i < buffs.length; i++) {
            let buff = encounter.buffs[buffs[i]];
            groupSynergiesAdd(synergies, buffs[i], buff);
        }
        for (let i = 0; i < debuffs.length; i++) {
            let debuff = encounter.debuffs[debuffs[i]];
            groupSynergiesAdd(synergies, debuffs[i], debuff);
        }

        return synergies
    }

    function calculateSynergies(players, synergies) {

    }

    function focus(enc: number, selected: string) {
        return () => focused[enc] = selected
    }

    function inspect(enc: number, selected: string) {
        return () => view[enc] = selected
    }
</script>

<div class="mx-3 my-14 mx-auto">
    <div class="columns-1 text-center ...">
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
        <br/>
        <p class="font-semibold text-pink-600">THE UI IS TEMP BARE MIN FOR DEV VELOCITY</p>
        <p>- faust ^.^</p>
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
                            <div>
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
                                    <span class="{getPartyColor(encounter.id, focused[encounter.id])} font-medium">{focused[encounter.id]}</span>
                                    <button on:click={focus(encounter.id, "")}>(back)</button>
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
                                                    <td class="{getPartyColor(encounter.id, player)} font-medium">
                                                        <button on:click={focus(encounter.id, player)}>{player}</button>
                                                    </td>
                                                    <td>{details[encounter.id].players.get(player).class}</td>
                                                    <td>{formatDamage(details[encounter.id].players.get(player).dps)}</td>
                                                    <td>{formatDamage(details[encounter.id].players.get(player).damage)}</td>
                                                </tr>
                                            {/each}
                                        </table>
                                    {:else}
                                        <table class="table-auto inline-block">
                                            <thead>
                                            <tr>
                                                <th>Name</th>
                                                <th>DPS</th>
                                                <th>DMG</th>
                                                <th>Crit</th>
                                                <th>Avg</th>
                                                <th>Max</th>
                                                <th>Casts</th>
                                                <th>Hits</th>
                                                <th>Casts/m</th>
                                                <th>Hits/m</th>
                                            </tr>
                                            </thead>
                                            <tr>
                                                <td class="{getPartyColor(encounter.id, focused[encounter.id])} font-medium">{focused[encounter.id]}</td>
                                                <td>{formatDamage(details[encounter.id].players.get(focused[encounter.id]).dps)}</td>
                                                <td>{formatDamage(details[encounter.id].players.get(focused[encounter.id]).damage)}</td>
                                            </tr>
                                            {#each details[encounter.id].players.get(focused[encounter.id]).skills as skill (skill.id)}
                                                {@const id = Number(skill.id)}
                                                {#if !(id >= 19090 && id <= 19099) && !(id >= 19280 && id <= 19287)}
                                                    <tr>
                                                        <td class="break-words">{skill.name}</td>
                                                        <td>{formatDamage(skill.dps)}</td>
                                                        <td>{formatDamage(skill.totalDamage)}</td>
                                                        <td>{formatPercent(skill.crits / skill.hits)}</td>
                                                        <td>{formatDamage(skill.totalDamage / skill.hits)}</td>
                                                        <td>{formatDamage(skill.maxDamage)}</td>
                                                        <td>{skill.casts}</td>
                                                        <td>{skill.hits}</td>
                                                        <td>{formatPercent(skill.casts / (encounter.duration / 1000 / 60) / 100)}</td>
                                                        <td>{formatPercent(skill.hits / (encounter.duration / 1000 / 60) / 100)}</td>
                                                    </tr>
                                                {/if}
                                            {/each}
                                        </table>
                                        <br/>
                                        {@const player = details[encounter.id].players.get(focused[encounter.id])}
                                        <div class="w-1/2 mx-auto">
                                            <p class="text-center">
                                                <span class="font-semibold">Opener: </span>
                                                {#each player.opener as skill, i}
                                                    {skill.name} {i === player.opener.length - 1 ? "" : "→ "}
                                                {/each}
                                            </p>
                                        </div>
                                    {/if}
                                {:else if view[encounter.id] === "buff"}
                                    coming soon
                                {/if}
                            </div>
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

    td {
        padding: 0 0.5em;
    }
</style>