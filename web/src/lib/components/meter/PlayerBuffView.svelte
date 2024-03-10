<script lang="ts">
    import {cards, getClassFromId, getClassIcon, getSkillIcon} from "$lib/game";
    import {formatPercent, sanitizeBuffDescription} from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";
    import SkillName from "$lib/components/meter/SkillName.svelte";
    import PlayerName from "$lib/components/meter/PlayerName.svelte";
    import { onMount } from "svelte";
    import { horizontalWheel } from "$lib/scroll";
    import IconUpwards from '~icons/lets-icons/up';

    export let encounter;
    export let focus;

    let party = 0;
    for (let i = 0; i < encounter.parties.length; i++) {
        if (encounter.parties[i].indexOf($focus) !== -1) {
            party = i
            break
        }
    }
    let synergies = encounter.data.synergies[party];

    let player = encounter.data.players[$focus]
    let skills = Object.keys(player.skillSynergy)
    skills.sort((a, b) => {
        let cmp = player.skillDamage[b].damage - player.skillDamage[a].damage;
        if (cmp !== 0) {
            return cmp
        } else {
            return player.skillDamage[b].casts - player.skillDamage[a].casts;
        }
    })
    skills = skills.filter((skill) => {
        return !cards[skill] || skill === "19282"
    })

    let skillCatalog = encounter.data.skillCatalog;

    let barColors = ["#eedede", "#eedede"];

    let rows = [];
    let bars = [];

    let cWidth = 0;

    let most = player.skillDamage[skills[0]].damage
    $: {
        for (let i = 0; i < rows.length; i++) {
            let percent;
            if (i === 0) {
                percent = 1;
            } else {
                percent = player.skillDamage[skills[i - 1]].damage / most;
            }
            bars[i].style.height = rows[i].clientHeight + "px";
            bars[i].style.width = percent * cWidth + "px";
            bars[i].style.backgroundColor = barColors[i % 2]
        }
    }

    let hovered;

    let div: HTMLElement;
    onMount(() => {
        horizontalWheel(div);
    });
</script>

<div class="custom-scroll h-full w-full overflow-scroll bg-bouquet-50" bind:this={div}>
    <table on:contextmenu|preventDefault={() => focus.set("")} class="table-fixed w-full min-w-[40rem]" bind:clientWidth={cWidth}>
        <thead class="bg-tapestry-500">
        <tr>
            <th class="w-8 rounded-tl-lg">
                <button class="flex items-center pl-1 sm:hidden" on:click={() => focus.set("")}>
                    <IconUpwards class="w-5 h-5"/>
                </button>
            </th>
            <th class="w-44"></th>
            <th class="w-full"></th>
            {#each synergies as synergy, i}
                {@const width = synergy.buffs.length > 1 ? `${synergy.buffs.length * 1.5 + 1.5}rem` : "3.5rem"}
                <th class:rounded-tr-lg={i === synergies.length - 1} style="width: {width}">
                    <div class="flex flex-row mx-2 items-center justify-center">
                        {#each synergy.buffs as buff}
                            {@const info = encounter.data.buffCatalog[buff]}
                            {@const hoverkey = synergy.name + "_" + buff}
                            <img alt={info.name}
                                 class="inline mx-0.5 rounded-sm h-6 w-6"
                                 src="{getSkillIcon(info.icon)}"
                                 on:mouseover={() => hovered = hoverkey}
                                 on:mouseleave={() => hovered = ""}
                            />
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-bouquet-50 border border-[#c58597] -translate-y-[calc(100%-1.5rem)] text-[#575279]">
                                    <img alt={info.name}
                                         class="inline rounded-sm h-6 w-6"
                                         src="{getSkillIcon(info.skill.icon)}"/>
                                    <p class="font-medium">
                                        [{getClassFromId(info.skill.class)}] {info.skill.name}</p>
                                    <p class="">{sanitizeBuffDescription(info.description)}</p>
                                </div>
                            {/if}
                        {/each}
                    </div>
                </th>
            {/each}
        </tr>
        </thead>
        <tr bind:this={rows[0]} class="relative">
            <PlayerName player={encounter.players[$focus]} difficulty={encounter.difficulty} />
            {#each synergies as synergy}
                {@const percent = player.synergy[synergy.name]?.percent}
                {@const hoverkey = "_player_" + synergy.name}
                <td>
                    <div class="flex relative justify-center">
                            <span class="my-auto z-10"
                                  on:mouseover={() => hovered = hoverkey}
                                  on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                        {#if hovered === hoverkey}
                            <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-bouquet-50 translate-y-[26px] border border-[#c58597] text-[#575279]">
                                {#each synergy.buffs as buff}
                                    {@const info = encounter.data.buffCatalog[buff]}
                                    <div class="w-[64px] flex items-center justify-start my-0.5">
                                        <img alt={info.skill.name}
                                             src="{getSkillIcon(info.skill.icon)}"
                                             class="h-6 w-6 mr-1.5 inline rounded-md"/>
                                        {formatPercent(player.synergy[synergy.name].buffs[buff] / encounter.players[$focus].damage)}
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </td>
            {/each}
            <div bind:this={bars[0]}
                 class="absolute z-0 left-0">
            </div>
        </tr>
        {#each skills as name, i}
            <tr bind:this={rows[i + 1]}>
                <SkillName info={skillCatalog[name]} />
                {#each synergies as synergy}
                    {@const percent = player.skillSynergy[name][synergy.name]?.percent}
                    {@const hoverkey = name + "_" + synergy.name}
                    <td>
                        <div class="flex relative justify-center">
                        <span class="my-auto z-10"
                              on:mouseover={() => hovered = hoverkey}
                              on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-bouquet-50 translate-y-[26px] border border-[#c58597] text-[#575279]">
                                    {#each synergy.buffs as buff}
                                        {@const info = encounter.data.buffCatalog[buff]}
                                        <div class="w-[64px] flex items-center justify-start my-0.5">
                                            <img alt={info.skill.name}
                                                 src="{getSkillIcon(info.skill.icon)}"
                                                 class="h-6 w-6 mr-1.5 inline rounded-md"/>
                                            {formatPercent(player.skillSynergy[name][synergy.name].buffs[buff] / player.skillDamage[name].damage)}
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    </td>
                {/each}
                <div bind:this={bars[i + 1]}
                     class="absolute z-0 left-0"
                     class:rounded-bl-lg={i === skills.length - 1}>
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
        padding: 0 6px;
        z-index: 10;
        position: relative;
        text-align: center;
        color: theme("colors.zinc.600");
    }
</style>
