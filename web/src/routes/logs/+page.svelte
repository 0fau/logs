<script lang="ts">
    import type {PageData} from './$types';
    import Sidebar from '$lib/components/Sidebar.svelte';
    import EncounterList from "$lib/components/EncounterList.svelte";
    import Settings from '$lib/components/Settings.svelte'
    import IconSearch from '~icons/gala/search'
    import {settings, user} from '$lib/store'
    import {settingsUI} from '$lib/menu'
    import {getRaidIcon} from "$lib/raids";
    import {browser} from "$app/environment"

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
                <p class="text-[#57517a] text-sm">(invite only atm)</p>
            </div>
        {/if}
        {#if settingsModal}
            <div class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-xl -translate-x-[50%] -translate-y-[80%] w-[376px] h-[330px] bg-[#faf4ed]">
                <Settings settings={settingsModal}/>
            </div>
        {/if}
        <div class="w-[20%] min-w-[220px] border-l-[1px] border-y-[1px] border-[#efdcc5] h-full bg-[#b4637a] rounded-l-lg flex flex-col">
            <Sidebar/>
        </div>
        <div class="w-[35%] min-w-[400px] border-y-[1px] border-[#efdcc5] h-full bg-[#debdc7] flex flex-col">
            <div class="w-4/5 border-[1px] border-[#efdcc5] h-[46px] text-[#575279] mx-auto mt-10 rounded-2xl bg-[#f2e9e7] flex justify-center items-center flex-row">
                <input
                        placeholder="owo owo owo"
                        class="bg-[#f2e9e7] placeholder-[#a7a3c1] w-4/5 outline-none outline-0"
                        autocomplete="off"
                        autocapitalize="off"
                        spellcheck="false"
                />
                <button>
                    <IconSearch class="w-6 h-6"/>
                </button>
            </div>
            <div class="mx-auto w-[78%] rounded-xl opacity-95 bg-[#f2e9e7] mt-2 flex flex-row">
                <div class="w-1/2 px-3 pt-3 py-2 h-full">
                    <button class="w-4/5 mb-1 h-[32px] hover:shadow-md transition flex flex-row rounded-xl bg-[#9db79d]">
                        <div class="w-fit h-full flex flex-row items-center px-1.5 rounded-l-xl bg-[#80a180]">
                            <img alt="Akkan" class="h-6 -translate-y-[0.25px] my-auto" src="{getRaidIcon('Akkan')}"/>
                        </div>
                    </button>
                    <button class="w-4/5 mb-1 h-[32px] hover:shadow-md transition rounded-xl bg-[#e0c261]">
                        <div class="w-fit h-full flex flex-row items-center px-1.5 rounded-l-xl bg-[#d8b236]">
                            <img alt="Kayangel" class="h-6 -translate-y-[0.25px] my-auto"
                                 src="{getRaidIcon('Kayangel')}"/>
                        </div>
                    </button>
                    <button class="w-4/5 mb-1 h-[32px] hover:shadow-md transition rounded-xl bg-[#b9acc6]">
                        <div class="w-fit h-full flex flex-row items-center px-1.5 rounded-l-xl bg-[#a08eb1]">
                            <img alt="Brelshaza" class="h-6 -translate-y-[0.25px] my-auto"
                                 src="{getRaidIcon('Brelshaza')}"/>
                        </div>
                    </button>
                </div>
                <div class="w-1/2 px-3 pt-3 py-2 h-full">
                    <button class="w-4/5 mb-1 h-[32px] hover:shadow-md transition rounded-xl bg-[#7ea6b2]">
                        <div class="w-fit h-full flex flex-row items-center px-1.5 rounded-l-xl bg-[#5f909e]">
                            <img alt="Valtan" class="h-6 -translate-y-[0.25px] my-auto" src="{getRaidIcon('Valtan')}"/>
                        </div>
                    </button>
                    <button class="w-4/5 mb-1 h-[32px] hover:shadow-md transition rounded-xl bg-[#d29eab]">
                        <div class="w-fit h-full flex flex-row items-center px-1.5 rounded-l-xl bg-[#c27b8d]">
                            <img alt="Vykas" class="h-6 -translate-y-[0.25px] my-auto" src="{getRaidIcon('Vykas')}"/>
                        </div>
                    </button>
                    <button class="w-4/5 mb-1 h-[32px] hover:shadow-md transition rounded-xl bg-[#d08e99]">
                        <div class="w-fit h-full flex flex-row items-center px-1.5 rounded-l-xl bg-[#c16a78]">
                            <img alt="Kakul Saydon" class="h-6 -translate-y-[0.25px] my-auto"
                                 src="{getRaidIcon('Kakul Saydon')}"/>
                        </div>
                    </button>
                </div>
            </div>
            {#if browser && !greeting}
                <button on:click={() => $settings.logs.announcement = true}
                        class="w-[136px] h-[42px] bg-[#faf4ed] text-[#413d5b] shadow-sm rounded-xl border-[#b4637a] border-[1px] py-1.5 px-4 mx-auto mb-5 my-auto flex justify-center items-center">
                    <span class="font-medium mr-1 text-sm">hewwo ^.^ </span>
                    <img alt="avatar" class="inline h-7"
                         src="https://cdn.discordapp.com/emojis/667830310741737473.gif?size=240&quality=lossless"/>
                </button>
            {/if}
        </div>
        <div class="w-[45%] min-w-[500px] border-y-[1px] border-r-[1px] border-[#efdcc5] h-full bg-[#e7cfd6] rounded-r-lg">
            <EncounterList />
        </div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
        -moz-osx-font-smoothing: grayscale;
    }
</style>