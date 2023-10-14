<script lang="ts">
    import {
        formatDamage,
        formatDuration,
        formatDate
    } from "$lib/components/meter/print";

    export let encounters;
</script>

<div class="flex flex-row m-auto mt-10 justify-center items-center flex flex-col">
    {#each encounters as encounter}
        <a href="/log/{encounter.id}"
           class="mb-6 h-[80px] flex transition shadow-sm rounded-md w-4/5 bg-[#F4EDE9]">
            <div class="w-full h-full flex flex-col ml-5 justify-center items-center">
                <div class="self-start text-[#575279]">
                    <p class="font-medium">[#{encounter.id}] {encounter.boss}</p>
                    <p class="text-sm">{formatDamage(encounter.damage)} damage dealt in {formatDuration(encounter.duration)}</p>
                    <p class="text-xs text-[#5d5978]">{formatDate(encounter.date)}</p>
                </div>
                <div class="self-end font-medium text-[#605d78] w-fit h-fit grid grid-cols-2">
                </div>
            </div>
        </a>
    {/each}
</div>