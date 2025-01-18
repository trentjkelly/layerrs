<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { goto } from "$app/navigation";

    let email = $state('')
    let password = $state('')

    async function login() {
        if ((email !== '') && (password !== '')) {
            try {
                const res = await fetch(`http://localhost:8080/api/authentication/login`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        email: email,
                        password: password
                    })
                })

                const jsonData = await res.json()

                // Set cookies for the refresh and jwt tokens
                try {
                    const res  = await fetch('/cookies', { 
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({ 
                            refreshToken: jsonData.refreshToken,
                            jwtToken: jsonData.token
                        })
                    })
                    if (!res.ok) {
                        console.error("Failed to set the JWT and refresh tokens", res.statusText);
                    }
                } catch (error) {
                    console.error("Failed to set the JWT")
                }

                // Go to homepage when logged in
                if (res.status == 200) {
                    goto('/')
                }

            } catch (err) {
                console.error(err)
            }
        }
    }

    function navigateSignUp() {
        goto('/signup')
    }

</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Log in" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center">

        <div class="outline outline-gray-600 rounded-3xl w-1/2 flex flex-col items-center mt-24">
            <h2 class="my-8 text-4xl">Log In</h2>
            <div class="flex flex-col">
                <label for="email">Email:</label>
                <input id="email" class="text-black mb-4" type="text" bind:value={email}>

                <label for="password">Password:</label>
                <input id="password" class="text-black" type="text" bind:value={password}>
                <div class="w-full h-auto flex flex-row items-center justify-center mt-8">
                    <button onclick={login} class="h-12 w-32 bg-indigo-500 hover:bg-indigo-400 rounded rounded-full">Log in!</button>
                </div>
            </div>
            <button onclick={navigateSignUp} class="h-20 w-32 hover:underline mt-8 text-gray-400">Sign up instead....</button>
        </div>
        
    </section>
</main>