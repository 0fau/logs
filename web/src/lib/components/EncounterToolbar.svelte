<script lang="ts">
    import IconSettings from '~icons/ep/setting';

    export let encounter;
    let uploader = encounter.uploader;
    export let user;

    let settingsDropdown = false;

    let settings = {};

    let job;

    function changeSetting(setting: string, value: boolean) {
        if (job) {
            clearTimeout(job)
        }

        settings[setting] = value;

        job = setTimeout(() => {
            changeSettingAPI(setting, value)
        }, 500)
    }

    async function changeSettingAPI(setting: string, value: boolean) {
        const resp = await fetch("/api/log/" + encounter.id + "/settings", {
            credentials: 'same-origin',
            method: 'PATCH',
            body: JSON.stringify(settings)
        });

        settings = {};
    }
</script>

<div class="h-4 mb-2 flex flex-row items-center max-w-xl w-full">
    <a href="/logs" class="mt-1 text-sm text-tapestry-800 hover:underline">← Logs</a>
</div>
<div class="h-8 w-full mb-2 flex items-center max-w-md px-1">
    <div class="flex p-0.5 items-center shadow-sm rounded-lg bg-tapestry-50 border border-tapestry-500">
        {#if !uploader}
            <img alt="avatar" class="rounded-md mr-1 w-6 h-6"
                 src="/icons/misc/sus.png"/>
            <span class="text-xs font-medium text-zinc-900 mr-1.5">Hidden</span>
        {:else}
            {#if uploader.avatar}
                <img alt="avatar" class="rounded-md mr-1.5 w-6 h-6"
                     src="/images/avatar/{uploader.id}"/>
            {:else}
                <img alt="avatar" class="rounded-md mr-1 w-6 h-6"
                     src="/icons/misc/play.png"/>
            {/if}
            <span class="text-sm text-[#413d5b] mr-2">{uploader.username ? uploader.username : uploader.discordTag}</span>
        {/if}
    </div>
    <div class="ml-auto mr-0 flex items-center">
        {#if uploader && user && uploader.id === user.id}
            <button class="bg-tapestry-500 text-tapestry-50 rounded-md p-0.5"
                    on:click={() => settingsDropdown =!settingsDropdown}>
                <IconSettings/>
            </button>
            {#if settingsDropdown}
                <div class="absolute flex flex-col items-start justify-center p-2 z-50 rounded-lg whitespace-nowrap bg-bouquet-50 border border-tapestry-600 translate-y-[calc(100%-1.9em)] -translate-x-[calc(100%-1.7rem)] text-gray-700">
                    <div class="flex items-center justify-center">
                        <div class="h-fit ml-1">
                            <p class="text-sm font-semibold">Log Visibility</p>
                            <div class="flex items-center">
                                <input on:click={() => changeSetting("names", 1)}
                                       class="w-3.5 h-3.5 mr-1"
                                       type="radio" name="logvisibility" value="1"
                                       checked={settings.visibility?.names === 1}/>
                                <p class="text-sm">Show all names</p>
                            </div>
                            <div class="flex items-center">
                                <input on:click={() => changeSetting("names", 2)}
                                       class="w-3.5 h-3.5 mr-1"
                                       type="radio" name="logvisibility" value="2"
                                       checked={settings.visibility?.names === 2}/>
                                <p class="text-sm">Show own name</p>
                            </div>
                            <div class="flex items-center">
                                <input on:click={() => changeSetting("names", 3)}
                                       class="w-3.5 h-3.5 mr-1"
                                       type="radio" name="logvisibility" value="3"
                                       checked={settings.visibility?.names === 3}/>
                                <p class="text-sm">Hide names</p>
                            </div>
                        </div>
                    </div>
                    <!--                    <div class="flex items-center justify-center">-->
                    <!--                        <input type="checkbox" bind:checked={canSearch}-->
                    <!--                               class="mr-1.5 focus:ring-0 focus:ring-offset-0 rounded w-4 h-4"/>-->
                    <!--                        <span class="text-sm font-medium">can search</span>-->
                    <!--                    </div>-->
                </div>
            {/if}
        {/if}
    </div>
</div>
