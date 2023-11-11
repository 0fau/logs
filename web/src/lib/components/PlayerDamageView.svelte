<script lang="ts">
    import {getSkillIcon} from "$lib/game";
    import {formatDamage, formatPercent} from "$lib/components/meter/print";

    export let encounter;
    export let focus;

    let data = encounter.data.players[$focus]

    let players = Object.keys(encounter.players)
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage)

    let skills = Object.keys(data.skillDamage)
    skills.sort((a, b) => data.skillDamage[b].damage - data.skillDamage[a].damage)

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
</script>

<table class="table-auto w-full inline-block">
    <thead class="bg-[#b96d83]">
    <tr>
        <th class="w-full rounded-tl-lg"></th>
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
    <tbody>
    {#each skills as sid}
        {@const skill = data.skillDamage[sid]}
        {@const info = encounter.data.skillCatalog[sid]}
        <tr class="relative">
            <td class="float-left">
                <div class="mt-1">
                    <img alt={info.name}
                         src="{getSkillIcon(info.icon)}"
                         class="h-6 mr-1 inline -translate-y-0.5 opacity-95"
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
            <td>{skill.casts}</td>
            <td>{skill.cpm}</td>
            <td>{skill.hits !== 0 ? skill.hits : ""}</td>
            <td>{skill.hpm !== '0.0' ? skill.hpm : ""}</td>
        </tr>
    {/each}
    </tbody>
</table>

<style>
    th {
        font-weight: normal;
        color: #F4EDE9;
        padding-top: 3px;
        padding-bottom: 3px;
    }

    td {
        padding: 0 5px;
        color: #524d72;
    }
</style>