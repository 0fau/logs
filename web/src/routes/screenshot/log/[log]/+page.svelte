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
    <EncounterPreview width="w-[462px]" screenshot={true} {encounter}/>
{/if}