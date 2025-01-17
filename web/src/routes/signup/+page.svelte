<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { goto } from "$app/navigation";

    let email = $state('')
    let password = $state('')
    let name = $state('')
    let username = $state('')


    async function signup() {
        if ((email !== '') && (password !== '') && (name !== '') && (username !== '')) {
            try {
                const res = await fetch(`http://localhost:8080/api/authentication/signup`, {
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
                console.error(err)
            }
        }
    }

    function navigateLogIn() {
        goto('/login')
    }
</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'}`}>

    <TopHeader pageName="Sign up" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center">

        <div class="border border-gray-600 bg-gray-800 rounded-3xl w-1/2 flex flex-col items-center">
            <h2 class="my-8 text-4xl">Sign Up</h2>

            <div class="flex flex-col">
                <label for="email">Email:</label>
                <input id="email" class="text-black" type="text" bind:value={email}>

                <label for="password">Password:</label>
                <input id="password" class="text-black" type="text" bind:value={password}>

                <label for="username">Username:</label>
                <input id="username" class="text-black" type="text" bind:value={username}>

                <label for="name">Name:</label>
                <input id="name" class="text-black" type="text" bind:value={name}>

                <button onclick={signup} class="h-20 w-20 bg-indigo-400">Sign up!</button>
            </div>

            <button onclick={navigateLogIn} class="h-20 w-20">Log in instead...</button>
        </div>
        
    </section>
</main>