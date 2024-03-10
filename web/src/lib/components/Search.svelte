<script>
    import { getRaidIcon } from "$lib/raids";
    import { classes, getClassIconNew } from "$lib/game";
    import IconNumberOne from "~icons/mdi/number-1-box";
    import IconNumberTwo from "~icons/mdi/number-2-box";
    import IconNumberThree from "~icons/mdi/number-3-box";
    import IconNumberFour from "~icons/mdi/number-4-box";
    import IconNumberFive from "~icons/mdi/number-5-box";
    import IconNumberSix from "~icons/mdi/number-6-box";
    import IconN from "~icons/mdi/letter-n-box";
    import IconH from "~icons/mdi/letter-h-box";
    import IconE from "~icons/mdi/letter-e-box";
    import IconI from "~icons/heroicons-solid/fire";
    import IconNumberOneO from "~icons/mdi/number-1-box-outline";
    import IconNumberTwoO from "~icons/mdi/number-2-box-outline";
    import IconNumberThreeO from "~icons/mdi/number-3-box-outline";
    import IconNumberFourO from "~icons/mdi/number-4-box-outline";
    import IconNumberFiveO from "~icons/mdi/number-5-box-outline";
    import IconNumberSixO from "~icons/mdi/number-6-box-outline";
    import IconNO from "~icons/mdi/letter-n-box-outline";
    import IconHO from "~icons/mdi/letter-h-box-outline";
    import IconEO from "~icons/mdi/letter-e-box-outline";
    import IconIO from "~icons/heroicons-outline/fire";
    import IconFrog from "~icons/icon-park-outline/frog";
    import IconTree from "~icons/typcn/tree";
    import IconMonkey from "~icons/emojione-monotone/monkey-face";
    import IconDragon from "~icons/fluent-emoji-high-contrast/dragon-face";
    import IconBird from "~icons/lucide/bird";
    import IconLion from "~icons/emojione-monotone/lion-face";
    import IconScaryDragon from "~icons/game-icons/spiked-dragon-head";
    import IconSearch from "~icons/gala/search";
    import IconFilter from "~icons/mi/filter";
    import { writable } from "svelte/store";

    let numbers = [
        IconNumberOne,
        IconNumberTwo,
        IconNumberThree,
        IconNumberFour,
        IconNumberFive,
        IconNumberSix
    ];

    let numbersO = [
        IconNumberOneO,
        IconNumberTwoO,
        IconNumberThreeO,
        IconNumberFourO,
        IconNumberFiveO,
        IconNumberSixO
    ];

    let difficulties = {
        Normal: IconN,
        Hard: IconH,
        Inferno: IconI,
        Extreme: IconE
    };

    let difficultiesO = {
        Normal: IconNO,
        Hard: IconHO,
        Inferno: IconIO,
        Extreme: IconEO
    };

    let raids = [
        [
            {
                name: "Ivory",
                primary: "bg-[#8fa6ad]",
                secondary: "bg-[#759098]",
                gates: 4,
                difficulties: ["Normal", "Hard"]
            },
            {
                name: "Akkan",
                primary: "bg-[#9db79d]",
                secondary: "bg-[#80a180]",
                gates: 3,
                difficulties: ["Normal", "Hard"]
            },
            {
                name: "Kayangel",
                primary: "bg-[#e0c261]",
                secondary: "bg-[#d8b236]",
                gates: 3,
                difficulties: ["Normal", "Hard"]
            },
            {
                name: "Brelshaza",
                primary: "bg-[#b9acc6]",
                secondary: "bg-[#a08eb1]",
                gates: 6,
                difficulties: ["Normal", "Hard", "Inferno"]
            }
        ],
        [
            {
                name: "Valtan",
                primary: "bg-[#7ea6b2]",
                secondary: "bg-[#5f909e]",
                gates: 2,
                difficulties: ["Normal", "Hard", "Extreme", "Inferno"]
            },
            {
                name: "Vykas",
                primary: "bg-[#d29eab]",
                secondary: "bg-[#c27b8d]",
                gates: 3,
                difficulties: ["Normal", "Hard", "Inferno"]
            },
            {
                name: "Kakul Saydon",
                primary: "bg-[#d08e99]",
                secondary: "bg-[#c16a78]",
                gates: 3,
                difficulties: ["Normal", "Inferno"]
            }
        ]
    ];

    let raidDropdown = "";

    let list = Object.keys(classes);
    list.sort();

    let guardians = [
        {
            name: "Gargadeth",
            icon: IconFrog
        },
        {
            name: "Sonavel",
            icon: IconTree
        },
        {
            name: "Hanumatan",
            icon: IconMonkey
        },
        {
            name: "Caliligos",
            icon: IconDragon
        },
        {
            name: "Deskaluda",
            icon: IconBird
        }
    ];

    let trials = [
        {
            name: "Caliligos",
            icon: IconScaryDragon
        },
        {
            name: "Achates",
            icon: IconLion
        }
    ];

    let shortcuts = [
        { name: "Endgame", color: "#b96d83" },
        { name: "Inferno", color: "#b96d83" }
    ];

    let selectedRaids = {};
    let selectedGuardians = {};
    let selectedTrials = {};
    let selectedClasses = {};

    function selectGate(raid, gate) {
        let selected = selectedRaids[raid] ?? { difficulties: {}, gates: {} };
        selected.gates[gate] = !selected.gates[gate];
        selectedRaids[raid] = selected;

        console.log(selectedRaids);

        update();
    }

    function selectDifficulty(raid, difficulty) {
        let selected = selectedRaids[raid] ?? { difficulties: {}, gates: {} };
        selected.difficulties[difficulty] = !selected.difficulties[difficulty];
        selectedRaids[raid] = selected;

        console.log(selectedRaids);

        update();
    }

    let task;
    function update() {
        clearTimeout(task);
        task = setTimeout(() => {
            search.update(() => ({
                raids: selectedRaids,
                guardians: selectedGuardians,
                trials: selectedTrials,
                classes: selectedClasses,
                search:
                    searching !== "" ? searching.charAt(0).toUpperCase() + searching.slice(1) : ""
            }));
        }, 500);
    }
    export let search;
    let searching = "";
    let hovered;
    let showFilter = writable(false);

</script>

<div class="relative mb-5">
    <div
        class="mx-auto flex h-[46px] w-4/5 flex-row items-center justify-center rounded-2xl border border-[#efdcc5] bg-[#f2e9e7] text-[#575279]">
        <button class="pl-2 pr-1" on:click={() => $showFilter = !$showFilter}>
            <IconFilter class="h-6 w-6" />
        </button>
        <input
            bind:value={searching}
            on:input={update}
            placeholder="owo owo owo"
            class="w-4/5 bg-[#f2e9e7] placeholder-[#a7a3c1] outline-none outline-0"
            autocomplete="off"
            autocapitalize="off"
            spellcheck="false" />
<!--        <button class="pl-1 pr-2">-->
<!--            <IconSearch class="h-6 w-6" />-->
<!--        </button>-->
    </div>
    {#if $showFilter}
        <div class="absolute bg-tapestry-200 z-50 rounded-xl mt-4 border border-tapestry-500 p-4 sm:w-[30rem] -translate-x-[10%] sm:-translate-x-1/4">
        <div
            class="mx-auto mt-2 flex w-[78%] flex-col items-center rounded-xl border-[0.75px] border-[#c17e91] bg-[#f2e9e7] opacity-95">
            <div class="flex w-full flex-row">
                {#each raids as column}
                    <div class="h-full w-1/2 px-3 py-2 pt-3">
                        {#each column as raid}
                            <div>
                                <button
                                    on:click={() =>
                                        (raidDropdown = raidDropdown === raid ? null : raid)}
                                    class="relative mb-1 flex h-[32px] w-4/5 flex-row justify-center rounded-xl transition hover:shadow-md {raid.primary}">
                                    <div
                                        class="ml-0 mr-auto flex h-full w-[36px] flex-row items-center rounded-l-xl px-1.5 {raid.secondary}">
                                        <img
                                            alt={raid.name}
                                            class="my-auto h-6 -translate-y-[0.25px]"
                                            src={getRaidIcon(raid.name)} />
                                    </div>
                                </button>
                            </div>
                            <div
                                class:h-fit={raidDropdown === raid}
                                class:p-2={raidDropdown === raid}
                                class:mb-1={raidDropdown === raid}
                                class:h-0={raidDropdown !== raid}
                                class="w-4/5 overflow-hidden rounded-xl {raid.primary}">
                                <div class="flex h-full w-full flex-col items-center">
                                    <div class="grid grid-cols-4">
                                        {#each { length: raid.gates } as _, i}
                                            {@const selected = selectedRaids[raid.name]
                                                ? selectedRaids[raid.name].gates[i + 1]
                                                : false}
                                            <button on:click={() => selectGate(raid.name, i + 1)}>
                                                {#if selected}
                                                    <svelte:component
                                                        this={numbers[i]}
                                                        class="h-6 w-6 text-[#f8f4f3]" />
                                                {:else}
                                                    <svelte:component
                                                        this={numbersO[i]}
                                                        class="h-6 w-6 text-[#f8f4f3]" />
                                                {/if}
                                            </button>
                                        {/each}
                                    </div>
                                    <div class="grid grid-cols-4">
                                        {#each raid.difficulties as difficulty}
                                            {@const selected = selectedRaids[raid.name]
                                                ? selectedRaids[raid.name].difficulties[difficulty]
                                                : false}
                                            <button
                                                on:click={() => {
                                                    selectDifficulty(raid.name, difficulty);
                                                }}>
                                                {#if selected}
                                                    <svelte:component
                                                        this={difficulties[difficulty]}
                                                        class="h-6 w-6 text-[#f8f4f3]" />
                                                {:else}
                                                    <svelte:component
                                                        this={difficultiesO[difficulty]}
                                                        class="h-6 w-6 text-[#f8f4f3]" />
                                                {/if}
                                            </button>
                                        {/each}
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/each}
            </div>
        </div>
        <div class="mx-auto mt-2 flex w-[78%] flex-row items-center justify-center">
            <div
                class="mr-2 flex flex-row items-center rounded-xl border-[0.75px] border-[#c17e91] bg-[#f2e9e7] p-1.5 opacity-95">
                {#each guardians as guardian}
                    {@const hoverkey = "guardian_" + guardian.name}
                    <button
                        on:click={() => {
                            selectedGuardians[guardian.name] = !selectedGuardians[guardian.name];
                            update();
                        }}
                        class="flex w-full items-center justify-center">
                        <div
                            class="flex items-center justify-center"
                            on:mouseover={() => (hovered = hoverkey)}
                            on:mouseleave={() => (hovered = "")}>
                            <svelte:component
                                this={guardian.icon}
                                style={selectedGuardians[guardian.name]
                                    ? "background-color: #76708f; color: #fff"
                                    : "color: #56516f"}
                                class="mx-1 my-auto  h-7 w-7 rounded-md p-0.5" />
                            {#if hovered === hoverkey}
                                <div
                                    class="absolute z-50 flex -translate-y-[calc(100%-0.15rem)] flex-row items-center justify-center whitespace-nowrap rounded-lg border border-[#c58597] bg-bouquet-50 p-1.5 text-[#575279]">
                                    <svelte:component
                                        this={guardian.icon}
                                        style={selectedGuardians[guardian.name]
                                            ? "background-color: #76708f; color: #fff"
                                            : "color: #56516f"}
                                        class="mx-1 my-auto h-6 w-6 rounded-md p-0.5" />
                                    <p class="text-sm font-medium">{guardian.name}</p>
                                </div>
                            {/if}
                        </div>
                    </button>
                {/each}
            </div>
            <div
                class="flex flex-row items-center rounded-xl border-[0.75px] border-[#c17e91] bg-[#f2e9e7] p-1.5 opacity-95">
                {#each trials as guardian}
                    {@const hoverkey = "trial_" + guardian.name}
                    <button
                        on:click={() => {
                            selectedTrials[guardian.name] = !selectedTrials[guardian.name];
                            update();
                        }}>
                        <div
                            class="flex items-center justify-center"
                            on:mouseover={() => (hovered = hoverkey)}
                            on:mouseleave={() => (hovered = "")}>
                            <svelte:component
                                this={guardian.icon}
                                style={selectedTrials[guardian.name]
                                    ? "background-color: #76708f; color: #fff"
                                    : "color: #56516f"}
                                class="mx-1 my-auto h-7 w-7 rounded-md p-0.5" />
                            {#if hovered === hoverkey}
                                <div
                                    class="absolute z-50 flex -translate-y-[calc(100%-0.15rem)] flex-row items-center justify-center whitespace-nowrap rounded-lg border border-[#c58597] bg-bouquet-50 p-1.5 text-[#575279]">
                                    <svelte:component
                                        this={guardian.icon}
                                        style={selectedTrials[guardian.name]
                                            ? "background-color: #76708f; color: #fff"
                                            : "color: #56516f"}
                                        class="mx-1 my-auto h-6 w-6 rounded-md p-0.5" />
                                    <p class="text-sm font-medium">Trial {guardian.name}</p>
                                </div>
                            {/if}
                        </div>
                    </button>
                {/each}
            </div>
        </div>
        <div
            class="mx-auto mt-2 flex flex-col items-center justify-center rounded-xl border-[0.75px] border-[#c17e91] bg-[#f2e9e7] p-3 opacity-95">
            <div class="grid grid-cols-8 gap-3">
                {#each list as c}
                    {@const hoverkey = "class_" + c}
                    <button
                        style="background-color: {selectedClasses[c] ? '#76708f' : '#c58799'}"
                        class="flex items-center justify-center rounded-lg p-1"
                        on:mouseover={() => (hovered = hoverkey)}
                        on:mouseleave={() => (hovered = "")}
                        on:click={() => {
                            selectedClasses[c] = !selectedClasses[c];
                            update();
                        }}>
                        <img
                            alt={c}
                            style="-webkit-backface-visibility: hidden;
                            -moz-backface-visibility: hidden;
                            -webkit-transform: translate3d(0, 0, 0);
                            -moz-transform: translate3d(0, 0, 0);"
                            class="z-50 m-auto h-5 w-5 blur-[0.1px]"
                            src={getClassIconNew(c)} />
                        {#if hovered === hoverkey}
                            <div
                                class="absolute z-50 flex -translate-y-[calc(100%)] flex-row items-center justify-center whitespace-nowrap rounded-lg border border-[#c58597] bg-bouquet-50 p-1.5 text-[#575279]">
                                <p class="text-sm font-medium">{c}</p>
                            </div>
                        {/if}
                    </button>
                {/each}
            </div>
        </div>
        </div>
    {/if}
</div>
