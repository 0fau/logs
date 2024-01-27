<script lang="ts">
    import type {PageData} from './$types';
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";

    export let data: PageData;
    $: encounter = data.encounter
</script>

<svelte:head>
    <title>Screenshot</title>
</svelte:head>

{#if encounter}
<div class="bg-white h-full">
    <div class="w-[28em]">
        <EncounterPreview screenshot={true} {encounter}/>
    </div>
</div>
{/if}