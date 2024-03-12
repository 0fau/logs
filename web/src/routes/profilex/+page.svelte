<script lang="ts">
    import type {PageData} from './$types';
    import Profile from "$lib/components/Profile.svelte";

    export let data: PageData;
    console.log(data)
</script>

<svelte:head>
    <title>Profile</title>
</svelte:head>

{#if data.profile}
    <Profile profile={data.profile}/>
{/if}

<style lang="postcss">
    :global(html) {
        background-color: #faf4ed;
    }
</style>