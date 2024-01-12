<script lang="ts">
    import {formatDamage, formatPercent} from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";

    export let encounter;
    export let focus;

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

    let barColors = ["#eedede", "#eedede"];

    let rows = [];
    let bars = [];

    let most = encounter.players[players[0]].damage
    $: {
        for (let i = 0; i < rows.length; i++) {
            let percent = encounter.players[players[i]].damage / most;
            bars[i].style.height = rows[i].clientHeight + "px";
            bars[i].style.width = percent * rows[i].clientWidth + "px";
            bars[i].style.backgroundColor = barColors[i % 2]
        }
    }
</script>

<div on:contextmenu|preventDefault={() => {}}
     class="bg-[#F4EDE9] overflow-hidden w-full h-full">
    <table class="table-auto w-full">
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
            <th class="rounded-tr-lg">B%</th>
        </tr>
        </thead>
        {#each players as name, i}
            {@const player = encounter.players[name]}
            {@const data = encounter.data.players[name]}
            <tr bind:this={rows[i]}>
                <div bind:this={bars[i]}
                     class="absolute z-0"
                     class:rounded-bl-lg={i === players.length - 1}>
                </div>
                <td class="float-left">
                    <button class="py-1 z-50 flex justify-center items-center"
                            on:click={() => focus.set(name)}>
                        <Player player={player} anonymized={encounter.anonymized} difficulty={encounter.difficulty}/>
                    </button>
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
                <td class="rounded-b-xl">{formatPercent(data.damage.brand / 100)}</td>
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
        padding: 0 5px;
        z-index: 10;
        position: relative;
        color: #524d72;
    }
</style>