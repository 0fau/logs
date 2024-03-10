<script lang="ts">
    import { blur } from "svelte/transition";
    import { settingsUI } from "$lib/menu";
    import { settings, user } from "$lib/store";

    import IconAnalytics from "~icons/material-symbols/analytics-outline";
    import IconAnalytics2 from "~icons/carbon/analytics";
    import IconCommunity from "~icons/iconoir/community";
    import IconLog from "~icons/ph/sword-light";
    import IconRanking from "~icons/solar/ranking-broken";
    import IconProfile from "~icons/mingcute/profile-line";
    import IconRecord from "~icons/ph/vinyl-record";
    import IconInbox from "~icons/ic/baseline-inbox";
    import IconAbout from "~icons/mdi/about-circle-outline";
    import IconDiscord from "~icons/prime/discord";
    import IconStar from "~icons/ph/shooting-star";
    import IconOpen from "~icons/nimbus/ellipsis";
    import IconLogout from "~icons/material-symbols/logout-sharp";

    import IconSetting from "~icons/uil/setting";

    import { page } from "$app/stores";
    import { browser } from "$app/environment";

    $: selected = $page.url.pathname.substring(1).split("/")[0];

    function style(select, item) {
        if (select == item) {
            return "p-1 pl-2 flex items-center rounded w-[65%] mb-1 bg-[#c58597]";
        } else {
            return "w-[65%] mb-2 opacity-95";
        }
    }

    let toggle;
    let panelOpen = false;

    function unfocus(box) {
        const click = (event) => {
            if (!box.contains(event.target) && !toggle.contains(event.target)) {
                panelOpen = false;
            }
        };

        document.addEventListener("click", click, true);
        return {
            destroy() {
                document.removeEventListener("click", click, true);
            }
        };
    }

    function openSettings() {
        settingsUI.set(true);
    }

    $: greeting = $settings?.logs.announcement;
</script>

<div class="hidden h-full w-72 min-w-72 flex-col bg-tapestry-600 md:flex">
    <div class="flex h-[20%] w-full flex-col items-center justify-center">
        <p class="text-2xl font-medium text-[#fcf9f6]">
            <IconAnalytics class="thin inline h-8 w-8 -translate-y-0.5" />
            Logs
        </p>
        <p class="translate-y-1 text-sm text-[#fcf9f6]">by faust</p>
    </div>
    <div class="mb-5 flex w-full flex-col pl-[16%] text-lg text-[#fff]">
        <p class="mb-1 text-sm font-semibold">General</p>
        <a href="/logs" class={style(selected, "logs")}>
            <IconLog class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Logs</span>
        </a>
        <a href="/rankings" class={style(selected, "rankings")}>
            <IconRanking class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Rankings</span>
        </a>
        <a href="/featured" class={style(selected, "featured")}>
            <IconStar class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Featured</span>
        </a>
        <a href="/trends" class={style(selected, "trends")}>
            <IconAnalytics2 class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Trends</span>
        </a>
        <p class="my-1 text-sm font-semibold">Social</p>
        <a href="/profile" class={style(selected, "profile")}>
            <IconProfile class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Profile</span>
        </a>
        <a href="/inbox" class={style(selected, "inbox")}>
            <IconInbox class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Inbox</span>
        </a>
        <a href="/community" class={style(selected, "community")}>
            <IconCommunity class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Community</span>
        </a>
        <p class="my-1 text-sm font-semibold">Misc</p>
        <a href="/musicbox" class={style(selected, "musicbox")}>
            <IconRecord class="inline h-6 w-6" />
            <span class="ml-1 text-sm">Music Box</span>
        </a>
        <a href="https://docs.fau.dev/logs/#faq" class={style(selected, "faq")}>
            <IconAbout class="inline h-6 w-6" />
            <span class="ml-1 text-sm">FAQ</span>
        </a>
         <button
            class="h-8 w-20 rounded-md bg-tapestry-400 p-1 text-sm"
            on:click={() => {
                $settings.darkMode = !$settings.darkMode;
            }}>
            dark mode
        </button>
        {#if browser && !greeting}
            <button
                transition:blur={{ duration: 10 }}
                on:click={() => ($settings.logs.announcement = !$settings.logs.announcement)}
                class="mt-4 shadow-sm rounded-xl py-1.5 px-4 w-36 flex justify-center items-center dark:bg-tapestry-400 hover:dark:bg-tapestry-500 border border-white text-gray-200"
            >
                <span class="font-medium mr-1 text-sm">hewwo ^.^ </span>
                <img
                    alt="avatar"
                    class="inline h-7"
                    src="https://cdn.discordapp.com/emojis/667830310741737473.gif?size=240&quality=lossless"
                />
            </button>
        {/if}
    </div>

    {#if $user && $user.id}
        <div class="mb-10 mt-auto flex flex-row items-center justify-center">
            <div class="m-1 rounded bg-[#9a4a61] py-1 opacity-0 shadow-sm">
                <IconOpen class="text-[#fcf9f6]" />
            </div>
            <div class="flex items-center justify-center rounded bg-[#9a4a61] px-4 py-2 shadow-sm">
                <span class="text-xs text-[#fcf9f6]"
                    >{$user.username ? $user.username : "@" + $user.discordTag}</span>
            </div>
            <button
                bind:this={toggle}
                on:click={() => (panelOpen = !panelOpen)}
                class="m-1 rounded bg-[#9a4a61] py-1 shadow-sm">
                <IconOpen class="text-[#fcf9f6]" />
            </button>
            {#if panelOpen}
                <div
                    transition:blur={{ duration: 10 }}
                    class="absolute flex h-fit w-[120px] -translate-y-[52px] flex-col items-center justify-center rounded bg-[#9a4a61] py-0.5 text-center text-xs text-[#fcf9f6] shadow-sm"
                    use:unfocus>
                    <button
                        on:click={openSettings}
                        class="flex h-[30px] w-[90%] items-center justify-center">
                        <IconSetting class="mr-1 inline" />
                        Settings
                    </button>
                    <div class="flex h-[30px] w-[90%] items-center justify-center">
                        <form action="/logout" method="post">
                            <button class="flex items-center justify-center">
                                <IconLogout class="mr-1 inline" />
                                Logout
                            </button>
                        </form>
                    </div>
                </div>
            {/if}
        </div>
    {:else}
        <div class="mx-auto mt-auto">
            <form action="/oauth2" method="post">
                <button
                    class="mb-10 w-[50%] min-w-[110px] rounded bg-[#9a4a61] p-1 py-1.5 text-center text-[#fcf9f6]">
                    <IconDiscord class="inline h-6 w-6" />
                    <span class="text-sm">Sign in</span>
                </button>
            </form>
        </div>
    {/if}
</div>
