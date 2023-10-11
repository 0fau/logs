<script lang="ts">
    import Meter from '$lib/components/meter/Meter.svelte'
    import type {PageData} from './$types';

    export let data: PageData;
    console.log(data);
</script>

<svelte:head>
    {#if data.encounter && data.details}
        <title>log {data.encounter.id}</title>
    {:else}
        <title>HUH</title>
    {/if}
</svelte:head>

<div class="mx-3 my-14 mx-auto text-center max-w-fit">
    {#if data.encounter && data.details}
        <Meter encounter={data.encounter} data={data.details}/>
    {:else}
        <p>uh oh log not found :p</p>
    {/if}
</div>