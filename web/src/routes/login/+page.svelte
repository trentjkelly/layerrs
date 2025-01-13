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

<main class={`transition-all duration-300 h-full w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'}`}>

    <TopHeader pageName="Log in" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center">

        <div class="border border-gray-600 bg-gray-800 rounded-3xl w-1/2 flex flex-col items-center">
            <h2 class="my-8 text-4xl">Log In</h2>
            <div class="flex flex-col">
                <label for="email">Email:</label>
                <input id="email" class="text-black" type="text" bind:value={email}>

                <label for="password">Password:</label>
                <input id="password" class="text-black" type="text" bind:value={password}>

                <button onclick={login} class="h-20 w-20 bg-indigo-400">Log in!</button>
            </div>
            <button onclick={navigateSignUp} class="h-20 w-20">Sign up instead....</button>
        </div>
        
    </section>
</main>