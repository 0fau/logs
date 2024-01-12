<script lang="ts">
    import {cards, getClassIcon, getSkillIcon} from "$lib/game";
    import {formatDamage, formatPercent} from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";

    export let encounter;
    export let focus;

    let player = encounter.players[$focus]
    let data = encounter.data.players[$focus]

    let players = Object.keys(encounter.players)
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage)

    let skills = Object.keys(data.skillDamage)
    skills = skills.filter((skill) => {
        return !cards[skill] || skill === "19282"
    })
    skills.sort((a, b) => {
        let cmp = data.skillDamage[b].damage - data.skillDamage[a].damage;
        if (cmp !== 0) {
            return cmp
        } else {
            return data.skillDamage[b].casts - data.skillDamage[a].casts;
        }
    })

    let hasCritDamage = false;
    let hasFA = false;
    let hasBA = false;

    for (let id of Object.keys(data.skillDamage)) {
        const skill = data.skillDamage[id]
        if (skill.critDamage !== "0.0") {
            hasCritDamage = true;
        }

        if (skill.fa !== "0.0") {
            hasFA = true;
        }

        if (skill.ba !== "0.0") {
            hasBA = true;
        }

        if (hasCritDamage && hasFA && hasBA) {
            break
        }
    }

    let barColors = ["#eedede", "#eedede"];

    let rows = [];
    let bars = [];

    let most = data.skillDamage[skills[0]].damage
    $: {
        for (let i = 0; i < rows.length; i++) {
            let percent;
            if (i === 0) {
                percent = 1;
            } else {
                percent = data.skillDamage[skills[i - 1]].damage / most;
            }
            bars[i].style.height = rows[i].clientHeight + "px";
            bars[i].style.width = percent * rows[i].clientWidth + "px";
            bars[i].style.backgroundColor = barColors[i % 2]
        }
    }
</script>

<div class="bg-[#F4EDE9] w-full h-full rounded-lg">
    <table on:contextmenu|preventDefault={() => focus.set("")}
           class="table-auto w-full inline-block">
        <thead class="bg-[#b96d83]">
        <tr>
            <th class="rounded-tl-lg"></th>
            <th>DMG</th>
            <th>DPS</th>
            <th>D%</th>
            <th>CRIT</th>
            {#if hasCritDamage}
                <th>CDMG</th>
            {/if}
            {#if hasFA}
                <th>FA</th>
            {/if}
            {#if hasBA}
                <th>BA</th>
            {/if}
            <th>Buff%</th>
            <th>B%</th>
            <th>APH</th>
            <th>APC</th>
            <th>Max</th>
            <th>Casts</th>
            <th>CPM</th>
            <th>Hits</th>
            <th class="rounded-tr-lg">HPM</th>
        </tr>
        </thead>
        <tr bind:this={rows[0]} class="relative">
            <div bind:this={bars[0]}
                 class="absolute z-0">
            </div>
            <td class="float-left">
                <div class="py-1 flex justify-center items-center">
                    <Player player={player} anonymized={encounter.anonymized} difficulty={encounter.difficulty}/>
                </div>
            </td>
            <td>{formatDamage(player.damage)}</td>
            <td>{formatDamage(player.dps)}</td>
            <td>{formatPercent(player.damage / encounter.damage)}</td>
            <td>{formatPercent(data.damage.crit / 100)}</td>
            {#if hasCritDamage}
                <td>
                    {#if data.damage.critDamage !== '0.0'}
                        {formatPercent(data.damage.critDamage / 100)}
                    {/if}
                </td>
            {/if}
            {#if hasFA}
                <td>
                    {#if data.damage.fa !== '0.0'}
                        {formatPercent(data.damage.fa / 100)}
                    {/if}
                </td>
            {/if}
            {#if hasBA}
                <td>
                    {#if data.damage.ba !== '0.0'}
                        {formatPercent(data.damage.ba / 100)}
                    {/if}
                </td>
            {/if}
            <td>{formatPercent(data.damage.buff / 100)}</td>
            <td>{formatPercent(data.damage.brand / 100)}</td>
            <td>-</td>
            <td>-</td>
            <td>-</td>
            <td>{formatDamage(data.damage.casts)}</td>
            <td>{data.damage.cpm}</td>
            <td>{formatDamage(data.damage.hits)}</td>
            <td>{data.damage.hpm}</td>
        </tr>
        {#each skills as name, i}
            {@const skill = data.skillDamage[name]}
            {@const info = encounter.data.skillCatalog[name]}
            <tr bind:this={rows[i + 1]}>
                <div bind:this={bars[i + 1]}
                     class="absolute z-0"
                     class:rounded-bl-lg={i === players.length - 1}>
                </div>
                <td class="float-left">
                    <div class="my-1 flex justify-center items-center">
                        <img alt={info.name}
                             src="{getSkillIcon(info.icon)}"
                             class="h-6 w-6 mr-1.5 inline opacity-95"
                        />
                        {info.name}
                    </div>
                </td>
                <td>{skill.damage !== 0 ? formatDamage(skill.damage) : ""}</td>
                <td>{skill.dps !== 0 ? formatDamage(skill.dps) : ""}</td>
                <td>{formatPercent(skill.percent / 100)}
                <td>{formatPercent(skill.crit / 100)}</td>
                {#if hasCritDamage}
                    <td>
                        {#if skill.critDamage !== '0.0'}
                            {formatPercent(skill.critDamage / 100)}
                        {/if}
                    </td>
                {/if}
                {#if hasFA}
                    <td>
                        {#if skill.fa !== '0.0'}
                            {formatPercent(skill.fa / 100)}
                        {/if}
                    </td>
                {/if}
                {#if hasBA}
                    <td>
                        {#if skill.ba !== '0.0'}
                            {formatPercent(skill.ba / 100)}
                        {/if}
                    </td>
                {/if}
                <td>{formatPercent(skill.buff / 100)}</td>
                <td>{formatPercent(skill.brand / 100)}</td>
                <td>{skill.aph !== '0.0' ? formatDamage(skill.aph / 100) : ""}</td>
                <td>{skill.apc !== '0.0' ? formatDamage(skill.apc / 100) : ""}</td>
                <td>{skill.max !== 0 ? formatDamage(skill.max) : ""}</td>
                <td>{formatDamage(skill.casts)}</td>
                <td>{skill.cpm}</td>
                <td>{skill.hits !== 0 ? formatDamage(skill.hits) : ""}</td>
                <td class="mr-1">{skill.hpm !== '0.0' ? skill.hpm : ""}</td>
            </tr>
        {/each}
    </table>
</div>

<style>
    th {
        font-weight: normal;
        color: #F4EDE9;
        padding-top: 3px;
        padding-bottom: 3px;
    }

    td {
        padding: 0 6px;
        z-index: 10;
        position: relative;
        color: #524d72;
    }
</style>