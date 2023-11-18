<script lang="ts">
    import {getClassIcon, getSkillIcon, cards} from "$lib/game";

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
    console.log(synergies)

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
    console.log(skills)

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
</script>

<div class="bg-[#F4EDE9] w-full h-full rounded-lg">
    <table on:contextmenu|preventDefault={() => focus.set("")} class="table-auto w-full">
        <thead class="bg-[#b96d83]">
        <tr>
            <th class="rounded-tl-lg"></th>
            {#each synergies as synergy, i}
                <th class:rounded-tr-lg={i === synergies.length - 1}>
                    <div class="flex flex-row mx-2 items-center justify-center">
                        {#each synergy.buffs as buff}
                            {@const info = encounter.data.buffCatalog[buff]}
                            <img alt={info.name}
                                 class="inline mx-0.5 rounded-sm h-6 w-6"
                                 src="{getSkillIcon(info.icon)}"/>
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
                    <img alt={encounter.players[$focus].class}
                         src="{getClassIcon(encounter.players[$focus].class)}"
                         class="h-6 mr-1.5 inline opacity-95"
                    />
                    {$focus}
                </div>
            </td>
            {#each synergies as synergy}
                {@const percent = player.synergy[synergy.name]?.percent}
                <td>
                    {percent ? percent : ""}
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
        padding: 0 6px;
        z-index: 10;
        position: relative;
        color: #524d72;
    }
</style>