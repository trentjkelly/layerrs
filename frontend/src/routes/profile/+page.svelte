<script lang="ts">
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";

    let username = $state("sampleartist");
    let bio = $state("Producer from Chicago. Making beats since 2018. Influences: J Dilla, Madlib, Flying Lotus.");
    let email = "artist@example.com";

    let profilePhotoSrc = $state<string | null>(null);
    let photoFile = $state<File | null>(null);

    function handlePhotoChange(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0] ?? null;
        photoFile = file;
        profilePhotoSrc = file ? URL.createObjectURL(file) : null;
    }

    function removeNewPhoto() {
        photoFile = null;
        profilePhotoSrc = null;
    }

    function saveChanges() {
        // TODO: wire up API request
    }
</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="Profile" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        <div class="outline outline-gray-600 rounded-3xl w-2/3 max-w-4xl flex flex-col items-center p-8">
            <h2 class="mb-8 text-3xl font-bold text-white">Edit Profile</h2>

            <!-- Email (read-only) -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Email</h3>
                <p class="w-full px-3 py-2 rounded-lg bg-gray-800 text-gray-400 border border-gray-700 select-none">
                    {email}
                </p>
            </div>

            <!-- Username -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Username</h3>
                <input
                    class="w-full px-3 py-2 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
                    type="text"
                    bind:value={username}
                    placeholder="Your username"
                    maxlength={30}
                />
                <p class="text-sm mt-1 text-gray-400">{username.length}/30 characters</p>
            </div>

            <!-- Bio -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Bio</h3>
                <textarea
                    class="w-full px-3 py-2 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 resize-none"
                    rows={4}
                    bind:value={bio}
                    placeholder="Tell people a bit about yourself..."
                    maxlength={300}
                ></textarea>
                <p class="text-sm mt-1 text-gray-400">{bio.length}/300 characters</p>
            </div>

            
            <!-- Profile Photo -->
            <div class="w-full mb-12 flex flex-col">
                <h3 class="text-xl font-semibold text-white mb-1">Profile Photo</h3>
                <div class="flex flex-row h-20 w-full items-center items-end gap-4">
                    <div class="h-20 w-20 shrink-0 rounded-xl bg-gray-700 border border-gray-600 overflow-hidden">
                        {#if profilePhotoSrc}
                            <img src={profilePhotoSrc} alt="Profile preview" class="h-full w-full object-cover" />
                        {/if}
                    </div>
                    <label for="photo" class="cursor-pointer h-10 ml-4 px-4 py-2 text-sm font-semibold bg-indigo-600 hover:bg-indigo-700 border border-gray-600 rounded-lg text-white transition-colors flex items-center">
                        Change / Add Photo
                    </label>
                    <input id="photo" class="hidden" type="file" accept="image/*" onchange={handlePhotoChange} />
                </div>
            </div>

            <!-- Save Buttons -->
            <button class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors">Save Changes</button>
        </div>
    </section>
</main>
