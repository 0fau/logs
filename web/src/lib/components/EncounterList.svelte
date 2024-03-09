<script lang="ts">
    import { settings, user } from "$lib/store";
    import { browser } from "$app/environment";

    import IconArrow from "~icons/carbon/next-filled";
    import IconBack from "~icons/ion/arrow-back-outline";
    import IconScope from "~icons/mdi/telescope";
    import IconSort from "~icons/basil/sort-outline";
    import IconMenu from '~icons/mingcute/menu-fill';

    import EncounterRecap from "$lib/components/EncounterRecap.svelte";
    import { onMount } from "svelte";
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";

    export let search;

    let encounters = [];
    let focused;

    function focus(encounter) {
        focused = encounter;
        if (encounter) {
            const el = document.getElementById("#" + encounter.id);
            if (!el) return;
            el.scrollIntoView({
                behavior: "smooth"
            });
        }
    }

    let page = 0;
    let more = false;
    let busy = false;

    let loading = true;
    onMount(async () => {
        busy = true;
        await load();
        console.log(encounters);
        busy = false;

        loading = false;

        search.subscribe((val) => {
            refresh();
        });
    });

    function loggedIn() {
        return $user && $user.id;
    }

    function prev() {
        if (busy) {
            return;
        }
        busy = true;
        if (page > 0) {
            page--;
        }
        busy = false;
    }

    function normalizeSelections(selected) {
        let selections = {
            raids: {},
            guardians: Object.keys(selected.guardians ?? {}).filter((g) => selected.guardians[g]),
            trials: Object.keys(selected.trials ?? {}).filter((t) => selected.trials[t]),
            classes: Object.keys(selected.classes ?? {}).filter((c) => selected.classes[c]),
            search: selected.search
        };
        for (let [key, val] of Object.entries(selected.raids ?? {})) {
            console.log(key, val);
            selections.raids[key] = {
                gates: Object.keys(val.gates)
                    .filter((g) => val.gates[g])
                    .map((v) => Number(v))
                    .sort(),
                difficulties: Object.keys(val.difficulties).filter((d) => val.difficulties[d])
            };
        }
        return selections;
    }

    export let point;

    async function load() {
        if ((scoped !== "Arkesia" || ($search.search && $search.search !== "")) && !loggedIn()) {
            return;
        }

        let url = location.protocol + "//" + location.host;
        if (point) {
            url = "https://logs.fau.dev";
        }

        url += "/api/logs?scope=" + scoped.toLowerCase();
        if (encounters.length > 0) {
            let last = encounters[encounters.length - 1];
            url += "&past_id=" + last.id;
            url += "&past_place=" + last.place;
            if (order === Order.Performance) {
                url += "&past_field=" + last.players[last.localPlayer].dps;
            } else if (order === Order.RecentClear) {
                url += "&past_field=" + last.date;
            } else if (order === Order.Duration) {
                url += "&past_field=" + last.duration;
            }
        }
        url += "&order=" + order.toLowerCase();

        if (gearScore) {
            url += "&gear_score=" + encodeURIComponent(gearScore)
        }

        const recent = await fetch(url, {
            credentials: "same-origin",
            method: "POST",
            body: JSON.stringify(normalizeSelections($search))
        }).then((resp) => {
            return resp.json();
        });
        process(recent.encounters);
        encounters = encounters.concat(recent.encounters);
        more = recent.more;
    }

    async function next() {
        if (busy) {
            return;
        }
        busy = true;
        await proceed();
        busy = false;
    }

    async function proceed() {
        if ((page + 1) * 5 < encounters.length) {
            page++;
            return;
        }

        if (!more) {
            return;
        }

        await load();

        if ((page + 1) * 5 < encounters.length) {
            page++;
            return;
        }
    }

    function process(encounters) {
        for (let encounter of encounters) {
            let names = Object.keys(encounter.players);
            let max = encounter.players[names[0]].damage;

            for (let name of names) {
                if (encounter.players[name].damage > max) {
                    max = encounter.players[name].damage;
                }
            }

            encounter.max = max;
        }
    }

    async function refresh() {
        focused = null;
        more = false;
        encounters = [];
        page = 0;
        loading = true;
        await load();
        loading = false;
    }

    async function changeScope(scope) {
        $settings.logs.scope = scope;
        scoped = scope;
        await refresh();
    }

    async function changeSort(sort) {
        showSortOptions = false;
        order = sort;
        await refresh();
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
        };

        document.addEventListener("click", click, true);
        return {
            destroy() {
                document.removeEventListener("click", click, true);
            }
        };
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
    $: scoped = browser && $settings.logs.scope;
    $: display = encounters.slice(page * 5, (page + 1) * 5);
</script>

<link
    rel="preload"
    as="image"
    href="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless" />
<div class="mx-auto mt-10 flex max-h-screen max-w-3xl flex-col items-center">
    <div class="mb-3 flex items-center">
        <button class="block md:hidden absolute left-[7%]">
            <IconMenu class="size-6 text-tapestry-600 hover:text-tapestry-700" />
        </button>
        <div class="flex flex-row items-center justify-center">
            <div
                class="flex h-10 w-60 flex-row items-center justify-center rounded-xl bg-tapestry-600 text-center text-sm text-tapestry-50">
                {#each ["Arkesia", "Friends", "Roster"] as scope}
                    <div
                        class="flex h-full w-full items-center justify-center"
                        class:rounded-lg={browser && scoped === scope}>
                        <div
                            class="flex h-[76%] w-[87%] items-center justify-center rounded-lg"
                            class:bg-tapestry-100={browser && scoped === scope}>
                            <button
                                class:font-medium={browser && scoped === scope}
                                class:text-tapestry-700={browser && scoped === scope}
                                on:click={() => changeScope(scope)}>
                                {scope}
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
            <div class="relative ml-2 size-6 rounded-md bg-tapestry-600 text-tapestry-50">
                <button
                    bind:this={toggle}
                    on:click={() => (showSortOptions = !showSortOptions)}
                    class="h-full w-full">
                    <IconSort class="size-6" />
                </button>
                {#if showSortOptions}
                    <div
                        use:unfocus
                        class="absolute z-50 float-right overflow-hidden rounded-md border border-bouquet-600 bg-white text-sm text-bouquet-950 shadow-sm">
                        <div class="bg-tapestry-600 text-center text-white">Sort</div>
                        {#each [Order.RecentClear, Order.RecentLog, Order.Duration, Order.Performance] as sort, i}
                            <button
                                class:underline={order === sort || (!order && i === 0)}
                                on:click={() => changeSort(sort)}
                                class="mx-auto my-0.5 whitespace-nowrap px-1 text-center hover:underline"
                                >{sort}</button>
                        {/each}
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
    </div>
    <div
        class="custom-scroll mb-3 flex w-[88%] flex-col items-center overflow-y-scroll rounded-md bg-tapestry-300 pl-2 pt-4"
        style="height: calc(100vh - 10rem)">
        {#if display.length === 0}
            {#if loading}
                <div class="mt-3 flex flex-col items-center justify-center">
                    <span class="text-semibold text-sm text-bouquet-950">Loading...</span>
                    <img
                        alt="loading"
                        class="size-10"
                        src="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless" />
                </div>
            {:else if (scoped !== "Arkesia" || ($search.search && $search.search !== "")) && !loggedIn()}
                <div class="mt-3 flex flex-col items-center justify-center">
                    <span class="text-semibold mb-0.5 text-sm text-bouquet-950"
                        >Not signed in.</span>
                    <img
                        alt="loading"
                        class="size-10"
                        src="https://cdn.discordapp.com/attachments/1154431161993535489/1177165751040360448/emoji_a_38.png?ex=65718409&is=655f0f09&hm=cb2e683112d257a9d89dcc7fc90a54b4a91d73ddf67c0b3e1fd6df225fbff4f6&" />
                </div>
            {:else}
                <div class="mt-3 flex flex-col items-center justify-center">
                    <span class="text-semibold mb-0.5 text-sm text-bouquet-950"
                        >No logs found.</span>
                    <img
                        alt="loading"
                        class="size-10"
                        src="https://cdn.discordapp.com/emojis/987954898094129172.webp?size=240&quality=lossless" />
                </div>
            {/if}
        {/if}
        {#each display as encounter}
            <button
                class="mb-2 h-20 w-full px-2"
                id={"#" + encounter.id}
                on:click={() => focus(focused && focused.id === encounter.id ? null : encounter)}>
                <EncounterPreview {encounter} />
            </button>
            {#if focused && focused.id === encounter.id}
                <div class="mb-2 flex w-full justify-between px-5 text-sm">
                    <!-- <div class="my-1">
                        <button
                            on:click={() => focus(null)}
                            class="flex items-center justify-center rounded-md border-[0.5px] border-[#b4637a] bg-[#f5efec] p-0.5 px-1.5 text-[#b4637a]">
                            <IconBack class="mr-0.5 inline" />
                            Back
                        </button>
                    </div> -->
                    <div class="mt-auto rounded-md bg-tapestry-600 p-0.5 px-6 text-white">
                        Preview
                    </div>
                    <a
                        class="flex items-center justify-center rounded-md border border-gray-600 bg-tapestry-50 p-0.5 px-1.5 text-gray-700"
                        href={"/log/" + focused.id}
                        target="_blank">
                        <IconScope class="mr-0.5 inline" />
                        Open
                    </a>
                </div>
                <EncounterRecap {focused} />
            {/if}
        {/each}
    </div>
    <!-- {#if !focused} -->
    <div class="flex flex-row items-center justify-center space-x-10 text-tapestry-600">
        {#if page > 0}
            <button
                on:click={prev}
                class="rounded-3xl border border-tapestry-400 bg-tapestry-50 p-0.5 shadow-sm">
                <IconArrow class="size-6 rotate-180" />
            </button>
        {/if}
        {#if more || (page + 1) * 5 < encounters.length}
            <button
                on:click={next}
                class="rounded-3xl border border-tapestry-400 bg-tapestry-50 p-0.5 shadow-sm">
                <IconArrow class="size-6" />
            </button>
        {/if}
    </div>
    <!-- {/if} -->
</div>


<style lang="postcss">
    /* * {
        -ms-overflow-style: none;
        scrollbar-width: none;
    }
    *::-webkit-scrollbar {
        display: none;
    } */

    .custom-scroll::-webkit-scrollbar {
        @apply right-0 block size-2;
    }
    .custom-scroll::-webkit-scrollbar-thumb {
        @apply rounded-md bg-tapestry-600;
    }
    .custom-scroll::-webkit-scrollbar-corner {
        @apply bg-zinc-800;
    }
</style>
