<script lang="ts">
    import { settings, user } from "$lib/store";
    import { browser } from "$app/environment";

    import IconArrow from "~icons/carbon/next-filled";
    import IconBack from "~icons/ion/arrow-back-outline";
    import IconScope from "~icons/mdi/telescope";
    import IconSort from "~icons/basil/sort-outline";
    import EncounterRecap from "$lib/components/EncounterRecap.svelte";
    import { onMount } from "svelte";
    import EncounterPreview from "$lib/components/EncounterPreview.svelte";

    export let search;

    let encounters = [];
    let focused;

    function focus(encounter) {
        focused = encounter;
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

    let order = Order.RecentClear;
    $: scoped = browser && $settings.logs.scope;
    $: display = focused ? [focused] : encounters.slice(page * 5, (page + 1) * 5);
</script>

<link
    rel="preload"
    as="image"
    href="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless" />
<div class="mx-auto mt-10 flex max-w-3xl flex-col items-center justify-center">
    <div class="mb-3 flex flex-row items-center justify-center">
        <!-- <div class="rounded bg-white opacity-0">
            <IconSort class="h-6 w-6" />
        </div> -->
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
                    class="border-bouquet-600 text-bouquet-950 absolute z-50 float-right overflow-hidden rounded-md border bg-white text-sm shadow-sm">
                    <div class="bg-tapestry-600 text-center text-white">Sort</div>
                    {#each [Order.RecentClear, Order.RecentLog, Order.Duration, Order.Performance] as sort, i}
                        <button
                            class:underline={order === sort || (!order && i === 0)}
                            on:click={() => changeSort(sort)}
                            class="mx-auto my-0.5 whitespace-nowrap px-1 text-center hover:underline"
                            >{sort}</button>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
    <div
        class="mb-3 flex min-h-28 w-[88%] flex-col items-center justify-center overflow-hidden rounded-md bg-tapestry-300 pt-4">
        {#if display.length === 0}
            {#if loading}
                <div class="flex flex-col items-center justify-center -mt-3">
                    <span class="text-semibold text-sm text-bouquet-950">Loading...</span>
                    <img
                        alt="loading"
                        class="size-10"
                        src="https://cdn.discordapp.com/emojis/1056373578733461554.gif?size=240&quality=lossless" />
                </div>
            {:else if (scoped !== "Arkesia" || ($search.search && $search.search !== "")) && !loggedIn()}
                <div class="flex flex-col items-center justify-center -mt-3">
                    <span class="text-semibold mb-0.5 text-sm text-bouquet-950">Not signed in.</span>
                    <img
                        alt="loading"
                        class="size-10"
                        src="https://cdn.discordapp.com/attachments/1154431161993535489/1177165751040360448/emoji_a_38.png?ex=65718409&is=655f0f09&hm=cb2e683112d257a9d89dcc7fc90a54b4a91d73ddf67c0b3e1fd6df225fbff4f6&" />
                </div>
            {:else}
                <div class="flex flex-col items-center justify-center -mt-3">
                    <span class="text-semibold mb-0.5 text-sm text-bouquet-950">No logs found.</span>
                    <img
                        alt="loading"
                        class="size-10"
                        src="https://cdn.discordapp.com/emojis/987954898094129172.webp?size=240&quality=lossless" />
                </div>
            {/if}
        {/if}
        {#each display as encounter}
            <!-- <button
                class="mb-4 h-20 w-[95%]"
                on:click={() => focus(focused ? null : encounter)}> -->
            <button
                class="mb-4 h-20 w-[95%]">
                <EncounterPreview {encounter} />
            </button>
        {/each}
        <!-- {#if focused}
            <div class="mt-0.5 flex w-[88%] flex-row text-sm">
                <div class="my-1">
                    <button
                        on:click={() => focus(null)}
                        class="flex items-center justify-center rounded-md border-[0.5px] border-[#b4637a] bg-[#f5efec] p-0.5 px-1.5 text-[#b4637a]">
                        <IconBack class="mr-0.5 inline" />
                        Back
                    </button>
                </div>
                <div class="mx-auto mb-1 mt-auto rounded-md bg-[#b96d83] p-0.5 px-6 text-[#f7f2ef]">
                    Preview
                </div>
                <div class="my-1">
                    <button
                        class="flex items-center justify-center rounded-md border-[0.5px] border-[#575279] bg-[#f3eeec] p-0.5 px-1.5 text-[#575279]"
                        on:click={() => window.open("/log/" + focused.id, "_blank").focus()}>
                        <IconScope class="mr-0.5 inline" />
                        Open
                    </button>
                </div>
            </div>
            <EncounterRecap {focused} />
        {/if} -->
    </div>
    <!-- {#if !focused} -->
        <div class="flex flex-row items-center justify-center text-tapestry-600 space-x-10">
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
