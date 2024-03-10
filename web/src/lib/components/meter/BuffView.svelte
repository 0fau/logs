<script lang="ts">
    import { getClassFromId, getSkillIcon } from "$lib/game";
    import { formatPercent, sanitizeBuffDescription } from "$lib/print.js";
    import Player from "$lib/components/meter/Player.svelte";
    import { onMount } from "svelte";
    import { horizontalWheel } from "$lib/scroll";
    import PlayerName from "$lib/components/meter/PlayerName.svelte";

    export let encounter;
    export let focus;

    let most = 0;
    for (let player of Object.values(encounter.players)) {
        if (player.damage > most) {
            most = player.damage;
        }
    }

    let parties = [];
    encounter.parties.forEach((p) => {
        let party = [];
        p.forEach((name) => {
            if (encounter.players[name]) {
                party.push(name);
            }
        });

        if (party.length > 0) {
            parties.push(party);
        }
    });

    let barColors = ["#eedede", "#eedede"];

    let rows = {};
    let bars = {};

    let cWidth = 0;

    $: {
        for (let [name, bar] of Object.entries(bars)) {
            let percent = encounter.players[name].damage / most;
            bar.style.height = rows[name].clientHeight + "px";
            bar.style.width = percent * cWidth + "px";
            bar.style.backgroundColor = barColors[0];
        }
    }

    let hovered;

    let div: HTMLElement;
    onMount(() => {
        horizontalWheel(div);
    });
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    on:contextmenu|preventDefault={() => {}}
    bind:this={div}
    class="custom-scroll h-full w-full overflow-scroll bg-bouquet-50 relative">
    {#each parties as players, party}
        {@const synergies = encounter.data.synergies[party]}
        <table
            class="w-full table-fixed min-w-[40rem]"
            class:mb-4={party !== parties.length - 1}
            bind:clientWidth={cWidth}>
            <thead class="bg-tapestry-500">
                <tr>
                    <th class="w-8 rounded-tl-lg"></th>
                    <th class="w-44"></th>
                    <th class="w-full"></th>
                    {#each synergies as synergy, i}
                        {@const width = synergy.buffs.length > 1 ? `${synergy.buffs.length * 1.5 + 1.5}rem` : "3.5rem"}
                        <th class:rounded-tr-lg={i === synergies.length - 1} style="width: {width}">
                            <div class="relative mx-2 flex flex-row items-center justify-center">
                                {#each synergy.buffs as buff, i}
                                    {@const info = encounter.data.buffCatalog[buff]}
                                    {@const hoverkey = party + "_" + synergy.name + "_" + buff}
                                    <div class="flex items-center justify-center">
                                        <img
                                            alt={info.name}
                                            class="inline h-6 w-6 rounded-sm"
                                            class:mr-0.5={i !== synergy.buffs.length - 1}
                                            src={getSkillIcon(info.icon)}
                                            on:mouseover={() => {
                                                hovered = hoverkey;
                                                console.log(synergy.name);
                                            }}
                                            on:mouseleave={() => (hovered = "")} />
                                        {#if hovered === hoverkey}
                                            <div
                                                class="absolute z-50 flex -translate-y-[calc(100%-1.5rem)] flex-col items-center justify-center whitespace-nowrap rounded-lg border border-[#c58597] bg-bouquet-50 p-2 text-[#575279]">
                                                <img
                                                    alt={info.name}
                                                    class="inline h-6 w-6 rounded-sm"
                                                    src={getSkillIcon(info.skill.icon)} />
                                                <p class="font-medium">
                                                    [{getClassFromId(info.skill.class)}] {info.skill
                                                        .name}
                                                </p>
                                                <p class="">
                                                    {sanitizeBuffDescription(info.description)}
                                                </p>
                                            </div>
                                        {/if}
                                    </div>
                                {/each}
                            </div>
                        </th>
                    {/each}
                </tr>
            </thead>
            {#each players as name, i}
                {@const player = encounter.players[name]}
                {@const data = encounter.data.players[name]}
                <tr bind:this={rows[name]}>
                    <PlayerName {player} difficulty={encounter.difficulty} {focus} />
                    {#each synergies as synergy}
                        {@const percent = data.synergy[synergy.name]?.percent}
                        {@const hoverkey = name + "_" + synergy.name}
                        <td>
                            <div class="relative flex justify-center">
                                <span
                                    class="z-10 my-auto"
                                    on:mouseover={() => (hovered = hoverkey)}
                                    on:mouseleave={() => (hovered = "")}
                                    >{percent ? percent : ""}</span>
                                {#if hovered === hoverkey}
                                    <div
                                        class="absolute z-50 flex translate-y-[26px] flex-col items-center justify-center whitespace-nowrap rounded-md border border-[#c58597] bg-bouquet-50 p-2 text-[#575279]">
                                        {#each synergy.buffs as buff}
                                            {@const info = encounter.data.buffCatalog[buff]}
                                            <div
                                                class="my-0.5 flex w-[64px] items-center justify-start">
                                                <img
                                                    alt={info.skill.name}
                                                    src={getSkillIcon(info.skill.icon)}
                                                    class="mr-1.5 inline h-6 w-6 rounded-md" />
                                                {formatPercent(
                                                    data.synergy[synergy.name].buffs[buff] /
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
                        bind:this={bars[name]}
                        class="absolute z-0 left-0"
                        class:rounded-bl-lg={i === players.length - 1}>
                    </div>
                </tr>
            {/each}
        </table>
    {/each}
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
