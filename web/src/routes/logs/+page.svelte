<script lang="ts">
    import type { PageData } from "./$types";
    import { blur } from "svelte/transition";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import Search from "$lib/components/Search.svelte";
    import EncounterList from "$lib/components/EncounterList.svelte";
    import Settings from "$lib/components/Settings.svelte";
    import { settings, user } from "$lib/store";
    import { settingsUI } from "$lib/menu";
    import { browser } from "$app/environment";
    import { writable } from "svelte/store";

    export let data: PageData;
    console.log(data);
    if (data.me.id) {
        user.set(data.me);
    }

    let settingsModal;
    settingsUI.subscribe(async (open) => {
        if (open) {
            let resp = await fetch("/api/settings", {
                credentials: "same-origin"
            });
            if (!resp.ok) {
                return;
            }

            settingsModal = await resp.json();
        } else {
            settingsModal = null;
        }
    });

    $: greeting = $settings?.logs.announcement;

    const search = writable({});
</script>

<svelte:head>
    <title>black meowket (new and improved)</title>
</svelte:head>

<div class="flex flex-row shadow rounded-lg max-w-[100rem] mx-auto h-screen">
    {#if browser && greeting}
        <div
            class="absolute text-[#413d5b] py-4 px-5 top-1/2 left-1/2 -translate-x-1/2 -translate-y-full shadow-md z-50 border-tapestry-300 border rounded-md w-96 bg-tapestry-50 dark:bg-tapestry-600 dark:text-white"
        >
            <div class="h-fit">
                <span class="font-semibold text-lg">hewwo ^.^ </span>
                <img
                    alt="avatar"
                    class="inline h-8 -translate-y-1"
                    src="https://cdn.discordapp.com/emojis/667830310741737473.gif?size=240&quality=lossless"
                />
                <button
                    on:click={() => ($settings.logs.announcement = false)}
                    class="float-right font-semibold"
                >[x]
                </button>
            </div>
            <p>
                you've stumbled upon... my website! this is a secret underground black market that deals
                in the forbidden teachings. FRICC THE POPO.
            </p>
            <p class="text-[#57517a] text-sm dark:text-gray-300">
                (sign up for the <a href="https://docs.fau.dev/logs/alpha/" class="underline hover:text-blue-500" target="_blank">alpha</a>)
            </p>
        </div>
    {/if}
    <!--{#if settingsModal}-->
    <!--    <div-->
    <!--        transition:blur={{ duration: 10 }}-->
    <!--        class="absolute text-[#413d5b] py-4 px-5 left-1/2 top-1/2 shadow-sm z-50 border-[#b4637a] border-[1px] rounded-xl -translate-x-[50%] -translate-y-[80%] w-[376px] h-[330px] bg-[#faf4ed]"-->
    <!--    >-->
    <!--        <Settings settings={settingsModal} />-->
    <!--    </div>-->
    <!--{/if}-->

    <Sidebar />

<!--    <div-->
<!--        class="w-[35%] min-w-[400px] border-y-[1px] border-[#efdcc5] h-full bg-[#debdc7] flex flex-col"-->
<!--    >-->
<!--        <Search {search} />-->
<!--    </div>-->
    <div
        class="bg-tapestry-100 w-full min-w-[30rem]"
    >
        <EncounterList {search} point={data.point} />
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
        -moz-osx-font-smoothing: grayscale;
    }
</style>
