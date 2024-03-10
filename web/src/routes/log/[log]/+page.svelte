<script lang="ts">
    import type { PageData } from "./$types";
    import { getRaid } from "$lib/raids";
    import Meter from "$lib/components/meter/Meter.svelte";
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";
    import EncounterToolbar from "$lib/components/EncounterToolbar.svelte";
    import EncounterSocial from "$lib/components/EncounterSocial.svelte";

    const difficultyColors = {
        Inferno: "#9a3148",
        Trial: "#9a3148",
        Challenge: "#625f77",
        Hard: "#b9982e",
        Normal: "#625f77"
    };

    export let data: PageData;
    $: encounter = data.encounter;
    let encounter = data.encounter
    let user = data.me
</script>

<svelte:head>
    <meta name="og:site_name" content="black meowket (alpha)" />
    <meta name="twitter:site" content="black meowket (alpha)" />
    {#if encounter}
        {@const raid = getRaid(encounter.boss)}
        <title
            >[#{encounter.id}] {encounter.difficulty.toLowerCase()}
            {raid ? raid.raid.toLowerCase() + " g" + raid.gate : encounter.boss.toLowerCase()}
            {encounter.anonymized ? "" : "(" + encounter.localPlayer.toLowerCase() + ")"}</title>
        <meta
            property="og:title"
            content="{encounter.difficulty} {raid
                ? raid.raid + ' G' + raid.gate
                : encounter.boss} - {encounter.players[encounter.localPlayer].class}" />
        <meta
            name="twitter:title"
            content="{encounter.difficulty} {raid
                ? raid.raid + ' G' + raid.gate
                : encounter.boss} - {encounter.players[encounter.localPlayer].class}" />
        <meta name="theme-color" content={difficultyColors[encounter.difficulty]} />
        {#if encounter.thumbnail}
            <meta
                property="og:image"
                content="https://logs.fau.dev/images/thumbnail/{encounter.id}" />
            <meta property="og:type" content="article" />
            <meta
                name="twitter:image:src"
                content="https://logs.fau.dev/images/thumbnail/{encounter.id}" />
            <meta name="twitter:card" content="summary_large_image" />
        {:else}
            <meta property="og:type" content="website" />
        {/if}
    {:else}
        <title>HUH</title>
    {/if}
</svelte:head>

<div class="w-full py-5 xs:py-10 sm:py-20">
    {#if encounter}
        <div class="flex flex-col items-center">
            <EncounterToolbar {encounter} {user}/>
            <div class="xs:w-[30rem] w-full px-1">
                <EncounterPreview {encounter} />
            </div>
        </div>
        <Meter {encounter} />
    {:else}
        <p>uh oh log not found :p</p>
    {/if}
</div>
