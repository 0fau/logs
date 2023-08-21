<script lang="ts">
    import type {PageData} from './$types';

    export let data: PageData;
</script>

<div class="container flex h-screen justify-center items-center">
    <div class="columns-3xs text-center ...">
        {#if data.username}
            <p>{data.username}</p>
            <form action="/logout" method="post">
                <button class="underline">Logout</button>
            </form>
        {:else}
            <form action="/oauth2" method="post">
                <button class="underline">Login</button>
            </form>
        {/if}
    </div>
</div>

<style lang="postcss">
    :global(html) {
        background-color: theme(colors.gray.100);
    }
</style>