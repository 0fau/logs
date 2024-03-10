<script lang="ts">
    import { formatDamage, formatPercent } from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";
    import { onMount } from "svelte";
    import { horizontalWheel } from "$lib/scroll";
    import PlayerName from "$lib/components/meter/PlayerName.svelte";

    export let encounter;
    export let focus;

    let players = Object.keys(encounter.players);
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage);

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
            break;
        }
    }

    let barColors = ["#eedede", "#eedede"];

    let rows = [];
    let bars = [];

    let cWidth = 0;

    let most = encounter.players[players[0]].damage;
    $: {
        for (let i = 0; i < rows.length; i++) {
            let percent = encounter.players[players[i]].damage / most;
            bars[i].style.height = rows[i].clientHeight + "px";
            bars[i].style.width = percent * cWidth + "px";
            bars[i].style.backgroundColor = barColors[i % 2];
        }
    }

    let div: HTMLElement;
    onMount(() => {
        horizontalWheel(div);
    });
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    bind:this={div}
    on:contextmenu|preventDefault={() => {}}
    class="custom-scroll h-full w-full overflow-scroll bg-bouquet-50">
    <table class="relative w-full table-fixed min-w-[40rem]" bind:clientWidth={cWidth}>
        <thead class="bg-tapestry-500">
            <tr>
                <th class="w-8 rounded-tl-lg"></th>
                <th class="w-14"></th>
                <th class="w-full"></th>
                <th class="w-16">DMG</th>
                <th class="w-16">DPS</th>
                <th class="w-14">D%</th>
                <th class="w-14">CRIT</th>
                {#if hasCritDamage}
                    <th class="w-14">CDMG</th>
                {/if}
                {#if hasFA}
                    <th class="w-14">FA</th>
                {/if}
                {#if hasBA}
                    <th class="w-14">BA</th>
                {/if}
                <th class="w-14">Buff%</th>
                <th class="w-14 rounded-tr-lg">B%</th>
            </tr>
        </thead>
        {#each players as name, i}
            {@const player = encounter.players[name]}
            {@const data = encounter.data.players[name]}
            <tr bind:this={rows[i]} class="">
                <PlayerName {player} difficulty={encounter.difficulty} {focus} />
<!--                <td class="float-left w-full truncate">-->
<!--                    <button-->
<!--                        class="z-50 flex items-center justify-start py-1"-->
<!--                        on:click={() => focus.set(name)}>-->
<!--                        <Player-->
<!--                            {player}-->
<!--                            anonymized={encounter.anonymized}-->
<!--                            difficulty={encounter.difficulty} />-->
<!--                    </button>-->
<!--                </td>-->
                <td>{formatDamage(player.damage)}</td>
                <td>{formatDamage(player.dps)}</td>
                <td>{formatPercent(player.damage / encounter.damage)}</td>
                <td>{formatPercent(data.damage.crit / 100)}</td>
                {#if hasCritDamage}
                    <td>
                        {#if data.damage.critDamage !== "0.0"}
                            {formatPercent(data.damage.critDamage / 100)}
                        {/if}
                    </td>
                {/if}
                {#if hasFA}
                    <td>
                        {#if data.damage.fa !== "0.0"}
                            {formatPercent(data.damage.fa / 100)}
                        {/if}
                    </td>
                {/if}
                {#if hasBA}
                    <td>
                        {#if data.damage.ba !== "0.0"}
                            {formatPercent(data.damage.ba / 100)}
                        {/if}
                    </td>
                {/if}
                <td>{formatPercent(data.damage.buff / 100)}</td>
                <td class="rounded-b-xl">{formatPercent(data.damage.brand / 100)}</td>
                <div
                    bind:this={bars[i]}
                    class="absolute left-0 z-0"
                    class:rounded-bl-lg={i === players.length - 1}>
                </div>
            </tr>
        {/each}
    </table>
</div>

<style>
    th {
        font-weight: normal;
        color: theme("colors.tapestry.50");
        padding-top: 3px;
        padding-bottom: 3px;
    }

    td {
        padding: 0 5px;
        z-index: 10;
        position: relative;
        text-align: center;
        color: theme("colors.zinc.600");
    }
</style>
