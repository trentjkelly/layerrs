export type TrackData = {
    description: string;
    artistId: string;
    likes: number;
    layerrs: number;
    waveformData: number[];
    duration: number;
}

export type ArtistData = {
    name: string;
}

export type TrackInfo = {
    id: number;
    description: string;
    artistId: number;
    artistName: string;
    artistPortraitUrl: string;
    likes: number;
    layerrs: number;
    duration: number;
    waveformData: number[];
    isLiked: boolean;
    color: string;
}