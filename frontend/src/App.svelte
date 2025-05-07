<script lang="ts">
    import {onMount} from "svelte";

    import {tryCatch} from "./lib/try-catch";
    import {GetError} from "../wailsjs/go/main/App";

    import Settings from "./pages/Settings.svelte";
    import Devices from "./pages/Devices.svelte";

    let error: string | null = null;

    let path = window.location.hash;
    onMount(async () => {
        const {err} = await tryCatch(GetError());
        if (err != null) {
            error = err.message;
            console.error("failed to load config:", err);
        }

        const updatePath = () => {
            path = window.location.hash || "#";
        };

        updatePath();

        window.addEventListener('hashchange', updatePath);
        return () => {
            window.removeEventListener('hashchange', updatePath);
        };
    });

</script>

{#if error}
    <div class="bg-red-100 dark:bg-red-800 rounded-lg p-6 mb-6">
        <p class="text-red-800 dark:text-red-100">{error}</p>
    </div>
{/if}

<div class="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <nav class="bg-white dark:bg-gray-800 shadow-md p-4">
        <div class="container mx-auto flex justify-between items-center">
            <a href="/#" class="text-xl font-bold">HA-Desktop</a>
            <div class="space-x-4">
                <a href="/#"  class="hover:text-blue-500 transition-colors">Devices</a>
                <a href="/#settings" class="hover:text-blue-500 transition-colors">Settings</a>
            </div>
        </div>
    </nav>

    <main>
        {#if path === "" || path === "#"}
            <Devices />
        {:else if path === "#settings"}
            <Settings />
        {/if}
    </main>
</div>
