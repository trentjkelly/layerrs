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
    <section class="w-full min-h-full flex flex-row justify-center items-center">
        <div class="w-full flex flex-col items-center">
            <h2 class="mb-8 text-2xl font-bold text-white">Log in to your <span class="text-violet-500">Layerrs</span> Account</h2>

            <!-- Email Input -->
            <div class="w-96 mb-6">
                <label for="email" class="block text-xl font-semibold text-white mb-3">Email</label>
                <input 
                    id="email" 
                    class="w-full px-2 py-2 rounded-lg bg-zinc-900 text-white placeholder-gray-200 border {error ? 'border-red-500' : 'border-gray-200'}" 
                    type="email" 
                    bind:value={email}
                    placeholder="Enter your email..."
                    onkeydown={(e) => e.key === 'Enter' && handleLogin()}
                />
            </div>

            <button class="py-2 px-6 bg-violet-500 hover:bg-violet-600">Continue</button>
        </div>
    </section>
</main>