<script lang="ts">
    export let encounter;
    let user = encounter.uploader;
</script>

<div class="h-[23px] w-[458px] mb-2 flex flex-row items-center">
    <div class="flex p-[1.5px] items-center shadow-sm rounded-lg bg-[#f4e8e6] border-[1px] border-[#c58597]">
        {#if encounter.anonymized}
            <img alt="avatar" class="rounded-md mr-1 w-6 h-6"
                 src="/icons/misc/sus.png"/>
            <span class="text-xs font-medium text-[#413d5b] mr-1.5">Hidden</span>
        {:else}
            <img alt="avatar" class="rounded-md mr-1.5 w-6 h-6"
                 src="https://cdn.discordapp.com/avatars/{user.discordId}/{user.avatar}.png"/>
            <span class="text-sm text-[#413d5b] mr-2">{user.username ? user.username : user.discordTag}</span>
        {/if}
    </div>
</div>