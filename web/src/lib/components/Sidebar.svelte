<script lang="ts">
    import IconAnalytics from '~icons/material-symbols/analytics-outline'
    import IconAnalytics2 from '~icons/carbon/analytics'
    import IconCommunity from '~icons/iconoir/community'
    import IconLog from '~icons/ph/sword-light'
    import IconRanking from '~icons/solar/ranking-broken'
    import IconProfile from '~icons/mingcute/profile-line'
    import IconRecord from '~icons/ph/vinyl-record'
    import IconInbox from '~icons/ic/baseline-inbox'
    import IconAbout from '~icons/mdi/about-circle-outline'
    import IconDiscord from '~icons/prime/discord'
    import IconStar from '~icons/ph/shooting-star'

    export let selected;
    export let user;

    function style(item) {
        if (selected == item) {
            return "p-1 pl-2 flex items-center rounded w-[65%] mb-1 bg-[#c58597]"
        } else {
            return "w-[65%] mb-2 opacity-95"
        }
    }
</script>

<div class="w-full h-[20%] flex flex-col justify-center items-center">
    <p class="text-[#fcf9f6] font-medium text-2xl">
        <IconAnalytics class="inline thin -translate-y-0.5 w-8 h-8"/>
        Logs
    </p>
    <p class="text-[#fcf9f6] translate-y-1 text-sm">by faust</p>
</div>
<div class="w-full flex flex-col text-lg text-[#fff] pl-[16%] mb-5">
    <p class="mb-1 text-sm font-semibold">General</p>
    <a href="/logs" class="{style('logs')}">
        <IconLog class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Logs</span>
    </a>
    <a href="/rankings" class="{style('rankings')}">
        <IconRanking class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Rankings</span>
    </a>
    <a href="/featured" class="{style('featured')}">
        <IconStar class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Featured</span>
    </a>
    <a href="/trends" class="{style('trends')}">
        <IconAnalytics2 class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Trends</span>
    </a>
    <p class="my-1 text-sm font-semibold">Social</p>
    <a href="/profile" class="{style('profile')}">
        <IconProfile class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Profile</span>
    </a>
    <a href="/inbox" class="{style('inbox')}">
        <IconInbox class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Inbox</span>
    </a>
    <a href="/community" class="{style('community')}">
        <IconCommunity class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Community</span>
    </a>
    <p class="my-1 text-sm font-semibold">Misc</p>
    <a href="/musicbox" class="{style('musicbox')}">
        <IconRecord class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">Music Box</span>
    </a>
    <a href="/faq" class="{style('faq')}">
        <IconAbout class="inline w-6 h-6"/>
        <span class="ml-1 text-sm">FAQ</span>
    </a>
</div>

{#if user}
    <div class="mx-auto">
        <p class="text-sm font-medium text-white">{user.username}</p>
    </div>
{:else}
    <div class="mx-auto">
        <form action="/oauth2" method="post">
            <button class="self-end min-w-[110px] p-1 py-1.5 rounded text-center text-[#fcf9f6] w-[50%] bg-[#9a4a61]">
                <IconDiscord class="inline w-6 h-6"/>
                <span class=" text-sm">Sign in</span>
            </button>
        </form>
    </div>
{/if}