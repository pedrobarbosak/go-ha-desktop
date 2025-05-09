<script lang="ts">
    import {onMount} from "svelte";
    import {tryCatch} from "../lib/try-catch";

    import {config, ha} from "../../wailsjs/go/models";
    import {GetConfig, UpdateConfig} from "../../wailsjs/go/main/App";
    import DeviceCard from "../components/DeviceCard.svelte";

    import {devices, error, loadDevices, loading} from '../lib/stores/devices';

    let isLoading: boolean = false;
    let errorMessage: string | null = null;
    let favoriteDevices: ha.Device[] = [];
    let missingFavorites: string[] = [];

    let cfg: config.Config = {
        version: 1,
        urls: [""],
        accessToken: "",
        pinnedDevices: [],
        scanInterval: 0,
    }

    $: errorMessage = $error;

    $: {
        favoriteDevices = [];
        missingFavorites = [];

        for (const device of $devices) {
            if (cfg.pinnedDevices.includes(device.ID)) {
                favoriteDevices.push(device);
            }
        }

        missingFavorites = cfg.pinnedDevices.filter(
            id => !$devices.some(device => device.ID === id)
        );
    }

    async function loadConfig() {
        const {data, err} = await tryCatch(GetConfig())
        if (err != null) {
            errorMessage = err.message;
            console.error("failed to load config:", err);
            return;
        }

        cfg = data;
    }

    async function toggleFavorite(device: ha.Device) {
        const index = cfg.pinnedDevices.indexOf(device.ID);
        if (index === -1) {
            cfg.pinnedDevices.push(device.ID);
        } else {
            cfg.pinnedDevices.splice(index, 1);
        }

        // Update config
        const {err} = await tryCatch(UpdateConfig(cfg));
        if (err != null) {
            errorMessage = err.message;
            console.error("failed to update config:", err);
            return;
        }
    }

    async function refreshDevices() {
        isLoading = true;
        await loadDevices();
        isLoading = false;
    }

    onMount(async () => {
        await loadConfig();
        await refreshDevices();
    })
</script>

<div class="container mx-auto p-6 max-w-2xl">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-3xl font-bold">Devices</h1>
		<button
				class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors flex items-center"
				disabled={$loading || isLoading}
				on:click={refreshDevices}
		>
			<span class={$loading ? "animate-spin mr-2" : "mr-2"}>⟳</span>
			Refresh
		</button>
	</div>

	{#if $loading}
		<div class="flex justify-center my-8">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if $error}
		<div class="bg-red-100 dark:bg-red-800 rounded-lg p-6 mb-6">
			<p class="text-red-800 dark:text-red-100">{$error}</p>
		</div>
	{:else}
		{#if favoriteDevices.length > 0 || missingFavorites.length > 0}
			<div class="mb-8">
				<h2 class="text-xl font-semibold mb-4">Favorites</h2>

				{#each favoriteDevices as device}
					<DeviceCard
							bind:device={device}
							isFavorite={true}
							on:toggleFavorite={() => toggleFavorite(device)}
					/>
				{/each}

				{#each missingFavorites as missingId}
					<div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 mb-4">
						<div class="flex items-center justify-between">
							<div>
								<h3 class="text-lg font-semibold">{missingId}</h3>
								<p class="text-sm text-red-500">Device not found</p>
							</div>
							<button
									on:click={() => toggleFavorite({ID: missingId, Name: "", Type: "", State: false, Error: null})}
									class="text-red-500 hover:text-red-700"
									title="Remove from favorites"
							>
								★
							</button>
						</div>
					</div>
				{/each}
			</div>

			<h2 class="text-xl font-semibold mb-4">All Devices</h2>
		{/if}

		{#each $devices.filter(d => !cfg.pinnedDevices.includes(d.ID)) as device}
			<DeviceCard
					bind:device={device}
					isFavorite={false}
					on:toggleFavorite={() => toggleFavorite(device)}
			/>
		{/each}
	{/if}
</div>