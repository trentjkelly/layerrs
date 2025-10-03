<script lang="ts">
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { goto } from "$app/navigation";
    import { logger } from "../../modules/lib/logger";
    import { handleBrowserLogin } from "../../modules/lib/session";
    import { loginServerRequest } from "../../modules/requests/auth-requests";

    let email = $state('')
    let password = $state('')
    let error = $state('')
    let isSubmitting = $state(false);

    async function handleLogin() {
        // Backend authentication request for logging in
        isSubmitting = true;

        const isValid = validateLoginInputs(email, password)
        if (!isValid) {
            error = 'Email and password are required.'
            isSubmitting = false;
            return
        }

        const res = await loginServerRequest(email, password)
        if (res === null) {
            logger.error('Failed to login')
            error = 'We\'re experiencing technical issues, please try again later.'
            isSubmitting = false;
            return
        }

        // Check status code from response
        const status = res.status
        if (status == 401 || status == 400) {
            error = 'Invalid email or password, please try again.'
            isSubmitting = false;
            return
        } else if (res.status !== 200) {
            error = 'We\'re experiencing technical issues, please try again later.'
            isSubmitting = false;
            return
        }

        // Set cookies for the refresh and jwt tokens
        const resJson = await res.json()
        const success = await handleBrowserLogin(resJson.token, resJson.refreshToken)
        if (success) {
            goto('/')
        } else {
            logger.error('Failed to login')
        }

        isSubmitting = false;
    }

    function validateLoginInputs(email : string, password : string) {
        if (email === '' || password === '') {
            return false
        }
        return true
    }

    function navigateSignUp() {
        goto('/signup')
    }

    function navigateForgotPassword() {
        goto('/forgot-password')
    }

</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>

    <TopHeader pageName="Log in" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        <div class="outline outline-indigo-800 outline-2 rounded-3xl w-2/3 max-w-2xl flex flex-col items-center p-8">
            <h2 class="mb-8 text-3xl font-bold text-white">Log In</h2>
            
            <div class="w-full space-y-6">
                <!-- Email Input -->
                <div class="w-full">
                    <label for="email" class="block text-xl font-semibold text-white mb-3">Email</label>
                    <input 
                        id="email" 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border {error ? 'border-red-500' : 'border-gray-600'} focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 autofill:bg-gray-700 autofill:text-white" 
                        type="email" 
                        bind:value={email}
                        placeholder="Enter your email..."
                        onkeydown={(e) => e.key === 'Enter' && handleLogin()}
                    />
                </div>

                <!-- Password Input -->
                <div class="w-full">
                    <label for="password" class="block text-xl font-semibold text-white mb-3">Password</label>
                    <input 
                        id="password" 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border {error ? 'border-red-500' : 'border-gray-600'} focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 autofill:bg-gray-700 autofill:text-white" 
                        type="password" 
                        bind:value={password}
                        placeholder="Enter your password..."
                        onkeydown={(e) => e.key === 'Enter' && handleLogin()}
                    />
                </div>
                
                <!-- Error Message -->
                {#if error}
                    <div class="w-full text-center">
                        <p class="text-red-400 text-md">{error}</p>
                    </div>
                {/if}
                
                <!-- Login Button -->
                <div class="w-full flex justify-center pt-4">
                    <button 
                        class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors {isSubmitting ? 'animate-pulse' : ''} disabled:opacity-50 disabled:cursor-not-allowed" 
                        onclick={handleLogin}
                        disabled={!email || !password}
                    >
                    {#if isSubmitting}
                        Logging in...
                    {:else}
                        Log In
                    {/if}
                    </button>
                </div>
                
                <!-- Sign Up Link -->
                <div class="w-full flex justify-center pt-4">
                    <button 
                        class="text-gray-400 hover:text-white transition-colors underline" 
                        onclick={navigateSignUp}
                    >
                        Don't have an account? Sign up instead
                    </button>
                </div>

                <!-- Forgot Password Link -->
                <div class="w-full flex justify-center">
                    <button 
                        class="text-gray-400 hover:text-white transition-colors underline" 
                        onclick={navigateForgotPassword}
                    >
                        Forgot password?
                    </button>
                </div>
            </div>
        </div>
    </section>
</main>