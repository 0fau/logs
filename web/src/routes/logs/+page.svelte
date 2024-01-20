<script lang="ts">
    import type {PageData} from './$types';
    import {blur} from 'svelte/transition';
    import Sidebar from '$lib/components/Sidebar.svelte';
    import Search from '$lib/components/Search.svelte';
    import EncounterList from "$lib/components/EncounterList.svelte";
    import Settings from '$lib/components/Settings.svelte'
    import {settings, user} from '$lib/store'
    import {settingsUI} from '$lib/menu'
    import {browser} from "$app/environment"
    import {writable} from "svelte/store";

    export let data: PageData;
    console.log(data)
    if (data.me.id) {
        user.set(data.me)
    }

    let settingsModal;
    settingsUI.subscribe(async (open) => {
        if (open) {
            let resp = await fetch("/api/settings", {
                credentials: 'same-origin'
            })
            if (!resp.ok) {
                return
            }

            settingsModal = await resp.json()
        } else {
            settingsModal = undefined
        }
    })

    $: greeting = $settings?.logs.announcement;

    const search = writable({});
</script>

<svelte:head>
    <title>black meowket (alpha)</title>
</svelte:head>

<div class="w-screen h-screen min-w-[1512px] min-h-[860px] bg-[#faf4ed] flex justify-center items-center">
    <div class="w-[1300px] h-[740px] flex flex-row shadow rounded-lg absolute">
        {#if browser && greeting}
            <div class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-md -translate-x-[50%] -translate-y-[100%] w-[376px] bg-[#faf4ed]">
                <div class="h-fit">
                    <span class="font-semibold text-lg">hewwo ^.^ </span>
                    <img alt="avatar" class="inline h-8 -translate-y-1"
                         src="https://cdn.discordapp.com/emojis/667830310741737473.gif?size=240&quality=lossless"/>
                    <button on:click={() => $settings.logs.announcement = false} class="float-right font-semibold">[x]
                    </button>
                </div>
                <p>you've stumbled upon... my website! this is a secret underground black market that deals in the
                    forbidden teachings. FRICC THE POPO.</p>
                <p class="text-[#57517a] text-sm">(sign up for the <a href="https://docs.fau.dev/logs/alpha/" class="underline">alpha</a>)</p>
            </div>
        {/if}
        {#if settingsModal}
            <div transition:blur={{duration: 10}}
                 class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-xl -translate-x-[50%] -translate-y-[80%] w-[376px] h-[330px] bg-[#faf4ed]">
                <Settings settings={settingsModal}/>
            </div>
        {/if}
        <div class="w-[20%] min-w-[220px] border-l-[1px] border-y-[1px] border-[#efdcc5] h-full bg-[#b4637a] rounded-l-lg flex flex-col">
            <Sidebar/>
        </div>
        <div class="w-[35%] min-w-[400px] border-y-[1px] border-[#efdcc5] h-full bg-[#debdc7] flex flex-col">
            <Search {search}/>
            {#if browser && !greeting}
                <button transition:blur={{duration: 10}} on:click={() => $settings.logs.announcement = true}
                        class="w-[136px] h-[42px] bg-[#faf4ed] text-[#413d5b] shadow-sm rounded-xl border-[#b4637a] border-[1px] py-1.5 px-4 mx-auto mb-5 my-auto flex justify-center items-center">
                    <span class="font-medium mr-1 text-sm">hewwo ^.^ </span>
                    <img alt="avatar" class="inline h-7"
                         src="https://cdn.discordapp.com/emojis/667830310741737473.gif?size=240&quality=lossless"/>
                </button>
            {/if}
        </div>
        <div class="w-[45%] min-w-[500px] border-y-[1px] border-r-[1px] border-[#efdcc5] h-full bg-[#e7cfd6] rounded-r-lg">
            <EncounterList {search}/>
        </div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
        -moz-osx-font-smoothing: grayscale;
    }
</style>