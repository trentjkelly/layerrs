<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { goto } from "$app/navigation";
    import { logger } from "../../modules/lib/logger";
    import { handleBrowserLogin } from "../../modules/lib/session";
    import { loginServerRequest } from "../../modules/requests/auth-requests";

    let email = $state('')
    let password = $state('')

    async function handleLogin() {
        // Backend authentication request for logging in
        const res = await loginServerRequest(email, password)
        if (res === null) {
            logger.error('Failed to login')
            return
        }

        console.log(res)

        // Set cookies for the refresh and jwt tokens
        const success = await handleBrowserLogin(res.token, res.refreshToken)
        if (success) {
            goto('/')
        } else {
            logger.error('Failed to login')
        }
    }

    function navigateSignUp() {
        goto('/signup')
    }

</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Log in" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        <div class="outline outline-gray-600 rounded-3xl w-2/3 max-w-2xl flex flex-col items-center p-8">
            <h2 class="mb-8 text-3xl font-bold text-white">Log In</h2>
            
            <div class="w-full space-y-6">
                <!-- Email Input -->
                <div class="w-full">
                    <label for="email" class="block text-xl font-semibold text-white mb-3">Email</label>
                    <input 
                        id="email" 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" 
                        type="email" 
                        bind:value={email}
                        placeholder="Enter your email..."
                    />
                </div>

                <!-- Password Input -->
                <div class="w-full">
                    <label for="password" class="block text-xl font-semibold text-white mb-3">Password</label>
                    <input 
                        id="password" 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" 
                        type="password" 
                        bind:value={password}
                        placeholder="Enter your password..."
                    />
                </div>
                
                <!-- Login Button -->
                <div class="w-full flex justify-center pt-4">
                    <button 
                        class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed" 
                        onclick={handleLogin}
                        disabled={!email || !password}
                    >
                        Log In
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
            </div>
        </div>
    </section>
</main>