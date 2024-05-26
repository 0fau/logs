<script lang="ts">

    import {getClassFromId, getSkillIcon} from "$lib/game";
    import {formatPercent, sanitizeBuffDescription} from "$lib/print.js";
    import Player from "$lib/components/meter/Player.svelte";

    export let encounter;
    export let focus;

    let most = 0;
    for (let player of Object.values(encounter.players)) {
        if (player.damage > most) {
            most = player.damage
        }
    }

    let parties = [];
    encounter.parties.forEach(p => {
        let party = [];
        p.forEach((name) => {
            if (encounter.players[name]) {
                party.push(name);
            }
        })

        if (party.length > 0) {
            parties.push(party);
        }
    })


    let barColors = ["#eedede", "#eedede"];

    let rows = {};
    let bars = {};

    $: {
        for (let [name, bar] of Object.entries(bars)) {
            let percent = encounter.players[name].damage / most;
            bar.style.height = rows[name].clientHeight + "px";
            bar.style.width = percent * rows[name].clientWidth + "px";
            bar.style.backgroundColor = barColors[0]
        }
    }

    let hovered;
</script>

{#each parties as players, party}
    {@const synergies = encounter.data.synergies[party]}
    <div on:contextmenu|preventDefault={() => {}}
         class="bg-[#F4EDE9] w-full h-full"
         class:mb-4={party !== parties.length - 1}>
        <table class="table-auto w-full">
            <thead class="bg-[#b96d83]">
            <tr>
                <th class="rounded-tl-lg"><span class="pl-2 float-left"></th>
                {#each synergies as synergy, i}
                    <th class:rounded-tr-lg={i === synergies.length - 1}>
                        <div class="flex flex-row mx-2 items-center justify-center relative">
                            {#each synergy.buffs as buff, i}
                                {@const info = encounter.data.buffCatalog[buff]}
                                {@const hoverkey = party + "_" + synergy.name + "_" + buff}
                                <div class="flex items-center justify-center">
                                    <img alt={info.name}
                                         class="inline rounded-sm h-6 w-6"
                                         class:mr-0.5={i !== synergy.buffs.length - 1}
                                         src="{getSkillIcon(info.icon)}"
                                         on:mouseover={() => { hovered = hoverkey; console.log(synergy.name)}}
                                         on:mouseleave={() => hovered = ""}
                                    />
                                    {#if hovered === hoverkey}
                                        <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-[#F4EDE9] border-[1px] border-[#c58597] -translate-y-[calc(100%-1.5rem)] text-[#575279]">
                                            <img alt={info.name}
                                                 class="inline rounded-sm h-6 w-6"
                                                 src="{getSkillIcon(info.skill.icon)}"/>
                                            <p class="font-medium">
                                                [{getClassFromId(info.skill.class)}] {info.skill.name}</p>
                                            <p class="">{sanitizeBuffDescription(info.description)}</p>
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
                    <div bind:this={bars[name]}
                         class="absolute z-0"
                         class:rounded-bl-lg={i === players.length - 1}>
                    </div>
                    <td class="float-left">
                        <button class="my-1 flex justify-center items-center"
                                on:click={() => focus.set(name)}>
                            <Player boss={encounter.boss} player={player} anonymized={encounter.anonymized}
                                    difficulty={encounter.difficulty}/>
                        </button>
                    </td>
                    {#each synergies as synergy}
                        {@const percent = data.synergy[synergy.name]?.percent}
                        {@const hoverkey = name + "_" + synergy.name}
                        <td>
                            <div class="flex relative justify-center">
                            <span class="my-auto z-10"
                                  on:mouseover={() => hovered = hoverkey}
                                  on:mouseleave={() => hovered = ""}>{percent ? percent : ""}</span>
                                {#if hovered === hoverkey}
                                    <div class="absolute flex flex-col items-center justify-center p-2 z-50 rounded-md whitespace-nowrap bg-[#F4EDE9] translate-y-[26px] border-[1px] border-[#c58597] text-[#575279]">
                                        {#each synergy.buffs as buff}
                                            {@const info = encounter.data.buffCatalog[buff]}
                                            <div class="w-[64px] flex items-center justify-start my-0.5">
                                                <img alt={info.skill.name}
                                                     src="{getSkillIcon(info.skill.icon)}"
                                                     class="h-6 w-6 mr-1.5 inline rounded-md"/>
                                                {formatPercent(data.synergy[synergy.name].buffs[buff] / player.damage)}
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
{/each}

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