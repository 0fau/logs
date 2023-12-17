<script lang="ts">
    import type {PageData} from './$types';
    import {getRaid} from "$lib/raids";
    import Meter from '$lib/components/meter/Meter.svelte'
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";
    import EncounterSocial from "$lib/components/EncounterSocial.svelte";

    export let data: PageData;
    $: encounter = data.encounter
</script>

<svelte:head>
    <meta name="og:site_name" content="black meowket (alpha)"/>
    <meta property="og:type" content="website"/>
    {#if encounter}
        {@const raid = getRaid(encounter.boss)}
        <title>[#{encounter.id}] {encounter.difficulty.toLowerCase()} {raid ? raid.raid.toLowerCase() + ' g' + raid.gate : encounter.boss.toLowerCase()} {encounter.anonymized ? "" : "(" + encounter.localPlayer.toLowerCase() + ")"}</title>
        <meta property="og:title"
              content="{encounter.difficulty} {raid ? raid.raid + ' G' + raid.gate : encounter.boss} - {encounter.players[encounter.localPlayer].class}">
        <meta property="og:description" content=""/>
    {:else}
        <title>HUH</title>
    {/if}
</svelte:head>

<div class="my-14 w-full min-w-[1512px] flex flex-col justify-center items-center text-center">
    {#if encounter}
        <EncounterSocial {encounter}/>
        <EncounterPreview width="w-[462px]" {encounter}/>
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