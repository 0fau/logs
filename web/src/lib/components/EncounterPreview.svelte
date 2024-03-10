<script lang="ts">
    import {getBossIcon, shortBossName} from "$lib/raids";
    import {formatDamage, formatDate, formatDateSolid, formatDuration} from "$lib/print";

    const difficultyColors = {
        "Inferno": "#9a3148",
        "Trial": "#9a3148",
        "Extreme": "#9a3148",
        "Challenge": "#625f77",
        "Hard": "#b9982e",
        "Normal": "#625f77",
    };

    export let encounter;
    export let screenshot = false;
    export let width;
    export let gearScore;

    $: player = encounter.players[encounter.localPlayer];
</script>

{#if screenshot}
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
{/if}

<div class="h-[80px] {width} flex border-[0.5px] border-[#c58597] shadow-sm rounded-md bg-[#F4EDE9]"
     class:screenshot={screenshot}>
    <div class="w-full h-full flex flex-row px-[18px] items-center">
        <div>
            <div class="self-start text-left text-[#575279]">
                <div>
                    <img alt={encounter.boss} src={getBossIcon(encounter.boss)}
                         class="inline w-6 h-6 -translate-y-0.5"/>
                    <span class="font-medium">{shortBossName(encounter.boss)}</span>
                </div>
                <p class="text-sm">{formatDamage(encounter.damage)} damage dealt
                    in {formatDuration(encounter.duration)}</p>
                <p class="text-xs text-[#5d5978]"><span class="font-medium">#{encounter.id} |</span> <span
                        class="font-semibold"
                        style="color: {difficultyColors[encounter.difficulty]}">{encounter.difficulty}</span> {screenshot ? formatDateSolid(encounter.date) : formatDate(encounter.date)}
                </p>
            </div>
        </div>
        <div class="py-1 h-full ml-auto self-end flex flex-col rounded-r-md text-white">
            <span class="text-xs text-center self-end text-[#F4EDE9] p-0.5 px-1 mt-1.5 rounded-sm bg-[#b4637a] font-medium"
                  class:bg-[#b4637a]={!encounter.anonymized}
                  class:bg-[#8F708A]={encounter.anonymized}
            >{encounter.localPlayer}</span>
            <span class="text-xs self-end text-right mt-0.5 font-medium"
                  class:text-[#b4637a]={!encounter.anonymized}
                  class:text-[#8F708A]={encounter.anonymized}
            >{gearScore ? Math.floor(player.gearScore) + " " : ""}{player.class}</span>
            <span class="text-[#575279] text-right my-auto text-lg font-medium">{formatDamage(player.dps)}</span>
        </div>
    </div>
</div>

{#if screenshot}
    <style>
        p, span {
            font-family: 'Inter', sans-serif;
        }
    </style>
{/if}