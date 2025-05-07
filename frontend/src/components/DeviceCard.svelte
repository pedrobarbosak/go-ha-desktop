<script lang="ts">
	import {createEventDispatcher, onMount} from "svelte";
	import {ha} from "../../wailsjs/go/models";
	import {tryCatch} from "../lib/try-catch";
	import {TurnOn, TurnOff} from "../../wailsjs/go/main/App";

	const dispatch = createEventDispatcher();

	let isLoading = false;
	let error: string | null = null;
	export let device: ha.Device;
	export let isFavorite: boolean = false;

	async function toggleDevice() {
		isLoading = true;

        const {data, err} = await tryCatch(device.State ? TurnOff(device.ID) : TurnOn(device.ID))
		if (err != null) {
			error = err.message;
			console.error("Failed to toggle device:", err);
			isLoading = false;
			return;
		}

        device = data;
        isLoading = false;
	}

    function toggleFavorite(event: MouseEvent) {
        event.stopPropagation();
        dispatch('toggleFavorite');
    }

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			toggleDevice();
		}
	}

	onMount(async () => {
		isLoading = false;
	})
</script>

<div 
	class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 mb-4 cursor-pointer hover:shadow-lg transition-shadow"
	on:click={toggleDevice}
	on:keydown={handleKeyDown}
	tabindex="0"
	role="button"
	aria-label={`Toggle ${device.Name || device.ID}`}
>
	{#if isLoading}
		<div class="flex justify-center my-2">
			<div class="animate-spin rounded-full h-6 w-6 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if device.Error}
		<div class="flex items-center justify-between">
			<div>
				<h3 class="text-lg font-semibold">{device.Name || device.ID}</h3>
				<p class="text-sm text-red-500">{device.Error}</p>
			</div>
			<div class="flex items-center gap-2">
				<button
					on:click={toggleFavorite}
					class={`text-xl ${isFavorite ? 'text-yellow-500' : 'text-gray-400 hover:text-yellow-500'}`}
					title={isFavorite ? "Remove from favorites" : "Add to favorites"}
				>
					★
				</button>
				<div class="bg-gray-300 dark:bg-gray-600 rounded-full w-10 h-6 flex items-center p-1">
					<div class="bg-gray-400 dark:bg-gray-500 rounded-full w-4 h-4"></div>
				</div>
			</div>
		</div>
	{:else}
		<div class="flex items-center justify-between">
			<div>
				<h3 class="text-lg font-semibold">{device.Name || device.ID}</h3>
				<p class="text-sm text-gray-500">{device.Type}</p>
			</div>
			<div class="flex items-center gap-2">
				<button
					on:click={toggleFavorite}
					class={`text-xl ${isFavorite ? 'text-yellow-500' : 'text-gray-400 hover:text-yellow-500'}`}
					title={isFavorite ? "Remove from favorites" : "Add to favorites"}
				>
					★
				</button>
				<div
					class={`rounded-full w-10 h-6 flex items-center p-1 ${device.State ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}
				>
					<div
						class={`bg-white rounded-full w-4 h-4 transform transition-transform ${device.State ? 'translate-x-4' : ''}`}
					></div>
				</div>
			</div>
		</div>
	{/if}
</div>
