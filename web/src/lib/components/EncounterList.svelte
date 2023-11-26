<script lang="ts">
    import {settings, user} from "$lib/store";
    import {browser} from "$app/environment";

    import IconArrow from '~icons/carbon/next-filled'
    import IconBack from '~icons/ion/arrow-back-outline'
    import IconScope from '~icons/mdi/telescope'
    import EncounterRecap from "$lib/components/EncounterRecap.svelte";
    import {onMount} from "svelte";
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";

    let encounters = [];
    let focused;

    function focus(encounter) {
        focused = encounter
    }

    let page = 0;
    let more = false;
    let busy = false;

    let loading = true;
    onMount(async () => {
        busy = true;
        await load()
        console.log(encounters)
        busy = false;

        loading = false;
    })

    function loggedIn() {
        return $user && $user.id
    }

    function prev() {
        if (busy) {
            return
        }
        busy = true
        if (page > 0) {
            page--
        }
        busy = false
    }

    async function load() {
        if (scoped != "Arkesia" && !loggedIn()) {
            return
        }

        let url = location.protocol + '//' + location.host;
        url += "/api/logs?scope=" + scoped.toLowerCase();
        if (encounters.length > 0) {
            let last = encounters[encounters.length - 1];
            url += "&past=" + last.date + "&id=" + last.id;
        }

        const recent = await fetch(url, {credentials: 'same-origin'})
            .then((resp) => {
                return resp.json()
            })
        process(recent.encounters)
        encounters = encounters.concat(recent.encounters)
        more = recent.more
    }

    async function next() {
        if (busy) {
            return
        }
        busy = true
        await proceed()
        busy = false
    }

    async function proceed() {
        if ((page + 1) * 5 < encounters.length) {
            page++
            return
        }

        if (!more) {
            return
        }

        await load()

        if ((page + 1) * 5 < encounters.length) {
            page++
            return
        }
    }

    function process(encounters) {
        for (let encounter of encounters) {
            let names = Object.keys(encounter.players)
            let max = encounter.players[names[0]].damage

            for (let name of names) {
                if (encounter.players[name].damage > max) {
                    max = encounter.players[name].damage
                }
            }

            encounter.max = max
        }
    }

    async function changeScope(scope) {
        focused = null
        $settings.logs.scope = scope
        scoped = scope
        more = false
        encounters = []
        page = 0
        loading = true
        await load()
        loading = false
    }

    $: scoped = browser && $settings.logs.scope
    $: display = focused ? [focused] : encounters.slice(page * 5, (page + 1) * 5)
</script>

<link rel="preload" as="image"
      href="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless">
<div class="m-auto mt-10 flex flex-col justify-center items-center">
    <div class="flex flex-row w-[88%] justify-center items-center">
        <div class="w-[270px] h-[40px] bg-[#b96d83] text-center text-[#F4EDE9] text-sm flex flex-row justify-center items-center rounded-xl mb-3">
            {#each ["Arkesia", "Friends", "Roster"] as scope}
                <div class="w-full h-full flex justify-center items-center"
                     class:rounded-lg={browser && scoped === scope}>
                    <div class="w-[88%] h-[76%] rounded-lg flex justify-center items-center"
                         class:bg-[#F4EDE9]={browser && scoped === scope}>
                        <button class:font-medium={browser && scoped === scope}
                                class:text-[#b96d83]={browser && scoped === scope}
                                on:click={() => changeScope(scope)}>
                            {scope}
                        </button>
                    </div>
                </div>
            {/each}
        </div>
    </div>
    <div class="w-[88%] min-h-[110px] flex flex-col overflow-hidden justify-center items-center bg-[#dec5cd] pt-4 mb-3 rounded-md">
        {#if display.length === 0}
            {#if loading}
                <div class="flex flex-col items-center justify-center -translate-y-[14%]">
                    <span class="text-[#413d5b] text-sm text-semibold">Loading...</span>
                    <img alt="loading"
                         class="w-10 h-10"
                         src="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless"/>
                </div>
            {:else if scoped !== "Arkesia" && !loggedIn()}
                <div class="flex flex-col items-center justify-center -translate-y-[14%]">
                    <span class="text-[#413d5b] text-sm text-semibold mb-0.5">Not signed in.</span>
                    <img alt="loading"
                         class="w-10 h-10"
                         src="https://cdn.discordapp.com/attachments/1154431161993535489/1177165751040360448/emoji_a_38.png?ex=65718409&is=655f0f09&hm=cb2e683112d257a9d89dcc7fc90a54b4a91d73ddf67c0b3e1fd6df225fbff4f6&"/>
                </div>
            {:else}
                <div class="flex flex-col items-center justify-center -translate-y-[14%]">
                    <span class="text-[#413d5b] text-sm text-semibold mb-0.5">No logs found.</span>
                    <img alt="loading"
                         class="w-10 h-10"
                         src="https://cdn.discordapp.com/emojis/987954898094129172.webp?size=240&quality=lossless"/>
                </div>
            {/if}
        {/if}
        {#each display as encounter}
            <button
                    class="{focused ? '' : 'mb-5'} h-[80px] w-[94%]"
                    on:click={() => focus(focused ? null : encounter)}>
                <EncounterPreview width="w-full" {encounter}/>
            </button>
        {/each}
        {#if focused}
            <div class="w-[88%] mt-0.5 flex flex-row text-sm">
                <div class="my-1">
                    <button on:click={() => focus(null)}
                            class="flex items-center justify-center bg-[#f5efec] p-0.5 px-1.5 border-[0.5px] border-[#b4637a] rounded-md text-[#b4637a]">
                        <IconBack class="inline mr-0.5"/>
                        Back
                    </button>
                </div>
                <div class="mx-auto mt-auto mb-1 p-0.5 px-6 rounded-md text-[#f7f2ef] bg-[#b96d83]">Preview</div>
                <div class="my-1">
                    <button class="flex items-center justify-center bg-[#f3eeec] p-0.5 px-1.5 border-[0.5px] border-[#575279] rounded-md text-[#575279]"
                            on:click={() => window.open("/log/" + focused.id, '_blank').focus()}>
                        <IconScope class="inline mr-0.5"/>
                        Open
                    </button>
                </div>
            </div>
            <EncounterRecap {focused}/>
        {/if}
    </div>
    {#if !focused}
        <div class="flex flex-row text-[#b4637a] items-center justify-center">
            {#if page > 0}
                <button on:click={prev}
                        class="bg-[#F4EDE9] border-[0.5px] border-[#c58597] p-0.5 rounded-3xl shadow-sm">
                    <IconArrow class="rotate-180 w-6 h-6"/>
                </button>
            {/if}
            {#if more || (page + 1) * 5 < encounters.length}
                <button on:click={next}
                        class="ml-5 bg-[#F4EDE9] border-[0.5px] border-[#c58597] p-0.5 rounded-3xl shadow-sm">
                    <IconArrow class="w-6 h-6"/>
                </button>
            {/if}
        </div>
    {/if}
</div>