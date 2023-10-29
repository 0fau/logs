<script lang="ts">
    import {formatDamage, formatDate, formatDuration} from "$lib/components/meter/print";
    import IconArrow from '~icons/carbon/next-filled'
    import IconBack from '~icons/ion/arrow-back-outline'
    import IconScope from '~icons/mdi/telescope'

    import EncounterRecap from "$lib/components/EncounterRecap.svelte";

    export let encounters;
    scan(encounters)
    process(encounters)
    console.log(encounters)

    const raids = {
        "Dark Mountain Predator": {
            icon: "/icons/raids/valtan.png"
        },
        "Ravaged Tyrant of Beasts": {
            icon: "/icons/raids/valtan.png"
        },
        "Saydon": {
            icon: "/icons/raids/clown.png"
        },
        "Kakul": {
            icon: "/icons/raids/clown.png"
        },
        "Encore-Desiring Kakul-Saydon": {
            icon: "/icons/raids/clown.png"
        },
        "Incubus Morphe": {
            icon: "/icons/raids/vykas.png"
        },
        "Nightmarish Morphe": {
            icon: "/icons/raids/vykas.png"
        },
        "Covetous Devourer Vykas": {
            icon: "/icons/raids/vykas.png"
        },
        "Covetous Legion Commander Vykas": {
            icon: "/icons/raids/vykas.png"
        },
        "Ashtarot": {
            icon: "/icons/raids/brelshaza.png"
        },
        "Primordial Nightmare": {
            icon: "/icons/raids/brelshaza.png"
        },
        "Gehenna Helkasirs": {
            icon: "/icons/raids/brelshaza.png"
        },
        "Evolved Maurug": {
            icon: "/icons/raids/akkan.png"
        },
        "Lord of Degradation Akkan": {
            icon: "/icons/raids/akkan.png"
        },
        "Plague Legion Commander Akkan": {
            icon: "/icons/raids/akkan.png"
        },
    };
    let focused = encounters[0];

    function focus(encounter) {
        focused = encounter
    }

    function getRaidIcon(boss: string) {
        return raids[boss] ? raids[boss].icon : "/icons/raids/surprised.png"
    }

    let page = 0;
    let more = encounters.length == 5;

    function prev() {
        if (page > 0) {
            page--
        }
    }

    async function next() {
        if (encounters.length === 0) {
            return
        }

        if ((page + 1) * 5 < encounters.length) {
            page++
            return
        }

        if (!more) {
            return
        }

        let last = encounters[encounters.length - 1];
        let url = location.protocol + '//' + location.host;
        const recent = await (await fetch(
            url + "/api/logs/@recent?past=" + last.date + "&id=" + last.id
        )).json()
        scan(recent)
        process(recent)

        encounters = encounters.concat(recent)
        if (recent.length < 5) {
            more = false
        }
        if (recent.length > 0) {
            page++
        }
    }

    // scan looks for encounters with missing data
    function scan(encounters) {
        for (let i = 0; i < encounters.length; i++) {
            let broken = false;
            let count = 0;

            const encounter = encounters[i]
            out: for (let party of encounter.parties) {
                count += party.length

                for (let player of party) {
                    if (!encounter.players[player]) {
                        broken = true
                        break out
                    }
                }
            }

            if (count != Object.keys(encounter.players).length) {
                broken = true
            }
            encounter.broken = broken
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

    $: display = focused ? [focused] : encounters.slice(page * 5, (page + 1) * 5)
</script>

<div class="m-auto mt-10 flex flex-col justify-center items-center">
    <div class="flex flex-row w-[88%] justify-center items-center">
        <div class="w-[270px] h-[36px] px-1 bg-[#b96d83] text-center text-[#F4EDE9] text-sm flex flex-row justify-center items-center rounded-sm mb-3">
            <div class="w-[80%] h-[76%] rounded-sm flex justify-center items-center bg-[#F4EDE9] shadow-sm">
                <button class="w-full font-medium text-[#b4637a]">Arkesia</button>
            </div>
            <div class="w-full h-full flex justify-center items-center">
                <button class="w-full">Friends</button>
            </div>
            <div class="w-full h-full flex justify-center items-center">
                <button class="w-full">Roster</button>
            </div>
        </div>
    </div>
    <div class="w-[88%] flex flex-col justify-center items-center bg-[#dec5cd] pt-4 mb-3 rounded-md">
        {#each display as encounter}
            <button
                    class="{focused ? '' : 'mb-5'} h-[80px] flex border-[0.5px] border-[#c58597] shadow-sm rounded-md w-[94%] bg-[#F4EDE9]"
                    on:click={() => focus(focused ? null : encounter)}>
                <div class="w-full h-full flex flex-row ml-5 items-center">
                    <div>
                        <div class="self-start text-left text-[#575279]">
                            <div>
                                <span class="font-medium">[#{encounter.id}]</span>
                                <img alt={encounter.boss} src={getRaidIcon(encounter.boss)}
                                     class="inline w-6 h-6 -translate-y-0.5"/>
                                <span class="font-medium">{encounter.boss}</span>
                            </div>
                            <p class="text-sm">{formatDamage(encounter.damage)} damage dealt
                                in {formatDuration(encounter.duration)}</p>
                            <p class="text-xs text-[#5d5978]">{formatDate(encounter.date)}</p>
                        </div>
                    </div>
                    <div class="py-1 px-1.5 h-full ml-auto self-end flex flex-col rounded-r-md text-white">
                        <span class="text-xs text-center self-end text-[#F4EDE9] p-0.5 px-1 mr-0.5 mt-1.5 rounded-sm bg-[#b4637a] font-medium">{encounter.localPlayer}</span>
                        <span class="text-xs text-[#b4637a] self-end text-right mr-0.5 mt-0.5 font-medium">{encounter.players[encounter.localPlayer].class}</span>
                        <span class="text-[#575279] text-right mr-1 my-auto text-lg font-medium">{formatDamage(encounter.players[encounter.localPlayer].dps)}</span>
                    </div>
                </div>
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
                <div class="mx-auto mt-auto mb-1 p-0.5 px-6 rounded-md text-[#f7f2ef] bg-[#b96d83]">Recap</div>
                <div class="my-1">
                    <button class="flex items-center justify-center bg-[#f3eeec] p-0.5 px-1.5 border-[0.5px] border-[#575279] rounded-md text-[#575279]"
                            on:click={() => window.open("/log/" + focused.id, '_blank').focus()}>
                        <IconScope class="inline mr-0.5"/>
                        Inspect
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