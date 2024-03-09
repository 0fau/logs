<script lang="ts">
    export let encounter;
    let user = encounter.uploader;
</script>

<div class="h-8 w-full mb-2 flex items-center max-w-md">
    <div class="flex py-1 px-2 items-center shadow-sm rounded-lg bg-tapestry-50 border border-tapestry-500">
        {#if encounter.anonymized}
            <img alt="avatar" class="rounded-md mr-1 w-6 h-6"
                 src="/icons/misc/sus.png"/>
            <span class="text-xs font-medium text-zinc-900 mr-1.5">Hidden</span>
        {:else}
            {#if user.avatar}
                <img alt="avatar" class="rounded-md mr-1.5 w-6 h-6"
                     src="/images/avatar/{user.id}"/>
            {:else}
                <img alt="avatar" class="rounded-md mr-1 w-6 h-6"
                     src="/icons/misc/play.png"/>
            {/if}
            <span class="text-sm text-zinc-900 mr-2">{user.username ? user.username : user.discordTag}</span>
        {/if}
    </div>
</div>