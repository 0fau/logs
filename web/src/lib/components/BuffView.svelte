<script lang="ts">

    import {getClassIcon, getSkillIcon} from "$lib/game";

    export let encounter;
    export let focus;

    let players = Object.keys(encounter.players)
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage)

    let barColors = ["#eedede", "#eedede"];

    let rows = [];
    let bars = [];

    let most = encounter.players[players[0]].damage
    $: {
        for (let i = 0; i < rows.length; i++) {
            let name = encounter.parties[Math.floor(i / 4)][i % 4]
            let percent = encounter.players[name].damage / most;
            bars[i].style.height = rows[i].clientHeight + "px";
            bars[i].style.width = percent * rows[i].clientWidth + "px";
            bars[i].style.backgroundColor = barColors[i % 2]
        }
    }
</script>

{#each encounter.parties as players, party}
    {@const synergies = encounter.data.synergies[party]}
    <div on:contextmenu|preventDefault={() => {}}
         class="bg-[#F4EDE9] w-full h-full"
         class:mb-4={party !== encounter.parties.length - 1}>
        <table class="table-auto w-full">
            <thead class="bg-[#b96d83]">
            <tr>
                <th class="rounded-tl-lg"><span class="pl-2 float-left"></th>
                {#each synergies as synergy, i}
                    <th class:rounded-tr-lg={i === synergies.length - 1}>
                        <div class="flex flex-row mx-2 items-center justify-center">
                            {#each synergy.buffs as buff, i}
                                {@const info = encounter.data.buffCatalog[buff]}
                                <img alt={info.name}
                                     class="inline rounded-sm h-6 w-6"
                                     class:mr-0.5={i !== synergy.buffs.length - 1}
                                     src="{getSkillIcon(info.icon)}"/>
                            {/each}
                        </div>
                    </th>
                {/each}
            </tr>
            </thead>
            {#each players as name, i}
                {@const player = encounter.players[name]}
                {@const data = encounter.data.players[name]}
                <tr bind:this={rows[i + (party * 4)]}>
                    <div bind:this={bars[i + (party * 4)]}
                         class="absolute z-0"
                         class:rounded-bl-lg={i === players.length - 1}>
                    </div>
                    <td class="float-left">
                        <button class="my-1 flex justify-center items-center"
                                on:click={() => focus.set(name)}>
                            <img alt={player.class}
                                 src="{getClassIcon(player.class)}"
                                 class="h-6 w-6 mr-1.5 inline opacity-95"
                            />
                            {name}
                        </button>
                    </td>
                    {#each synergies as synergy}
                        {@const percent = data.synergy[synergy.name]?.percent}
                        <td>
                            <span class="my-auto">{percent ? percent : ""}</span>
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
        z-index: 10;
        position: relative;
        color: #524d72;
    }
</style>