<script lang="ts">
    import {onMount} from 'svelte';
    import {tryCatch} from "../lib/try-catch";

    import {config} from "../../wailsjs/go/models";
    import {GetConfig} from '../../wailsjs/go/main/App';

    let isLoading: boolean = false;
    let error: string | null = null;
    let currentConfig: config.Config = {
        version: 1,
        urls: [""],
        accessToken: "",
        pinnedDevices: [],
        scanInterval: 0,
    }

    onMount(async () => {
        isLoading = true;

        const {data, err} = await tryCatch(GetConfig())
        if (err != null) {
            error = err.message;
            console.error("failed to load config:", err);
            isLoading = false;
            return;
        }

        console.log("loaded config:", data);
        currentConfig = data

        isLoading = false;
    });

    async function saveSettings() {
        console.log("saving config:", currentConfig);
    }

    async function testConnectionHandler(url) {
        console.log("testing connection:", url);
    }
</script>

<div class="container mx-auto p-6 max-w-2xl ">
	<h1 class="text-2xl font-bold mb-6">Settings</h1>

	{#if isLoading}
		<div class="flex justify-center my-8">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if error}
		<div class="bg-red-100 dark:bg-red-800 rounded-lg p-6 mb-6">
			<p class="text-red-800 dark:text-red-100">{error}</p>
		</div>
	{:else}
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
			<div class="mb-6">
				<h2 class="text-xl font-semibold mb-4">Home Assistant Connection</h2>

				<div class="mb-4">
					<label class="block text-sm font-medium mb-2">Home Assistant URLs</label>
					<p class="text-sm text-gray-500 mb-2">
						Add multiple URLs for different networks (e.g., local IP for home network, domain for external
						access)
					</p>

					{#each currentConfig.urls as url, index}
						<div class="flex mb-2 gap-2">
							<input
									type="text"
									bind:value={currentConfig.urls[index]}
									placeholder="http://192.168.1.123:8123"
									class="flex-grow text-black p-2 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							/>
							<button
									on:click={() => testConnectionHandler(url)}
									class="px-3 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 rounded-md transition-colors"
									disabled={!url || !currentConfig.accessToken}
									title="Test connection"
							>
								<!--{#if testStatus.testing && testingUrl === url}-->
								<!--  <span class="inline-block animate-spin">⟳</span>-->
								<!--{:else}-->
								✓
								<!--{/if}-->
							</button>
							<button
									class="px-3 py-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
									disabled={currentConfig.urls.length <= 1}
									title="Remove URL"
							>
								✕
							</button>
						</div>
					{/each}

					<button
							class="mt-2 px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-md transition-colors flex items-center"
					>
						<span class="mr-1">+</span> Add URL
					</button>

					<p class="text-sm text-gray-500 mt-2">Example: http://192.168.1.123:8123,
						https://ha.mydomain.local:443</p>
				</div>

				<div class="mb-4">
					<label for="token" class="block text-sm font-medium mb-1">Access Token</label>
					<input
							type="password"
							id="token"
							bind:value={currentConfig.accessToken}
							placeholder="Your long-lived access token"
							class="w-full text-black p-2 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
					/>
					<p class="text-sm text-gray-500 mt-1">
						<a href="https://developers.home-assistant.io/docs/auth_api/#long-lived-access-token"
						   target="_blank"
						   class="text-blue-500 hover:underline">
							How to get a long-lived access token
						</a>
					</p>
				</div>

				<div class="mb-6">
					<label for="scanInterval" class="block text-sm font-medium mb-1">Scan Interval (seconds)</label>
					<input
							type="number"
							id="scanInterval"
							bind:value={currentConfig.scanInterval}
							class="w-full text-black p-2 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
					/>
					<p class="text-sm text-gray-500 mt-1">How often to refresh device states (in seconds)</p>
				</div>

				<!--{#if testStatus.message}-->
				<!--  <div class={`mb-4 p-3 rounded-md ${testStatus.success ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100' : 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100'}`}>-->
				<!--    {testStatus.message}-->
				<!--  </div>-->
				<!--{/if}-->

				<div class="flex justify-end">
					<button
							on:click={saveSettings}
							class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors"
							disabled={isLoading || !currentConfig.accessToken || currentConfig.urls.every(url => !url)}
					>
						Save Settings
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>