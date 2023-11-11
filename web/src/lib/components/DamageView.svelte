<script lang="ts">
    import {getClassIcon} from "$lib/game";
    import {formatDamage, formatPercent} from "$lib/components/meter/print";

    export let encounter;

    let players = Object.keys(encounter.players)
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage)

    let hasCritDamage = false;
    let hasFA = false;
    let hasBA = false;

    for (let player of Object.values(encounter.data.players)) {
        if (player.damage.critDamage !== "0.0") {
            hasCritDamage = true;
        }

        if (player.damage.fa !== "0.0") {
            hasFA = true;
        }

        if (player.damage.ba !== "0.0") {
            hasBA = true;
        }

        if (hasCritDamage && hasFA && hasBA) {
            break
        }
    }
</script>

<div class="bg-[#F4EDE9] w-full h-full rounded-lg">
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
            <th class="rounded-tr-lg">B%</th>
        </tr>
        </thead>
        <tbody>
        {#each players as name}
            {@const player = encounter.players[name]}
            {@const data = encounter.data.players[name]}
            <tr class="relative">
                <td class="float-left">
                    <div class="mt-1">
                        <img alt={player.class}
                             src="{getClassIcon(player.class)}"
                             class="h-6 mr-1 inline -translate-y-0.5 opacity-95"
                        />
                        {name}
                    </div>
                </td>
                <td>{formatDamage(player.damage)}</td>
                <td>{formatDamage(player.dps)}</td>
                <td>{formatPercent(player.damage / encounter.damage)}
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
            </tr>
        {/each}
        </tbody>
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
        padding: 0 5px;
        color: #524d72;
    }
</style>