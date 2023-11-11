<script lang="ts">

    import {getClassIcon, getSkillIcon} from "$lib/game";

    export let encounter;
    export let focus;

    let players = Object.keys(encounter.players)
    players.sort((a, b) => encounter.players[b].damage - encounter.players[a].damage)
</script>

{#each encounter.parties as players, party}
    {@const synergies = encounter.data.synergies[party]}
    <div class="bg-[#F4EDE9] w-full h-full rounded-lg"
         class:mb-4={party !== encounter.parties.length - 1}>
        <table class="table-auto w-full inline-block">
            <thead class="bg-[#b96d83]">
            <tr>
                <th class="w-full rounded-tl-lg"><span class="pl-2 float-left">Party {party + 1}</span></th>
                {#each synergies as synergy, i}
                    <th
                            class:rounded-tr-lg={i === synergies.length - 1}>
                        <div class="flex flex-row items-center justify-center">
                            {#each synergy.buffs as buff}
                                {@const info = encounter.data.buffCatalog[buff]}
                                <img alt={info.name} class="inline rounded-sm h-6 w-6" src="{getSkillIcon(info.icon)}"/>
                            {/each}
                        </div>
                    </th>
                {/each}
            </tr>
            </thead>
            <tbody>
            {#each players as name}
                {@const player = encounter.players[name]}
                {@const data = encounter.data.players[name]}
                <tr class="relative">
                    <td class="float-left">
                        <div class="mt-1">
                            <img alt={player.class}
                                 src="{getClassIcon(player.class)}"
                                 class="h-6 mr-1 inline -translate-y-0.5 opacity-95"
                            />
                            {name}
                        </div>
                    </td>
                    {#each synergies as synergy}
                        {@const percent = data.synergy[synergy.name]?.percent}
                        <td>
                            {percent ? percent : ""}
                        </td>
                    {/each}
                </tr>
            {/each}
            </tbody>
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
        color: #524d72;
    }
</style>