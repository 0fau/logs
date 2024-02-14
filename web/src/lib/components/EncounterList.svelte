<script lang="ts">
    import {settings, user} from "$lib/store";
    import {browser} from "$app/environment";

    import IconArrow from '~icons/carbon/next-filled'
    import IconBack from '~icons/ion/arrow-back-outline'
    import IconScope from '~icons/mdi/telescope'
    import IconSort from '~icons/basil/sort-outline'
    import EncounterRecap from "$lib/components/EncounterRecap.svelte";
    import {onMount} from "svelte";
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";

    export let search;

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

        search.subscribe((val) => {
            refresh()
        })
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

    function normalizeSelections(selected) {
        let selections = {
            raids: {},
            guardians: Object.keys(selected.guardians ?? {}).filter((g) => selected.guardians[g]),
            trials: Object.keys(selected.trials ?? {}).filter((t) => selected.trials[t]),
            classes: Object.keys(selected.classes ?? {}).filter((c) => selected.classes[c]),
            search: selected.search,
        }
        for (let [key, val] of Object.entries(selected.raids ?? {})) {
            console.log(key, val)
            selections.raids[key] = {
                gates: Object.keys(val.gates).filter((g) => val.gates[g]).map((v) => Number(v)).sort(),
                difficulties: Object.keys(val.difficulties).filter((d) => val.difficulties[d])
            }
        }
        return selections
    }

    export let point;

    async function load() {
        if ((scoped !== "Arkesia" || $search.search && $search.search !== "") && !loggedIn()) {
            return
        }

        let url = location.protocol + '//' + location.host;
        if (point) {
            url = "https://logs.fau.dev"
        }

        url += "/api/logs?scope=" + scoped.toLowerCase();
        if (encounters.length > 0) {
            let last = encounters[encounters.length - 1];
            url += "&past_id=" + last.id
            url += "&past_place=" + last.place
            if (order === Order.Performance) {
                url += "&past_field=" + last.players[last.localPlayer].dps;
            } else if (order === Order.RecentClear) {
                url += "&past_field=" + last.date
            } else if (order === Order.Duration) {
                url += "&past_field=" + last.duration
            }
        }
        url += "&order=" + order.toLowerCase()

        if (gearScore) {
            url += "&gear_score=" + encodeURIComponent(gearScore)
        }

        const recent = await fetch(url, {
            credentials: 'same-origin',
            method: 'POST',
            body: JSON.stringify(normalizeSelections($search))
        }).then((resp) => {
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

    async function refresh() {
        focused = null
        more = false
        encounters = []
        page = 0
        loading = true
        await load()
        loading = false
    }

    async function changeScope(scope) {
        $settings.logs.scope = scope
        scoped = scope
        await refresh()
    }

    async function changeSort(sort) {
        showSortOptions = false;
        order = sort
        await refresh()
    }

    async function changeGearScore(score) {
        showSortOptions = false;
        if (gearScore !== score) {
            gearScore = score
        } else {
            gearScore = null
        }
        await refresh()
    }

    let showSortOptions = false;

    let toggle;

    function unfocus(box) {
        const click = (event) => {
            if (!box.contains(event.target) && !toggle.contains(event.target)) {
                showSortOptions = false;
            }
        }

        document.addEventListener('click', click, true)
        return {
            destroy() {
                document.removeEventListener('click', click, true)
            }
        }
    }

    enum Order {
        RecentClear = "Recent Clear",
        RecentLog = "Recent Log",
        Duration = "Raid Duration",
        Performance = "Performance"
    }

    enum GearScoreRange {
        Range1540To1560 = "1540-1560",
        Range1560To1580 = "1560-1580",
        Range1580To1600 = "1580-1600",
        Range1600To1610 = "1600-1610",
        Range1610To1620 = "1610-1620",
        Range1620Plus = "1620+"
    }

    const gearScores = [
        GearScoreRange.Range1540To1560,
        GearScoreRange.Range1560To1580,
        GearScoreRange.Range1580To1600,
        GearScoreRange.Range1600To1610,
        GearScoreRange.Range1610To1620,
        GearScoreRange.Range1620Plus
    ]

    let gearScore;

    let order = Order.RecentClear;
    $: scoped = browser && $settings.logs.scope
    $: display = focused ? [focused] : encounters.slice(page * 5, (page + 1) * 5)
</script>

<link rel="preload" as="image"
      href="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless">
<div class="m-auto mt-10 flex flex-col justify-center items-center">
    <div class="flex mb-3 flex-row w-[88%] justify-center items-center">
        <div class="bg-white border-[0.5px] mr-1.5 rounded opacity-0">
            <IconSort class="w-6 h-6"/>
        </div>
        <div class="w-[270px] h-[40px] bg-[#b96d83] text-center text-[#F4EDE9] text-sm flex flex-row justify-center items-center rounded-xl">
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
        <div class="rounded-md h-6 w-6 ml-1.5 bg-[#a7738b] border-[0.5px] border-[#e8d2d7] text-[#F4EDE9] relative">
            <button bind:this={toggle} on:click={() => showSortOptions = !showSortOptions} class="w-full h-full">
                <IconSort class="w-6 h-6"/>
            </button>
            {#if showSortOptions}
                <div use:unfocus transition:blur={{duration: 10}}
                     class="absolute overflow-hidden float-right text-[#575279] text-sm z-50">
                    <div class="bg-[#F4EDE9] border-[#a7738b] border-[0.5px] shadow-sm rounded-md overflow-hidden">
                        <div class="text-center text-[#F4EDE9] bg-[#a7738b]">Sort</div>
                        {#each [Order.RecentClear, Order.RecentLog, Order.Duration, Order.Performance] as sort, i}
                            <button class:underline={order === sort || (!order && i === 0)}
                                    on:click={() => changeSort(sort)}
                                    class="whitespace-nowrap mx-auto text-center px-1 my-0.5">{sort}</button>
                        {/each}
                    </div>
                    <div class="mt-1 bg-[#F4EDE9] border-[#a7738b] border-[0.5px] shadow-sm rounded-md overflow-hidden">
                        <div class="text-center text-[#F4EDE9] bg-[#a7738b]">Gear Score</div>
                        {#each gearScores as range, i}
                            <button class:underline={gearScore === range}
                                    on:click={() => changeGearScore(range)}
                                    class="whitespace-nowrap mx-auto text-center px-1 my-0.5">{range}</button>
                        {/each}
                    </div>
                </div>
            {/if}
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
            {:else if (scoped !== "Arkesia" || $search.search && $search.search !== "") && !loggedIn()}
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
                <EncounterPreview gearScore={gearScore} width="w-full" {encounter}/>
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