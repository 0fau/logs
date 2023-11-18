<script lang="ts">
    import type {LayoutData} from './$types';
    import Sidebar from '$lib/components/Sidebar.svelte';
    import {user} from "$lib/store";
    import {settingsUI} from "$lib/menu";
    import Settings from "$lib/components/Settings.svelte";

    export let data: LayoutData;
    console.log(data)

    if (data.me.id) {
        user.set(data.me)
    }

    let settings;
    settingsUI.subscribe(async (open) => {
        if (open) {
            let resp = await fetch("/api/settings", {
                credentials: 'same-origin'
            })
            if (!resp.ok) {
                return
            }

            settings = await resp.json()
        } else {
            settings = undefined
        }
    })
</script>

<svelte:head>
    <title>logs by faust</title>
</svelte:head>

<div class="w-screen h-screen min-w-[1512px] min-h-[860px] bg-[#faf4ed] flex justify-center items-center">
    {#if settings}
        <div class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-xl -translate-x-[50%] -translate-y-[80%] w-[376px] h-[330px] bg-[#faf4ed]">
            <Settings {settings}/>
        </div>
    {/if}
    <div class="w-[1300px] h-[740px] flex flex-row shadow rounded-lg">
        <div class="w-[20%] min-w-[220px] border-l-[1px] border-y-[1px] border-[#efdcc5] h-full bg-[#b4637a] rounded-l-lg flex flex-col">
            <Sidebar user={data.me}/>
        </div>
        <div class="w-[80%] min-w-[400px] border-y-[1px] border-r-[1px] rounded-r-lg border-[#efdcc5] h-full bg-[#debdc7]">
            <div class="w-full h-full text-[#fff] flex flex-col items-center justify-center">
                <slot/>
                <img alt="munchers" class="w-20 h-20 -translate-y-5"
                     src="https://cdn.discordapp.com/emojis/1046992385395138660.webp?size=240&quality=lossless"/>
            </div>
        </div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
    }
</style>