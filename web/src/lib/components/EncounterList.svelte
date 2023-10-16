<script lang="ts">
    import {formatDamage, formatDate, formatDuration} from "$lib/components/meter/print";
    import IconArrow from '~icons/carbon/next-filled'

    export let encounters;
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
        "Plague Legion Commander Akkan": {
            icon: "/icons/raids/akkan.png"
        },
    };

    function getRaidIcon(boss) {
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

        encounters = encounters.concat(recent)
        if (recent.length < 5) {
            more = false
        }
        if (recent.length > 0) {
            page++
        }
    }
</script>

<div class="flex flex-row m-auto mt-10 justify-center items-center flex flex-col">
    <div class="w-[46%] h-[36px] bg-[#6d6797] text-center text-[#F4EDE9] text-sm flex flex-row justify-evenly items-center rounded-md shadow-sm mb-5">
        <div class="w-full h-fit border-r-[#dbd5d1] border-r-[0.5px]">
            <button class="w-full">Arkesia</button>
        </div>
        <div class="w-full h-fit border-r-[#dbd5d1] border-r-[0.5px]">
            <button class="w-full">Friends</button>
        </div>
        <div class="w-full h-fit">
            <button class="w-full">Me</button>
        </div>
    </div>
    {#each encounters.slice(page * 5, (page + 1) * 5) as encounter}
        <a href="/log/{encounter.id}"
           class="mb-6 h-[80px] flex shadow-sm rounded-md w-4/5 bg-[#F4EDE9]">
            <div class="w-full h-full flex flex-row ml-5 items-center">
                <div>
                    <div class="self-start text-[#575279]">
                        <div>
                            <span class="font-medium">[#{encounter.id}]</span>
                            <img src={getRaidIcon(encounter.boss)}
                                 class="inline w-6 h-6 -translate-y-0.5"/>
                            <span class="font-medium">{encounter.boss}</span>
                        </div>
                        <p class="text-sm">{formatDamage(encounter.damage)} damage dealt
                            in {formatDuration(encounter.duration)}</p>
                        <p class="text-xs text-[#5d5978]">{formatDate(encounter.date)}</p>
                    </div>
                </div>
                <div class="py-1 px-1.5 h-full ml-auto self-end flex flex-col rounded-r-md text-white">
                    <span class="text-xs text-center self-end text-[#F4EDE9] p-0.5 px-1 mr-0.5 mt-1.5 rounded-sm bg-[#6d6797] font-medium">{encounter.localPlayer}</span>
                    <span class="text-xs text-[#6d6797] self-end text-right mr-0.5 mt-0.5 font-medium">{encounter.players[encounter.localPlayer].class}</span>
                    <span class="text-[#575279] text-right mr-1 my-auto text-lg font-medium">{formatDamage(encounter.players[encounter.localPlayer].dps)}</span>
                </div>
            </div>
        </a>
    {/each}
    <div class="flex flex-row text-[#6d6797] items-center justify-center">
        {#if page > 0}
            <button on:click={prev} class="bg-[#F4EDE9] p-0.5 rounded-3xl shadow-sm">
                <IconArrow class="rotate-180 w-6 h-6"/>
            </button>
        {/if}
        {#if more || (page + 1) * 5 < encounters.length}
            <button on:click={next} class="ml-5 bg-[#F4EDE9] p-0.5 rounded-3xl shadow-sm">
                <IconArrow class="w-6 h-6"/>
            </button>
        {/if}
    </div>
</div>