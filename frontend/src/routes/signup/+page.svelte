<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { goto } from "$app/navigation";
    import { urlBase } from "../../stores/environment";
    import { logger } from "../../modules/lib/logger";

    let email = $state('')
    let password = $state('')
    let name = $state('')
    let username = $state('')


    async function signup() {
        if ((email !== '') && (password !== '') && (name !== '') && (username !== '')) {
            try {
                const res = await fetch(`${$urlBase}/api/authentication/signup`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        email: email,
                        password: password,
                        name: name,
                        username: username
                    })
                })
            } catch (err) {
                logger.error(err)
            }
        }
    }

    function navigateLogIn() {
        goto('/login')
    }
    
</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>

    <TopHeader pageName="Sign up" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        <div class="outline outline-gray-600 rounded-3xl w-2/3 max-w-2xl flex flex-col items-center p-8">
            <h2 class="mb-8 text-3xl font-bold text-white">Sign Up</h2>

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

                <!-- Username Input -->
                <div class="w-full">
                    <label for="username" class="block text-xl font-semibold text-white mb-3">Username</label>
                    <input 
                        id="username" 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" 
                        type="text" 
                        bind:value={username}
                        placeholder="Choose a username..."
                    />
                </div>

                <!-- Name Input -->
                <div class="w-full">
                    <label for="name" class="block text-xl font-semibold text-white mb-3">Full Name</label>
                    <input 
                        id="name" 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" 
                        type="text" 
                        bind:value={name}
                        placeholder="Enter your full name..."
                    />
                </div>
                
                <!-- Sign Up Button -->
                <div class="w-full flex justify-center pt-4">
                    <button 
                        class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed" 
                        onclick={signup}
                        disabled={!email || !password || !username || !name}
                    >
                        Sign Up
                    </button>
                </div>
                
                <!-- Log In Link -->
                <div class="w-full flex justify-center pt-4">
                    <button 
                        class="text-gray-400 hover:text-white transition-colors underline" 
                        onclick={navigateLogIn}
                    >
                        Already have an account? Log in instead
                    </button>
                </div>
            </div>
        </div>
    </section>
</main>