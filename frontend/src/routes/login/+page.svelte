<script lang="ts">
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { logger } from "../../modules/lib/logger";
    import { loginServerRequest } from "../../modules/requests/auth-requests";
    import { isLoggedIn } from "../../stores/auth";
    import { goto, invalidateAll } from '$app/navigation';
    import { browser } from '$app/environment';

    let email = $state('')
    let error = $state('')
    let isSubmitting = $state(false);
    let emailEntered = $state(false);
    let timeElapsed = $state(0);

    let sendButtonDisabled = $state(false);

    const REFRESH_TIME = 30;

    // Just a check if they're already logged in, they shouldn't be allowed on the login page
    $effect(() => {
        if ($isLoggedIn) {
            goto('/');
        }
    });

    // For if a user switches back to this tab (usually after checking login email)
    // refreshes the page so they don't see login page
    $effect(() => {
        if (!browser) return;

        const handler = () => {
            if (document.visibilityState === 'visible') {
                invalidateAll();
            }
        };

        document.addEventListener('visibilitychange', handler);
        return () => document.removeEventListener('visibilitychange', handler);
    });

    async function handleLogin() {
        isSubmitting = true;
        sendButtonDisabled = true;

        const isValid = validateEmailInput(email)
        if (!isValid) {
            error = 'Email is required.'
            isSubmitting = false;
            return
        }

        const res = await loginServerRequest(email)
        if (res === null) {
            logger.error('Failed to login')
            error = 'We\'re experiencing technical issues, please try again later.'
            isSubmitting = false;
            return
        }

        // Check status code from response
        if (res.status == 200) {
            isSubmitting = false;
            emailEntered = true;
        } else {
            isSubmitting = false;
            error = 'We\'re experiencing technical issues, please try again later.'
            return
        }

        startCountdown();
    }

    function validateEmailInput(email : string) {
        if (email === '') {
            return false
        }
        return true
    }

    function startCountdown() {
        timeElapsed = REFRESH_TIME

        const interval = setInterval(() => {
            timeElapsed -= 1;
            if (timeElapsed <= 0) {
                clearInterval(interval);
                timeElapsed = REFRESH_TIME;
                sendButtonDisabled = false;
            }
        }, 1000);

    }

</script>

<main class={`transition-all duration-300 h-screen w-full flex flex-col ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="" pageIcon=""/>
    <section class="w-full flex-1 flex flex-col justify-center items-center">
        <div class="w-full flex flex-col items-center">
            <h2 class="mb-8 text-2xl font-bold text-white">Log in to your <span class="text-violet-500">Layerrs</span> Account</h2>

            <!-- Email Input -->
            <div class="w-96 mb-4">
                <label for="email" class="block text-xl font-semibold text-white mb-3">Email</label>
                <input 
                    id="email" 
                    class="w-full px-2 py-2 rounded-lg bg-zinc-900 text-white placeholder-zinc-200 focus:outline-none border {error ? 'border-red-500' : 'border-zinc-200'}" 
                    type="email" 
                    bind:value={email}
                    placeholder="Enter your email..."
                    onkeydown={(e) => e.key === 'Enter' && handleLogin()}
                />
            </div>

            {#if emailEntered}
                <p class="mb-2">Please check your email for a link to log in.</p>
                {#if sendButtonDisabled}
                    <p class="text-sm text-zinc-400 mb-2">You can resend the link in {timeElapsed} seconds</p>
                {/if}
                <button class="py-2 px-6 rounded rounded-lg bg-violet-500 hover:bg-violet-600 disabled:opacity-50 disabled:cursor-not-allowed" onclick={handleLogin} disabled={sendButtonDisabled}>Resend Link</button>
            {:else}
                <button class="py-2 px-6 rounded rounded-lg bg-violet-500 hover:bg-violet-600" onclick={handleLogin}>Continue</button>
            {/if}
        </div>
    </section>
</main>