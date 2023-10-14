<script lang="ts">
    import type {LayoutData} from './$types';
    import Sidebar from '$lib/components/Sidebar.svelte';

    export let data: LayoutData;
    console.log(data)
</script>

<svelte:head>
    <title>logs by faust</title>
</svelte:head>

<div class="w-[1512px] h-[860px] bg-[#faf4ed] flex justify-center items-center">
    <div class="w-[86%] h-[86%] flex flex-row shadow rounded-lg">
        <div class="w-[20%] min-w-[220px] border-l-[1px] border-y-[1px] border-[#efdcc5] h-full bg-[#b4637a] rounded-l-lg flex flex-col">
            <Sidebar user={data.me}/>
        </div>
        <div class="w-[80%] min-w-[400px] border-y-[1px] border-r-[1px] rounded-r-lg border-[#efdcc5] h-full bg-[#d6adb9]">
            <div class="w-full h-full text-[#fff] flex flex-col items-center justify-center">
                <slot />
            </div>
        </div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
    }
</style>