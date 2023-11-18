<script lang="ts">

    import {cards, getClassIcon, getSkillIcon} from "$lib/game";

    export let encounter;
    export let focus;

    let data = encounter.data.players[$focus]

    let skills = Object.keys(data.skillDamage)
    skills = skills.filter((skill) => {
        return !cards[skill] || skill === "19282"
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
                            <img alt={info.name} class="rounded-sm mx-0.5 h-6 w-6" src="{getSkillIcon(info.icon)}"/>
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
            {#each data.skillSelfBuffs as buffGroup}
                {@const percent = data.skillSelfBuff["_player"] ? data.skillSelfBuff["_player"][buffGroup.name]?.percent : ""}
                <td>
                    {percent ? percent : ""}
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