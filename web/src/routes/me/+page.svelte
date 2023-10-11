<script lang="ts">
    import type {PageData} from './$types';
    import {formatDistance} from 'date-fns';
    import numeral from 'numeral';
    import IconArrow from '~icons/bxs/up-arrow';
    import {formatSeconds} from "$lib/components/meter/print";

    export let data: PageData;
    console.log(data);
    let toggled = {};
    let details = {};
    let view = {};
    let focused = {};
    let sort = {};
    let hovered = {};

    let more = data.recent.length == 5;

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
        let ret = numeral(percent * 100).format('0.0');
        return ret == "0.0" ? "" : ret;
    }

    function toggleEncounter(enc: number): Promise<Response> {
        return async () => {
            if (details[enc] == undefined) {
                let resp = await (await fetch("/api/logs/" + enc + "/details")).json()
                details[enc] = {
                    players: new Map<string, object>(),
                    buffs: resp.buffs,
                    debuffs: resp.debuffs,
                    hpLog: resp.hpLog,
                    partyInfo: resp.partyInfo,
                    hasPartyInfo: resp.partyInfo !== null,
                    partyBuffs: {},
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

                resp.entities.forEach((entity) => {
                    if (entity.enttype !== "PLAYER") {
                        return;
                    }

                    if (details[enc].partyInfo) {
                        entity.party = partyLookup.get(entity.name);
                    } else {
                        entity.party = "0"
                    }

                    entity.skills.sort((a, b) => b.totalDamage - a.totalDamage);
                    entity.opener = getOpener(entity);
                    details[enc].players.set(entity.name, entity);
                })
                sort[enc] = [...details[enc].players.keys()];
                sort[enc].sort((a, b) => details[enc].players.get(b).damage - details[enc].players.get(a).damage);

                details[enc].synergies = groupSynergies(enc);

                if (details[enc].partyInfo) {
                    let parties = Object.keys(details[enc].partyInfo);
                    parties.forEach((p) => {
                        details[enc].partyInfo[parties[p]].sort((a, b) => details[enc].players.get(b).damage - details[enc].players.get(a).damage);
                        details[enc].partyBuffs[p] = calculateSynergies(enc, details[enc].partyInfo[parties[p]], details[enc].synergies);

                    })
                } else {
                    details[enc].partyInfo = {"0": sort[enc]}
                    details[enc].partyBuffs["0"] = calculateSynergies(enc, Array.from(details[enc].players.keys()), details[enc].synergies);
                }

                for (let [name, player] of details[enc].players) {
                    if (!player.party) {
                        player.skillBuffs = new Map<string, Map<string, number>>();
                        continue;
                    }

                    player.skillBuffs = calculateSkillSynergies(details[enc].synergies, details[enc].partyBuffs[player.party ?? "0"], player)
                }
                console.log(details[enc].players);
            }

            toggled[enc] = !toggled[enc]
        }
    }

    function getPartyColor(enc, name) {
        if (!details[enc].hasPartyInfo) {
            return "text-amber-700"
        }

        let player = details[enc].players.get(name);
        switch (player.party) {
            case "0":
                return "text-fuchsia-700"
            case "1":
                return "text-teal-700"
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
                || skill.name.includes("Basic Attack")
                || skill.icon === "") {
                continue;
            }
            for (let j = 0; j < skill.castLog.length; j++) {
                if (skill.castLog[j] < 20000) {
                    opener.push({
                        name: player.skills[i].name,
                        time: skill.castLog[j],
                        icon: skill.icon,
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

        let key = (buff.source.skill?.classId ?? 0) + "_" + (buff.uniqueGroup ? buff.uniqueGroup : buff.source.skill?.name);
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

    function calculateSynergies(enc, players, synergies) {
        let syns = new Map<string, Map<string, number>>();
        let encounter = details[enc];
        for (let i = 0; i < players.length; i++) {
            let player = encounter.players.get(players[i]);

            let buffed = player.buffed;
            let debuffed = player.debuffed;
            synergies.forEach((buffs, key) => {
                buffs.forEach((buff, id) => {
                    if (buffed[id]) {
                        if (!syns.has(key)) {
                            syns.set(key, new Map<string, number>([
                                [player.name, buffed[id]]
                            ]));
                        } else {
                            let val = syns.get(key).get(player.name) ?? 0;
                            syns.get(key).set(player.name, val + buffed[id]);
                        }
                    } else if (debuffed[id]) {
                        if (!syns.has(key)) {
                            syns.set(key, new Map<string, number>([
                                [player.name, debuffed[id]]
                            ]));
                        } else {
                            let val = syns.get(key).get(player.name) ?? 0;
                            syns.get(key).set(player.name, val + debuffed[id]);
                        }
                    }
                })
            })
        }

        return syns
    }

    function calculateSkillSynergies(buffInfo, buffed, player) {
        let m = new Map<string, Map<string, number>>();
        for (let syn of buffed.keys()) {
            for (let buff of buffInfo.get(syn).keys()) {
                for (let skill of player.skills) {
                    let dmg = 0;
                    if (m.has(syn)) {
                        dmg = m.get(syn).get(skill.id) ?? 0;
                    }
                    dmg += skill.buffed[buff] ?? 0;
                    dmg += skill.debuffed[buff] ?? 0;

                    if (dmg > 0) {
                        if (!m.has(syn)) {
                            m.set(syn, new Map<string, number>())
                        }
                        m.get(syn).set(skill.id, dmg)
                    }
                }
            }
        }

        return m
    }

    const classesMap = {
        0: "Unknown",
        101: "Warrior (Male)",
        102: "Berserker",
        103: "Destroyer",
        104: "Gunlancer",
        105: "Paladin",
        111: "Female Warrior",
        112: "Slayer",
        201: "Mage",
        202: "Arcanist",
        203: "Summoner",
        204: "Bard",
        205: "Sorceress",
        301: "Martial Artist (Female)",
        302: "Wardancer",
        303: "Scrapper",
        304: "Soulfist",
        305: "Glaivier",
        311: "Martial Artist (Male)",
        312: "Striker",
        401: "Assassin",
        402: "Deathblade",
        403: "Shadowhunter",
        404: "Reaper",
        501: "Gunner (Male)",
        502: "Sharpshooter",
        503: "Deadeye",
        504: "Artillerist",
        505: "Machinist",
        511: "Gunner (Female)",
        512: "Gunslinger",
        601: "Specialist",
        602: "Artist",
        603: "Aeromancer",
        604: "Alchemist"
    };

    function sortSyn(synergies) {
        let sorted = [...synergies.keys()]
        sorted.sort((a, b) => {
            let asplit = a.split("_");
            let bsplit = b.split("_");

            if (Number(asplit[0]) == Number(bsplit[0])) {
                return Number(asplit[1]) - Number(bsplit[1]);
            } else {
                return Number(bsplit[0]) - Number(asplit[0])
            }

        })
        return sorted;
    }

    async function expand() {
        if (data.recent.length === 0) {
            return
        }

        let past = data.recent[data.recent.length - 1].date

        let url = location.protocol + '//' + location.host;

        const recent = await (await fetch(
            url + "/api/logs/recent?past=" + past
        )).json()

        if (recent.length < 5) {
            more = false
        }

        data.recent = data.recent.concat(recent)
    }

    function sanitizeSynergyDesc(desc) {
        return desc.replaceAll(/<FONT.*>(.*)<\/FONT>/g, "$1")
    }

    function focus(enc: number, selected: string) {
        return () => {
            focused[enc] = selected
            hoverBuff(enc, null)()
        }
    }

    function inspect(enc: number, selected: string) {
        return () => {
            view[enc] = selected
            hoverBuff(enc, null)()
        }
    }

    function hoverBuff(enc: number, buff: object) {
        return () => {
            if (hovered[enc] === buff) {
                hovered[enc] = null
            } else {
                hovered[enc] = buff
            }
        }
    }
</script>

<svelte:head>
    <title>logs by faust</title>
</svelte:head>

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
                                                        <td>
                                                            {skill.name}
                                                        </td>
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
                                        <div class="w-1/2 mx-auto text-center">
                                            <p class="font-semibold">Opener</p>
                                            <div class="leading-4">
                                                {#each player.opener as skill, i}
                                                    <div class="inline-block w-12 h-12">
                                                        <span class="text-xs font-medium">{formatSeconds(skill.time)}
                                                            s</span>
                                                        <img on:mouseenter={() => console.log(skill.name)}
                                                             class="w-8 h-8 m-auto inline border-2 rounded border-black"
                                                             src="/icons/skills/{skill.icon !== "" ? skill.icon : 'unknown.png'}"/>
                                                    </div>
                                                    {i === player.opener.length - 1 ? "" : "→ "}
                                                {/each}
                                            </div>
                                        </div>
                                    {/if}
                                {:else if view[encounter.id] === "buff"}
                                    {#if !focused[encounter.id]}
                                        {@const parties = Object.keys(details[encounter.id].partyInfo)}
                                        {#each parties as p}
                                            {#if parties.length > 0}
                                                <p class="font-semibold">Party {Number(p) + 1}</p>
                                            {/if}
                                            {@const buffInfo = details[encounter.id].synergies}
                                            {@const buffs = details[encounter.id].partyBuffs[p]}
                                            {@const sortedSyns = sortSyn(buffs)}
                                            <table class="table-auto inline-block">
                                                <thead>
                                                <th>Player</th>
                                                {#each sortedSyns as buff}
                                                    {@const syns = buffInfo.get(buff)}
                                                    <th>
                                                        {#each [...syns.values()] as syn}
                                                            <button on:click={hoverBuff(encounter.id, syn)}>
                                                                <img class="inline w-6 h-6 m-0.5 border-2 border-red-600 rounded"
                                                                     src="/icons/skills/{syn.source.icon}"
                                                                     alt="{syn.source.name}"
                                                                />
                                                            </button>
                                                        {/each}
                                                    </th>
                                                {/each}
                                                </thead>
                                                {#each details[encounter.id].partyInfo[p] as player}
                                                    <tr>
                                                        <td class="{getPartyColor(encounter.id, player)} font-medium">
                                                            <button on:click={focus(encounter.id, player)}>{player}</button>
                                                        </td>
                                                        {#each sortedSyns as syn}
                                                            <td>
                                                                {
                                                                    formatPercent(
                                                                        details[encounter.id].partyBuffs[p].get(syn).get(player) /
                                                                        details[encounter.id].players.get(player).damage
                                                                    )
                                                                }
                                                            </td>
                                                        {/each}
                                                    </tr>
                                                {/each}
                                            </table>
                                        {/each}
                                    {:else}
                                        {@const player = details[encounter.id].players.get(focused[encounter.id])}
                                        {@const buffInfo = details[encounter.id].synergies}
                                        {@const buffs = details[encounter.id].partyBuffs[player.party ?? "0"]}
                                        {@const sortedSyns = sortSyn(buffs)}
                                        <table class="table-auto inline-block">
                                            <thead>
                                            <th>Name</th>
                                            {#each sortedSyns as buff}
                                                {#if player.skillBuffs.has(buff)}
                                                    {@const syns = buffInfo.get(buff)}
                                                    <th>
                                                        {#each [...syns.values()] as syn}
                                                            <button on:click={hoverBuff( encounter.id, syn)}>
                                                                <img class="inline w-6 h-6 m-0.5 border-2 border-red-600 rounded"
                                                                     src="/icons/skills/{syn.source.icon}"
                                                                     alt="{syn.source.name}"
                                                                />
                                                            </button>
                                                        {/each}
                                                    </th>
                                                {/if}
                                            {/each}
                                            </thead>
                                            {#each player.skills as skill}
                                                {@const id = Number(skill.id)}
                                                {#if skill.totalDamage > 0 && !(id >= 19090 && id <= 19099) && !(id >= 19280 && id <= 19287)}
                                                    <tr>
                                                        <td>{skill.name}</td>
                                                        {#each sortedSyns as syn}
                                                            {#if player.skillBuffs.has(syn)}
                                                                <td>
                                                                    {
                                                                        formatPercent(
                                                                            (player.skillBuffs.get(syn).get(skill.id) ?? 0) /
                                                                            skill.totalDamage
                                                                        )
                                                                    }
                                                                </td>
                                                            {/if}
                                                        {/each}
                                                    </tr>
                                                {/if}
                                            {/each}
                                        </table>
                                    {/if}
                                    {#if hovered[encounter.id]}
                                        {@const buff = hovered[encounter.id]}
                                        <p class="font-semibold underline">Tooltip</p>
                                        <p>[{classesMap[buff.source.skill.classId]}] <span
                                                class="text-purple-600 font-semibold">{buff.source.name}</span></p>
                                        <p>{sanitizeSynergyDesc(buff.source.desc)}</p>
                                        <img src="/icons/skills/{buff.source.skill.icon}"
                                             class="w-5 h-5 mx-auto border-2 border-red-600 rounded inline"
                                             alt="{buff.source.skill.name}"/>
                                        <span>{buff.source.skill.name}</span>
                                    {/if}
                                {/if}
                            </div>
                        {/if}
                    </div>
                    <p>---</p>
                {/each}
                {#if more}
                    <button class="underline" on:click={expand}>show more</button>
                {:else}
                    end
                {/if}
            </div>
        {/if}
    </div>
</div>

<style lang="postcss">
    td {
        padding: 0 0.5em;
    }
</style>