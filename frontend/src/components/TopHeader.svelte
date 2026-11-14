<script>
    let { pageName, pageIcon } = $props();
</script>

{#if pageName}
<section class="flex h-20 w-full flex-row items-center px-4 sm:px-6 lg:px-8">
    <div class="flex flex-row items-center text-2xl font-bold">
        {#if pageIcon !== ""}
            <img class="h-6 w-6 mr-2" src={pageIcon} alt="">
        {/if}
        <h1>{ pageName }</h1>
    </div>
</section>
{/if}
