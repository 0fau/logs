<script lang="ts">
    import type {PageData} from './$types';
    import {getBossIcon, getRaid} from "$lib/raids";
    import {formatDamage, formatDate, formatDuration} from "$lib/components/meter/print";
    import Meter from '$lib/components/Meter.svelte'

    export let data: PageData;
    $: encounter = data.encounter
</script>

<svelte:head>
    {#if encounter}
        {@const raid = getRaid(encounter.boss)}
        <title>[#{encounter.id}] {raid ? raid.raid.toLowerCase() + ' g' + raid.gate : encounter.boss.toLowerCase()} ({encounter.localPlayer.toLowerCase()})</title>
    {:else}
        <title>HUH</title>
    {/if}
</svelte:head>

<div class="my-14 w-full min-w-[1512px] flex flex-col justify-center items-center text-center">
    {#if encounter}
        <div class="'h-[80px] w-[462px] flex border-[0.5px] border-[#c58597] shadow-sm rounded-md bg-[#F4EDE9]">
            <div class="w-full h-full flex flex-row ml-5 items-center">
                <div>
                    <div class="self-start text-left text-[#575279]">
                        <div>
                            <span class="font-medium">[#{encounter.id}]</span>
                            <img alt={encounter.boss} src={getBossIcon(encounter.boss)}
                                 class="inline w-6 h-6 -translate-y-0.5"/>
                            <span class="font-medium">{encounter.boss}</span>
                        </div>
                        <p class="text-sm">{formatDamage(encounter.damage)} damage dealt
                            in {formatDuration(encounter.duration)}</p>
                        <p class="text-xs text-[#5d5978]">{formatDate(encounter.date)}</p>
                    </div>
                </div>
                <div class="py-1 px-1.5 h-full ml-auto self-end flex flex-col rounded-r-md text-white">
                    <span class="text-xs text-center self-end text-[#F4EDE9] p-0.5 px-1 mr-0.5 mt-1.5 rounded-sm bg-[#b4637a] font-medium">{encounter.localPlayer}</span>
                    <span class="text-xs text-[#b4637a] self-end text-right mr-0.5 mt-0.5 font-medium">{encounter.players[encounter.localPlayer].class}</span>
                    <span class="text-[#575279] text-right mr-1 my-auto text-lg font-medium">{formatDamage(encounter.players[encounter.localPlayer].dps)}</span>
                </div>
            </div>
        </div>
        <Meter {encounter}/>
    {:else}
        <p>uh oh log not found :p</p>
    {/if}
</div>

<style lang="postcss">
    :global(html) {
        background-color: #faf4ed;
    }
</style>