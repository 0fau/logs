<script lang="ts">
    import {settingsUI} from "$lib/menu";
    import {user} from "$lib/store";

    export let settings;

    let username = "";
    let token = "";

    console.log(settings)

    async function saveUsername() {
        let resp = await fetch("/api/settings/username?username=" + username, {
            method: 'PUT',
            credentials: 'same-origin'
        })
        if (resp.ok) {
            user.update(user => {
                user.username = username
                return user
            })
        }
    }

    async function generateAccessToken() {
        let resp = await fetch("/api/users/@me/token", {
            method: 'POST',
            credentials: 'same-origin'
        })
        if (resp.ok) {
            token = (await resp.json()).token
        }
    }

    async function revokeAccessToken() {
        let resp = await fetch("/api/users/@me/token?revoke=true", {
            method: 'POST',
            credentials: 'same-origin'
        })
        if (resp.ok) {
            token = null
            settings.hasToken = false
        }
    }

    function copy() {
        navigator.clipboard.writeText(token)
    }
</script>

<div class="w-full h-full">
    <span class="font-semibold bg-[#575279] text-[#faf4ed] p-0.5 px-1">Settings</span>
    <button class="inline font-semibold float-right" on:click={() => settingsUI.set(false)}>[x]</button>
    <p class="mt-1">um janky ui alert i'll revisit later :D</p>
    <p class="mt-1 font-semibold">Username</p>
    <div class="flex items-center">
        <input bind:value={username} class="m-1 rounded-md py-0.5 px-2 w-[140px]"
               placeholder="{$user && $user.username ? $user.username : '(unset)'}"/>
        <button on:click={saveUsername} class="bg-[#575279] text-[#faf4ed] px-1 py-0.5 ml-1 rounded-sm h-full text-sm">
            save
        </button>
    </div>
    <p class="mt-1 font-semibold">Access Token</p>
    <div class="m-1 mb-1.5 w-[200px] h-[26px] px-2 flex items-center rounded-md bg-white text-sm">
        {#if token}
            <p class="w-full overflow-hidden whitespace-nowrap text-ellipsis">{token}</p>
            <button on:click={copy} class="font-semibold">copy</button>
        {:else if settings.hasToken}
            (hidden)
        {:else}
            (unset)
        {/if}
    </div>
    <button on:click={generateAccessToken} class="bg-[#575279] text-[#faf4ed] px-1 py-0.5 ml-1 rounded-sm text-sm">
        generate
    </button>
    <button on:click={revokeAccessToken} class="bg-[#575279] text-[#faf4ed] px-1 py-0.5 ml-1 rounded-sm text-sm">
        revoke
    </button>
</div>