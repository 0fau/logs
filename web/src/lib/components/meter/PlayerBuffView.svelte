<script lang="ts">
    import {cards, getClassFromId, getClassIcon, getSkillIcon} from "$lib/game";
    import {formatPercent, sanitizeBuffDescription} from "$lib/print";
    import Player from "$lib/components/meter/Player.svelte";

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
            bars[i].style.width = percent * rows[i].clientWidth + "px";
            bars[i].style.backgroundColor = barColors[i % 2]
        }
    }

    let hovered;
</script>

<div class="bg-bouquet-50 w-full h-full rounded-lg">
    <table on:contextmenu|preventDefault={() => focus.set("")} class="table-auto w-full">
        <thead class="bg-tapestry-500">
        <tr>
            <th class="rounded-tl-lg"></th>
            {#each synergies as synergy, i}
                <th class:rounded-tr-lg={i === synergies.length - 1}>
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
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-bouquet-50 border-[1px] border-[#c58597] -translate-y-[calc(100%-1.5rem)] text-[#575279]">
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
            <div bind:this={bars[0]}
                 class="absolute z-0">
            </div>
            <td class="float-left">
                <div class="my-1 flex justify-center items-center">
                    <Player player={encounter.players[$focus]} anonymized={encounter.anonymized} difficulty={encounter.difficulty}/>
                </div>
            </td>
            {#each synergies as synergy}
                {@const percent = player.synergy[synergy.name]?.percent}
                {@const hoverkey = "_player_" + synergy.name}
                <td>
                    <div class="flex relative justify-center">
                            <span class="my-auto z-10"
                                  on:mouseover={() => hovered = hoverkey}
                                  on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                        {#if hovered === hoverkey}
                            <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-bouquet-50 translate-y-[26px] border-[1px] border-[#c58597] text-[#575279]">
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
        </tr>
        {#each skills as name, i}
            <tr bind:this={rows[i + 1]}>
                <div bind:this={bars[i + 1]}
                     class="absolute z-0"
                     class:rounded-bl-lg={i === skills.length - 1}>
                </div>
                <td class="float-left">
                    <div class="my-1 flex justify-center items-center"
                         class:mb-1.5={i === skills.length - 1}>
                        <img alt={skillCatalog[name].name}
                             src="{getSkillIcon(skillCatalog[name].icon)}"
                             class="h-6 w-6 mr-1.5 inline opacity-95"
                        />
                        <span>{skillCatalog[name].name}</span>
                    </div>
                </td>
                {#each synergies as synergy}
                    {@const percent = player.skillSynergy[name][synergy.name]?.percent}
                    {@const hoverkey = name + "_" + synergy.name}
                    <td>
                        <div class="flex relative justify-center">
                        <span class="my-auto z-10"
                              on:mouseover={() => hovered = hoverkey}
                              on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-bouquet-50 translate-y-[26px] border-[1px] border-[#c58597] text-[#575279]">
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