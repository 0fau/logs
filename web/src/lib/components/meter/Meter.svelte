<script lang="ts">
    import {formatDamage, formatDate, formatDuration, formatPercent, formatSeconds} from "$lib/components/meter/print";

    export let encounter;
    export let data;
    let view;
    let focused;
    let hovered;

    let details = {
        players: new Map<string, object>(),
        buffs: data.buffs,
        debuffs: data.debuffs,
        hpLog: data.hpLog,
        partyInfo: data.partyInfo,
        hasPartyInfo: data.partyInfo !== null,
        partyBuffs: {},
    };

    let partyLookup = new Map<string, string>();
    if (details.partyInfo) {
        let parties = Object.keys(details.partyInfo);

        for (let i = 0; i < parties.length; i++) {
            let party = details.partyInfo[parties[i]];
            for (let j = 0; j < party.length; j++) {
                partyLookup.set(party[j], parties[i]);
            }
        }
    }

    data.entities.forEach((entity) => {
        if (entity.enttype !== "PLAYER") {
            return;
        }

        if (details.partyInfo) {
            entity.party = partyLookup.get(entity.name);
        } else {
            entity.party = "0"
        }

        entity.skills.sort((a, b) => b.totalDamage - a.totalDamage);
        entity.opener = getOpener(entity);
        console.log(entity.opener);
        details.players.set(entity.name, entity);
    })
    let sort = [...details.players.keys()];
    sort.sort((a, b) => details.players.get(b).damage - details.players.get(a).damage);

    details.synergies = groupSynergies();

    if (details.partyInfo) {
        let parties = Object.keys(details.partyInfo);
        parties.forEach((p) => {
            details.partyInfo[parties[p]].sort((a, b) => details.players.get(b).damage - details.players.get(a).damage);
            details.partyBuffs[p] = calculateSynergies(details.partyInfo[parties[p]], details.synergies);

        })
    } else {
        details.partyInfo = {"0": sort}
        details.partyBuffs["0"] = calculateSynergies(Array.from(details.players.keys()), details.synergies);
    }

    console.log(details)

    for (let [name, player] of details.players) {
        if (!player.party) {
            player.skillBuffs = new Map<string, Map<string, number>>();
            continue;
        }

        player.skillBuffs = calculateSkillSynergies(details.synergies, details.partyBuffs[player.party ?? "0"], player)
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

    function groupSynergies() {
        let buffs = Object.keys(details.buffs);
        let debuffs = Object.keys(details.debuffs)

        let synergies = new Map<string, Map<string, object>>();
        for (let i = 0; i < buffs.length; i++) {
            let buff = details.buffs[buffs[i]];
            groupSynergiesAdd(synergies, buffs[i], buff);
        }
        for (let i = 0; i < debuffs.length; i++) {
            let debuff = details.debuffs[debuffs[i]];
            groupSynergiesAdd(synergies, debuffs[i], debuff);
        }

        return synergies
    }

    function calculateSynergies(players, synergies) {
        let syns = new Map<string, Map<string, number>>();
        for (let i = 0; i < players.length; i++) {
            let player = details.players.get(players[i]);

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

    function getOpener(player) {
        let opener = [];
        console.log(player.skills)
        for (let i = 0; i < player.skills.length; i++) {
            let skill = player.skills[i];
            console.log(skill.name)
            if (skill.id === 0
                || skill.name.includes("(Summon)")
                || skill.name === "Stand Up"
                || skill.name === "Weapon Attack"
                || skill.name.includes("Basic Attack")
                || skill.icon === "") {
                continue;
            }
            console.log("cast log length")
            console.log(skill.castLog.length)
            for (let j = 0; j < skill.castLog.length; j++) {
                console.log(skill.castLog[j])
                if (skill.castLog[j] < 180000) {
                    opener.push({
                        name: player.skills[i].name,
                        time: skill.castLog[j],
                        icon: skill.icon,
                    })
                }
            }
        }
        opener.sort((a, b) => a.time - b.time)
        console.log(opener)
        return opener;
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


    function focus(selected: string) {
        return () => {
            focused = selected
            hoverBuff(null)()
        }
    }

    function inspect(selected: string) {
        return () => {
            view = selected
            hoverBuff(null)()
        }
    }

    function hoverBuff(buff: object) {
        return () => {
            if (hovered === buff) {
                hovered = null
            } else {
                hovered = buff
            }
        }
    }

    function getPartyColor(name) {
        if (!details.hasPartyInfo) {
            return "text-amber-700"
        }

        let player = details.players.get(name);
        switch (player.party) {
            case "0":
                return "text-fuchsia-700"
            case "1":
                return "text-teal-700"
        }
    }

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

    function sanitizeSynergyDesc(desc) {
        return desc.replaceAll(/<FONT.*>(.*)<\/FONT>/g, "$1")
    }
</script>

<p>
    [#{encounter.id}]
    <span class="underline">Raid: {encounter.raid}</span> {formatDate(encounter.date)}
</p>
<p>{formatDamage(encounter.damage)} damage dealt in {formatDuration(encounter.duration)}</p>
<div>
    <span class="font-semibold">Inspect:</span>
    <button on:click={inspect("")}
            class="text-red-600"
            class:underline={!view}
            class:font-semibold={!view}>
        DMG
    </button>
    <span class="font-semibold">|</span>
    <button on:click={inspect("buff")}
            class="text-purple-600"
            class:underline={view === "buff"}
            class:font-semibold={view === "buff"}>
        BUFF
    </button>
    <span class="font-semibold ml-2">Focus: </span>
    {#if !focused}
        <span class="text-yellow-600 font-semibold">Party</span>
    {:else}
        <span class="{getPartyColor(focused)} font-medium">{focused}</span>
        <button on:click={focus("")}>(back)</button>
    {/if}
    <br/>
    {#if !view}
        {#if !focused}
            <table class="table-auto inline-block">
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Class</th>
                    <th>DPS</th>
                    <th>Damage</th>
                </tr>
                </thead>
                {#each sort as player}
                    <tr>
                        <td class="{getPartyColor(player)} font-medium">
                            <button on:click={focus(player)}>{player}</button>
                        </td>
                        <td>{details.players.get(player).class}</td>
                        <td>{formatDamage(details.players.get(player).dps)}</td>
                        <td>{formatDamage(details.players.get(player).damage)}</td>
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
                    <td class="{getPartyColor(focused)} font-medium">{focused}</td>
                    <td>{formatDamage(details.players.get(focused).dps)}</td>
                    <td>{formatDamage(details.players.get(focused).damage)}</td>
                </tr>
                {#each details.players.get(focused).skills as skill (skill.id)}
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
            {@const player = details.players.get(focused)}
            <div class="w-1/2 mx-auto text-center">
                <p class="font-semibold">Skill Cast Log</p>
                <div class="leading-4">
                    {#each player.opener as skill, i}
                        <div class="inline-block w-12 h-12">
                            <span class="text-xs font-medium">{formatSeconds(skill.time)}s</span>
                            <img on:mouseenter={() => console.log(skill.name)}
                                 class="w-8 h-8 m-auto inline border-2 rounded border-black"
                                 src="/icons/skills/{skill.icon !== "" ? skill.icon : 'unknown.png'}"/>
                        </div>
                        {i === player.opener.length - 1 ? "" : "→ "}
                    {/each}
                </div>
            </div>
        {/if}
    {:else if view === "buff"}
        {#if !focused}
            {@const parties = Object.keys(details.partyInfo)}
            {#each parties as p}
                {#if parties.length > 0}
                    <p class="font-semibold">Party {Number(p) + 1}</p>
                {/if}
                {@const buffInfo = details.synergies}
                {@const buffs = details.partyBuffs[p]}
                {@const sortedSyns = sortSyn(buffs)}
                <table class="table-auto inline-block">
                    <thead>
                    <th>Player</th>
                    {#each sortedSyns as buff}
                        {@const syns = buffInfo.get(buff)}
                        <th>
                            {#each [...syns.values()] as syn}
                                <button on:click={hoverBuff(syn)}>
                                    <img class="inline w-6 h-6 m-0.5 border-2 border-red-600 rounded"
                                         src="/icons/skills/{syn.source.icon}"
                                         alt="{syn.source.name}"
                                    />
                                </button>
                            {/each}
                        </th>
                    {/each}
                    </thead>
                    {#each details.partyInfo[p] as player}
                        <tr>
                            <td class="{getPartyColor(player)} font-medium">
                                <button on:click={focus(player)}>{player}</button>
                            </td>
                            {#each sortedSyns as syn}
                                <td>
                                    {
                                        formatPercent(
                                            details.partyBuffs[p].get(syn).get(player) /
                                            details.players.get(player).damage
                                        )
                                    }
                                </td>
                            {/each}
                        </tr>
                    {/each}
                </table>
            {/each}
        {:else}
            {@const player = details.players.get(focused)}
            {@const buffInfo = details.synergies}
            {@const buffs = details.partyBuffs[player.party ?? "0"]}
            {@const sortedSyns = sortSyn(buffs)}
            <table class="table-auto inline-block">
                <thead>
                <th>Name</th>
                {#each sortedSyns as buff}
                    {#if player.skillBuffs.has(buff)}
                        {@const syns = buffInfo.get(buff)}
                        <th>
                            {#each [...syns.values()] as syn}
                                <button on:click={hoverBuff(syn)}>
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
        {#if hovered}
            {@const buff = hovered}
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

<style lang="postcss">
    td {
        padding: 0 0.5em;
    }
</style>