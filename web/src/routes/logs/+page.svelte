<script lang="ts">
    import type {PageData} from './$types';
    import Sidebar from '$lib/components/Sidebar.svelte';
    import EncounterList from "$lib/components/EncounterList.svelte";
    import Settings from '$lib/components/Settings.svelte'
    import IconSearch from '~icons/gala/search'
    import {user} from '$lib/store'
    import {settingsUI} from '$lib/menu'

    export let data: PageData;
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

    let greeting = false;
</script>

<svelte:head>
    <title>logs by faust</title>
</svelte:head>

<div class="w-screen h-screen min-w-[1512px] min-h-[860px] bg-[#faf4ed] flex justify-center items-center">
    <div class="w-[1300px] h-[740px] flex flex-row shadow rounded-lg absolute">
        {#if greeting}
            <div class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-md -translate-x-[50%] -translate-y-[100%] w-[376px] h-[154px] bg-[#faf4ed]">
                <div class="h-fit">
                    <span class="font-semibold text-lg">hewwo ^.^ </span>
                    <img class="inline w-7 -translate-y-1"
                         src="https://cdn.discordapp.com/emojis/667830310741737473.gif?size=240&quality=lossless"/>
                    <button on:click={() => greeting = false} class="float-right font-semibold">[x]</button>
                </div>
                <p>you've stumbled upon... my website! this is a secret underground black market that deals in the
                    forbidden teachings. FRICC THE POPO.</p>
                <p class="text-[#57517a] text-sm">(invite only atm)</p>
            </div>
        {/if}
        {#if settings}
            <div class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-xl -translate-x-[50%] -translate-y-[80%] w-[376px] h-[330px] bg-[#faf4ed]">
                <Settings {settings}/>
            </div>
        {/if}
        <div class="w-[20%] min-w-[220px] border-l-[1px] border-y-[1px] border-[#efdcc5] h-full bg-[#b4637a] rounded-l-lg flex flex-col">
            <Sidebar/>
        </div>
        <div class="w-[35%] min-w-[400px] border-y-[1px] border-[#efdcc5] h-full bg-[#debdc7] flex flex-col">
            <div class="w-4/5 border-[1px] border-[#efdcc5] h-[46px] text-[#575279] mx-auto mt-10 rounded-2xl bg-[#f2e9e7] flex justify-center items-center flex-row">
                <input
                        placeholder="owo owo owo" class="bg-[#f2e9e7] placeholder-[#a7a3c1] w-4/5 outline-0"
                        autocomplete="off"
                        autocapitalize="off"
                        spellcheck="false"
                />
                <button>
                    <IconSearch class="w-6 h-6"/>
                </button>
            </div>
        </div>
        <div class="w-[45%] min-w-[500px] border-y-[1px] border-r-[1px] border-[#efdcc5] h-full bg-[#e7cfd6] rounded-r-lg">
            <EncounterList encounters={data.recent}/>
        </div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
    }
</style>