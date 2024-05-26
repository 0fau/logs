<script>
    export let profile;
</script>

<div class="w-fit px-2 py-2 mb-2 flex flex-col justify-center items-center border-[0.5px] border-[#c58597] shadow-sm rounded-md bg-[#f8f3ed]">
    <div class="w-full h-full flex flex-col items-center justify-evenly">
        <div class="grid grid-rows-2 grid-cols-4 gap-2">
            {#each profile as character}
                <div class="h-[88px] w-[187.5px] bg-[#F4EDE9] rounded-sm flex flex-col justify-center items-center">
                    <div class="flex flex-col items-center justify-evenly text-[#a2596d]">
                        <p class="text-sm font-medium">{character.character}</p>
                        <p class="text-xs">{Math.floor(character.gearScore)} {character.class}</p>
                        <p class="font-medium">{character.logCount}</p>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>