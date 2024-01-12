<script lang="ts">

    import {cards, getSkillIcon} from "$lib/game";
    import {formatPercent, sanitizeBuffDescription} from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";

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

    let hovered;
</script>

<div on:contextmenu|preventDefault={() => focus.set("")}
     class="bg-[#F4EDE9] w-full h-full rounded-lg">
    <table class="table-auto w-full">
        <thead class="bg-[#b96d83]">
        <tr>
            <th class="rounded-tl-lg"></th>
            {#each buffGroups as buffGroup, i}
                <th class:rounded-tr-lg={i === buffGroups.length - 1}>
                    <div class="flex flex-row mx-2 items-center justify-center">
                        {#each buffGroup.buffs as buff}
                            {@const info = encounter.data.buffCatalog[buff]}
                            {@const hoverkey = buffGroup.name + "_" + buff}
                            <img alt={info.name} class="rounded-sm mx-0.5 h-6 w-6" src="{getSkillIcon(info.icon)}"
                                 on:mouseover={() => hovered = hoverkey}
                                 on:mouseleave={() => hovered = ""}
                            />
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-[#F4EDE9] border-[1px] border-[#c58597] -translate-y-[calc(100%-1.5rem)] text-[#575279]">
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
            <div bind:this={bars[0]}
                 class="absolute z-0">
            </div>
            <td class="float-left">
                <div class="my-1 flex justify-center items-center">
                    <Player player={encounter.players[$focus]} anonymized={encounter.anonymized} difficulty={encounter.difficulty}/>
                </div>
            </td>
            {#each data.skillSelfBuffs as buffGroup}
                {@const percent = data.skillSelfBuff["_player"] ? data.skillSelfBuff["_player"][buffGroup.name]?.percent : ""}
                {@const hoverkey = "_player_" + buffGroup.name}
                <td>
                    <div class="flex relative justify-center">
                        <span class="my-auto z-10"
                              on:mouseover={() => hovered = hoverkey}
                              on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                        {#if hovered === hoverkey}
                            <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-[#F4EDE9] translate-y-[26px] border-[1px] border-[#c58597] text-[#575279]">
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
        </tr>
        {#each skills as name, i}
            {@const info = encounter.data.skillCatalog[name]}
            <tr bind:this={rows[i + 1]}>
                <div bind:this={bars[i + 1]}
                     class="absolute z-0"
                     class:rounded-bl-lg={i === skills.length - 1}>
                </div>
                <td class="float-left">
                    <div class="my-1 flex justify-center items-center">
                        <img alt={info.name}
                             src="{getSkillIcon(info.icon)}"
                             class="h-6 w-6 mr-1.5 inline opacity-95"
                        />
                        <span>{info.name}</span>
                    </div>
                </td>
                {#each data.skillSelfBuffs as buffGroup}
                    {@const percent = data.skillSelfBuff[name] ? data.skillSelfBuff[name][buffGroup.name]?.percent : ""}
                    {@const hoverkey = name + "_" + buffGroup.name}
                    <td>
                        <div class="flex relative justify-center">
                        <span class="my-auto z-10"
                              on:mouseover={() => hovered = hoverkey}
                              on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-[#F4EDE9] translate-y-[26px] border-[1px] border-[#c58597] text-[#575279]">
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
        position: relative;
        color: #524d72;
    }
</style>