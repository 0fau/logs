<script lang="ts">

    import {getClassIcon, getSkillIcon} from "$lib/game";
    import {formatPercent, sanitizeBuffDescription} from "$lib/print";

    export let encounter;
    export let focus;

    let players = Object.keys(encounter.players)
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage)

    let buffGroups = encounter.data.selfBuffs;

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

    let hovered;

    function getBuffTitle(synergy, buff) {
        console.log(synergy)
        console.log(buff)
        if (buff.set) {
            return "[Set] " + buff.set + ": " + buff.name
        }
        if (synergy.name === "cook") {
            return buff.name
        }
        if (synergy.name === "battleitem") {
            return "[Battle Item] " + buff.name
        }
        if (synergy.name.startsWith("bracelet_")) {
            return "[Bracelet] " + buff.name
        }
        return buff.name
    }
</script>

<div on:contextmenu|preventDefault={() => {}}
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
                            <img alt={info.name}
                                 class="inline mx-0.5 rounded-sm h-6 w-6"
                                 src="{getSkillIcon(info.icon)}"
                                 on:mouseover={() => hovered = hoverkey}
                                 on:mouseleave={() => hovered = ""}
                            />
                            {#if hovered === hoverkey}
                                <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-[#F4EDE9] border-[1px] border-[#c58597] -translate-y-[calc(100%-1.5rem)] text-[#575279]">
                                    <img alt={info.name}
                                         class="inline rounded-sm h-6 w-6"
                                         src="{getSkillIcon(info.icon)}"/>
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
                <div bind:this={bars[i]}
                     class="absolute z-0"
                     class:rounded-bl-lg={i === players.length - 1}>
                </div>
                <td class="float-left">
                    <button on:click={() => focus.set(name)} class="my-1 flex justify-center items-center">
                        <img alt={player.class}
                             src="{getClassIcon(player.class)}"
                             class="h-6 mr-1.5 inline opacity-95"
                        />
                        {name}
                    </button>
                </td>
                {#each buffGroups as buffGroup}
                    {@const percent = data.selfBuff[buffGroup.name]?.percent}
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
                                            <img alt={info.name}
                                                 src="{getSkillIcon(info.icon)}"
                                                 class="h-6 w-6 mr-1.5 inline rounded-md"/>
                                            {formatPercent(data.selfBuff[buffGroup.name].buffs[buff] / player.damage)}
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
        padding: 0 5px;
        position: relative;
        color: #524d72;
    }
</style>