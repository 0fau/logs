<script lang="ts">
    export let encounter;
    let user = encounter.user;
</script>

<div class="h-[23px] w-[458px] mb-2 flex items-center">
    <div class="flex p-[1.5px] items-center shadow-sm rounded-lg bg-[#f4e8e6] border-[1px] border-[#c58597]">
        <img alt="avatar" class="rounded-md mr-1.5 w-6 h-6"
             src="https://cdn.discordapp.com/avatars/{user.discord_id}/{user.avatar}.png"/>
        <span class="text-sm text-[#413d5b] mr-2">{user.username ? user.username : user.discord_tag}</span>
    </div>
</div>