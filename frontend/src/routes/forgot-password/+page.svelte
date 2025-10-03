<script lang="ts">
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { urlBase } from "../../stores/environment";

    let email = $state('');
    let text = $state('');
    let error = $state('');
    let isSubmitting = $state(false);

    async function handleResetPassword() {
        isSubmitting = true;
        error = '';
        text = '';

        if (!email) {
            error = 'Email is required.';
            isSubmitting = false;
            return;
        }

        try {
            const response = await fetch(`${$urlBase}/api/auth/reset-password`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email }),
            });

            if (response.status === 404) {
                error = 'Email not found. Please check your email address.';
            } else if (response.status === 400) {
                error = 'Invalid email format.';
            } else if (!response.ok) {
                error = 'We\'re experiencing technical issues, please try again later.';
            } else {
                text = 'We have sent a reset link to your email.';
            }
        } catch (err) {
            error = 'We\'re experiencing technical issues, please try again later.';
        }
        isSubmitting = false;
    }
</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="Forgot Password" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        <div class="outline outline-indigo-800 outline-2 rounded-3xl w-2/3 max-w-2xl flex flex-col items-center p-8">
            <h2 class="mb-8 text-3xl font-bold text-white">Forgot Password</h2>
            
            <div class="w-full space-y-6">
                <p class="text-gray-400 text-center">We will send a reset link to your email.</p>
                
                <!-- Email Input -->
                <div class="w-full">
                    <label for="email" class="block text-xl font-semibold text-white mb-3">Email</label>
                    <input 
                        id="email"
                        type="email" 
                        placeholder="Enter your email..." 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 autofill:bg-gray-700 autofill:text-white"
                        bind:value={email}
                        onkeydown={(e) => e.key === 'Enter' && handleResetPassword()}
                    />
                </div>
                
                <!-- Error Message -->
                {#if error}
                    <div class="w-full text-center">
                        <p class="text-red-400 text-md">{error}</p>
                    </div>
                {/if}
                
                <!-- Success Message -->
                {#if text}
                    <div class="w-full text-center">
                        <p class="text-green-400 text-md">{text}</p>
                    </div>
                {/if}
                
                <!-- Send Reset Link Button -->
                <div class="w-full flex justify-center pt-4">
                    <button 
                        class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors {isSubmitting ? 'animate-pulse' : ''} disabled:opacity-50 disabled:cursor-not-allowed" 
                        onclick={handleResetPassword} 
                        disabled={!email}
                    >
                        {#if isSubmitting}
                            Sending...
                        {:else}
                            Send Reset Link
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </section>
</main>