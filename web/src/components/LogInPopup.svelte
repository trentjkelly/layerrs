<script>
    import { goto } from "$app/navigation";

    function navigateLogin() {
        goto('/login')
    }

</script>

<div class="w-96 h-48 flex flex-col items-center justify-center mt-40 outline outline-gray-600 rounded rounded-3xl">
    <h3 class="mb-8 text-lg">Not sure who you are. Let's fix that.</h3>
    <button class="w-48 h-12 rounded rounded-full bg-indigo-500 text-xl" onclick={navigateLogin}>Log in</button>
</div>