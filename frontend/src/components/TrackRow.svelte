<script>
    import LikeButton from "./LikeButton.svelte";
    import LayerrButton from "./LayerrButton.svelte";

    // Track Details
    let imageURL = $state('')
    let trackName = $state('Yessirski')
    let artistName = $state('trentjkelly')
    let timeAgo = $state('3 days ago')
    let trackLength = $state('3:13')
    let trackId = $state(24)
    let numLikes = $state(3)
    let numLayerrs = $state(4)

</script>

<div class="w-full h-20 bg-slate-400 flex flex-row items-center">

    <div class="h-20 w-20 flex flex-row items-center justify-center">
        <img src="/play.png" class="h-10 w-10" alt="Play button">
    </div>

    <div class="h-20 w-20 bg-slate-200">
        <img src={imageURL} alt="Cover art">
    </div>

    <p class="text-xl font-bold ml-8">{trackName}</p>
    <p class="text-lg ml-4">@{artistName}</p>

    <div class="flex flex-row items-center ml-8">
        <img class="h-6 w-6" src="/calendar.png" alt="Calendar icon">
        <p class="text-lg ml-2">{timeAgo}</p>
    </div>

    <div class="flex flex-row items-center ml-8">
        <img class="h-6 w-6" src="/clock.png" alt="clock icon">
        <p class="text-lg ml-2">{trackLength}</p>
    </div>

    <LikeButton trackId={trackId} numLikes={numLikes} />
    <LayerrButton trackId={trackId} numLayerrs={numLayerrs} />

</div>
