<script lang="ts">

    import {cards, getSkillIcon} from "$lib/game";
    import {formatPercent, sanitizeBuffDescription} from "$lib/print";
    import PlayerName from "$lib/components/meter/PlayerName.svelte";
    import SkillName from "$lib/components/meter/SkillName.svelte";
    import { onMount } from "svelte";
    import { horizontalWheel } from "$lib/scroll";
    import IconUpwards from '~icons/lets-icons/up';

    export let encounter;
    export let focus;

    let data = encounter.data.players[$focus]

    let skills = Object.keys(data.skillSelfBuff)
    skills = skills.filter((skill) => {
        return (!cards[skill] || skill === "19282") && skill !== "_player"
    })
    skills.sort((a, b) => {
        let cmp = data.skillDamage[b].damage - data.skillDamage[a].damage;
        if (cmp !== 0) {
            return cmp
        } else {
            return data.skillDamage[b].casts - data.skillDamage[a].casts;
        }
    })

    let buffGroups = encounter.data.players[$focus].skillSelfBuffs;

    let barColors = ["#eedede", "#eedede"];

    let rows = [];
    let bars = [];

    let cWidth = 0;

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
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div on:contextmenu|preventDefault={() => focus.set("")}
     bind:this={div}
     class="custom-scroll h-full w-full overflow-scroll bg-bouquet-50">
    <table class="table-fixed w-full min-w-[40rem]"
           bind:clientWidth={cWidth}>
        <thead class="bg-tapestry-500">
        <tr>
            <th class="w-8 rounded-tl-lg">
                <button class="flex items-center pl-1 sm:hidden" on:click={() => focus.set("")}>
                    <IconUpwards class="w-5 h-5"/>
                </button>
            </th>
            <th class="w-44"></th>
            <th class="w-full"></th>
            {#each buffGroups as buffGroup, i}
                {@const width =
                    buffGroup.buffs.length > 1
                        ? `${buffGroup.buffs.length * 1.5 + 1.5}rem`
                        : "3.5rem"}
                <th class:rounded-tr-lg={i === buffGroups.length - 1} style="width: {width}">
                    <div class="flex flex-row mx-2 items-center justify-center">
                        {#each buffGroup.buffs as buff}
                            {@const info = encounter.data.buffCatalog[buff]}
                            {@const hoverkey = buffGroup.name + "_" + buff}
                            <img alt={info.name} class="rounded-sm mx-0.5 h-6 w-6" src="{getSkillIcon(info.icon)}"
                                 on:mouseover={() => hovered = hoverkey}
                                 on:mouseleave={() => hovered = ""}
                            />
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-bouquet-50 border border-[#c58597] -translate-y-[calc(100%-1.5rem)] text-[#575279]">
                                    <img alt={info.name}
                                         class="inline rounded-sm h-6 w-6"
                                         src="{getSkillIcon(info.skill?.icon ?? info.icon)}"/>
                                    <p class="font-medium">{info.skill?.name ?? info.name}</p>
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
            {#each data.skillSelfBuffs as buffGroup}
                {@const percent = data.skillSelfBuff["_player"] ? data.skillSelfBuff["_player"][buffGroup.name]?.percent : ""}
                {@const hoverkey = "_player_" + buffGroup.name}
                <td>
                    <div class="flex relative justify-center">
                        <span class="my-auto z-10"
                              on:mouseover={() => hovered = hoverkey}
                              on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                        {#if hovered === hoverkey}
                            <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-bouquet-50 translate-y-[26px] border border-[#c58597] text-[#575279]">
                                {#each buffGroup.buffs as buff}
                                    {@const info = encounter.data.buffCatalog[buff]}
                                    <div class="w-[64px] flex items-center justify-start my-0.5">
                                        <img alt={info.skill?.name ?? info.name}
                                             src="{getSkillIcon(info.skill?.icon ?? info.icon)}"
                                             class="h-6 w-6 mr-1.5 inline rounded-md"/>
                                        {formatPercent(data.skillSelfBuff["_player"][buffGroup.name].buffs[buff] / encounter.players[$focus].damage)}
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
            {@const info = encounter.data.skillCatalog[name]}
            <tr bind:this={rows[i + 1]}>
                <SkillName info={info} />
                {#each data.skillSelfBuffs as buffGroup}
                    {@const percent = data.skillSelfBuff[name] ? data.skillSelfBuff[name][buffGroup.name]?.percent : ""}
                    {@const hoverkey = name + "_" + buffGroup.name}
                    <td>
                        <div class="flex relative justify-center">
                        <span class="my-auto z-10"
                              on:mouseover={() => hovered = hoverkey}
                              on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-bouquet-50 translate-y-[26px] border border-[#c58597] text-[#575279]">
                                    {#each buffGroup.buffs as buff}
                                        {@const info = encounter.data.buffCatalog[buff]}
                                        <div class="w-[64px] flex items-center justify-start my-0.5">
                                            <img alt={info.skill?.name ?? info.name}
                                                 src="{getSkillIcon(info.skill?.icon ?? info.icon)}"
                                                 class="h-6 w-6 mr-1.5 inline rounded-md"/>
                                            {formatPercent(data.skillSelfBuff[name][buffGroup.name].buffs[buff] / data.skillDamage[name].damage)}
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
