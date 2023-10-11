<script lang="ts">
    import type {PageData} from './$types';
    import Sidebar from '$lib/components/Sidebar.svelte';
    import EncounterList from "$lib/components/EncounterList.svelte";
    import IconSearch from '~icons/gala/search'

    export let data: PageData;
</script>

<svelte:head>
    <title>logs by faust</title>
</svelte:head>

<div class="w-screen h-screen bg-[#faf4ed] flex justify-center items-center">
    <div class="w-[80%] h-[80%] flex flex-row shadow rounded-lg">
        <div class="w-[20%] border-l-[1px] border-y-[1px] border-[#efdcc5] h-full bg-[#b4637a] rounded-l-lg flex flex-col">
            <Sidebar selected="logs"/>
        </div>
        <div class="w-[35%] border-y-[1px] border-[#efdcc5] h-full bg-[#d6adb9]">
            <div class="w-4/5 border-[1px] border-[#efdcc5] h-[50px] text-[#575279] mx-auto mt-10 rounded-2xl bg-[#faf4ed] flex justify-center items-center flex-row">
                <input placeholder="owo owo owo" class="bg-[#faf4ed] w-4/5 outline-0"/>
                <button>
                    <IconSearch class="w-6 h-6"/>
                </button>
            </div>
        </div>
        <div class="w-[45%] border-y-[1px] border-r-[1px] border-[#efdcc5] h-full bg-[#e7cfd6] rounded-r-lg">
            <EncounterList encounters={data.recent}/>
        </div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
    }
</style>