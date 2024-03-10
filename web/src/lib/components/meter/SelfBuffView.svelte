<script lang="ts">
    import { getSkillIcon } from "$lib/game";
    import { formatPercent, sanitizeBuffDescription } from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";
    import PlayerName from "$lib/components/meter/PlayerName.svelte";
    import { onMount } from "svelte";
    import { horizontalWheel } from "$lib/scroll";

    export let encounter;
    export let focus;

    let players = Object.keys(encounter.players);
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage);

    let buffGroups = encounter.data.selfBuffs;

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

    let hovered;

    function getBuffTitle(synergy, buff) {
        console.log(synergy);
        console.log(buff);
        if (buff.set) {
            return "[Set] " + buff.set + ": " + buff.name;
        }
        if (synergy.name === "cook") {
            return buff.name;
        }
        if (synergy.name === "battleitem") {
            return "[Battle Item] " + buff.name;
        }
        if (synergy.name.startsWith("bracelet_")) {
            return "[Bracelet] " + buff.name;
        }
        if (synergy.name.startsWith("elixir")) {
            return "[Elixir] " + buff.name;
        }
        return buff.name;
    }

    let div: HTMLElement;
    onMount(() => {
        horizontalWheel(div);
    });
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    on:contextmenu|preventDefault={() => {}}
    bind:this={div}
    class="custom-scroll relative h-full w-full overflow-scroll bg-bouquet-50">
    <table class="w-full min-w-[40rem] table-fixed" bind:clientWidth={cWidth}>
        <thead class="bg-tapestry-500">
            <tr>
                <th class="w-8 rounded-tl-lg"></th>
                <th class="w-44"></th>
                <th class="w-full"></th>
                {#each buffGroups as buffGroup, i}
                    {@const width =
                        buffGroup.buffs.length > 1
                            ? `${buffGroup.buffs.length * 1.5 + 1.5}rem`
                            : "3.5rem"}
                    <th class:rounded-tr-lg={i === buffGroups.length - 1} style="width: {width}">
                        <div class="mx-2 flex flex-row items-center justify-center">
                            {#each buffGroup.buffs as buff}
                                {@const info = encounter.data.buffCatalog[buff]}
                                {@const hoverkey = buffGroup.name + "_" + buff}
                                <img
                                    alt={info.name}
                                    class="mx-0.5 inline h-6 w-6 rounded-sm"
                                    src={getSkillIcon(info.icon)}
                                    on:mouseover={() => (hovered = hoverkey)}
                                    on:mouseleave={() => (hovered = "")} />
                                {#if hovered === hoverkey}
                                    <div
                                        class="absolute z-50 flex -translate-y-[calc(100%-1.5rem)] flex-col items-center justify-center whitespace-nowrap rounded-lg border border-[#c58597] bg-bouquet-50 p-2 text-[#575279]">
                                        <img
                                            alt={info.name}
                                            class="inline h-6 w-6 rounded-sm"
                                            src={getSkillIcon(info.icon)} />
                                        <p class="font-medium">{getBuffTitle(buffGroup, info)}</p>
                                        <p class="">{sanitizeBuffDescription(info.description)}</p>
                                    </div>
                                {/if}
                            {/each}
                        </div>
                    </th>
                {/each}
            </tr>
        </thead>
        {#each players as name, i}
            {@const player = encounter.players[name]}
            {@const data = encounter.data.players[name]}
            <tr bind:this={rows[i]}>
                <PlayerName {player} difficulty={encounter.difficulty} {focus} />
                {#each buffGroups as buffGroup}
                    {@const percent = data.selfBuff[buffGroup.name]?.percent}
                    {@const hoverkey = name + "_" + buffGroup.name}
                    <td>
                        <div class="relative flex justify-center">
                            <span
                                class="z-10 my-auto"
                                on:mouseover={() => (hovered = hoverkey)}
                                on:mouseleave={() => (hovered = "")}>{percent ? percent : ""}</span>
                            {#if hovered === hoverkey}
                                <div
                                    class="absolute z-50 flex translate-y-[26px] flex-col items-center justify-center whitespace-nowrap rounded-md border border-[#c58597] bg-bouquet-50 p-2 text-[#575279]">
                                    {#each buffGroup.buffs as buff}
                                        {@const info = encounter.data.buffCatalog[buff]}
                                        <div
                                            class="my-0.5 flex w-[64px] items-center justify-start">
                                            <img
                                                alt={info.name}
                                                src={getSkillIcon(info.icon)}
                                                class="mr-1.5 inline h-6 w-6 rounded-md" />
                                            {formatPercent(
                                                data.selfBuff[buffGroup.name].buffs[buff] /
                                                    player.damage
                                            )}
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    </td>
                {/each}
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
        color: #f4ede9;
        padding-top: 3px;
        padding-bottom: 3px;
    }

    td {
        padding: 0 5px;
        position: relative;
        color: #524d72;
    }
</style>
