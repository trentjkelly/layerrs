<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { goto } from "$app/navigation";
    import { urlBase } from "../../stores/environment";
    import { logger } from "../../lib/logger";

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

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Sign up" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center">

        <div class="outline outline-gray-600 rounded-3xl w-1/2 flex flex-col items-center mt-16">
            <h2 class="my-8 text-4xl">Sign Up</h2>

            <div class="flex flex-col">
                <label for="email">Email:</label>
                <input id="email" class="text-black mb-4" type="text" bind:value={email}>

                <label for="password">Password:</label>
                <input id="password" class="text-black mb-4" type="text" bind:value={password}>

                <label for="username">Username:</label>
                <input id="username" class="text-black mb-4" type="text" bind:value={username}>

                <label for="name">Name:</label>
                <input id="name" class="text-black mb-4" type="text" bind:value={name}>

                <div class="w-full h-auto flex flex-row items-center justify-center mt-8">
                    <button onclick={signup} class="h-12 w-32 bg-indigo-500 hover:bg-indigo-400 rounded rounded-full">Sign up!</button>
                </div>
            </div>

            <button onclick={navigateLogIn} class="h-20 w-32 hover:underline mt-8 text-gray-400">Log in instead...</button>
        </div>
        
    </section>
</main>