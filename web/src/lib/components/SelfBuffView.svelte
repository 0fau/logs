<script lang="ts">

    import {getClassIcon, getSkillIcon} from "$lib/game";

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
                            <img alt={info.name} class="inline rounded-sm h-6 w-6" src="{getSkillIcon(info.icon)}"/>
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
                    <button on:click={() => focus.set(name)} class="mt-1">
                        <img alt={player.class}
                             src="{getClassIcon(player.class)}"
                             class="h-6 ml-1 inline -translate-y-0.5 opacity-95"
                        />
                        {name}
                    </button>
                </td>
                {#each buffGroups as buffGroup}
                    {@const percent = data.selfBuff[buffGroup.name]?.percent}
                    <td>
                        {percent ? percent : ""}
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
        z-index: 10;
        position: relative;
        color: #524d72;
    }
</style>