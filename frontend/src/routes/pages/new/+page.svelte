<script lang="ts">
    import { goto } from "$app/navigation";
    import { isLoggedIn } from "../../../stores/auth";
    import { urlBase } from "../../../stores/environment";
    import { createPage } from "../../../modules/requests/page-requests";
    import TopHeader from "../../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../../stores/player";

    let name = $state("");
    let description = $state("");
    let errorMessage = $state("");
    let isSubmitting = $state(false);

    async function handleSubmit(event: SubmitEvent) {
        event.preventDefault();
        errorMessage = "";

        const trimmedName = name.trim();
        if (trimmedName === "") {
            errorMessage = "Page name is required";
            return;
        }
        if (trimmedName.length > 255) {
            errorMessage = "Page name must be 255 characters or less";
            return;
        }

        isSubmitting = true;
        const page = await createPage($urlBase, trimmedName, description.trim());
        isSubmitting = false;

        if (page) {
            goto(`/pages/${page.id}`);
        } else {
            errorMessage = "Failed to create page. Please try again.";
        }
    }
</script>

<main class={`transition-all duration-300 h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-950`}>
    <TopHeader pageName="New Page" pageIcon="" />

    <section class="p-8 max-w-2xl mx-auto">
        {#if !$isLoggedIn}
            <p class="text-center text-lg">You must be logged in to create a Page.</p>
        {:else}
            <h1 class="text-3xl font-bold mb-6">Create a Page</h1>

            <form onsubmit={handleSubmit} class="flex flex-col gap-4 bg-zinc-900 border border-zinc-700 p-6">
                <div>
                    <label for="name" class="block font-semibold mb-1">Name *</label>
                    <input
                        id="name"
                        type="text"
                        bind:value={name}
                        maxlength="255"
                        class="w-full p-2 bg-zinc-800 text-white border border-zinc-700 focus:border-focus outline-none"
                        placeholder="e.g. Underground Hits"
                    />
                </div>

                <div>
                    <label for="description" class="block font-semibold mb-1">Description</label>
                    <textarea
                        id="description"
                        bind:value={description}
                        rows="4"
                        class="w-full p-2 bg-zinc-800 text-white border border-zinc-700 focus:border-focus outline-none"
                        placeholder="What is this Page about?"
                    ></textarea>
                </div>

                {#if errorMessage}
                    <p class="text-danger">{errorMessage}</p>
                {/if}

                <button
                    type="submit"
                    disabled={isSubmitting}
                    class="self-start px-6 py-2 bg-primary hover:bg-primary-hover active:bg-primary-active disabled:opacity-50 font-semibold"
                >
                    {isSubmitting ? "Creating..." : "Create Page"}
                </button>
            </form>
        {/if}
    </section>
</main>
