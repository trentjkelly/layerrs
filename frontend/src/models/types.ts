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

export type Recommendation = {
    id: number;
    description: string;
    artistId: number;
    artistName: string;
    likes: number;
    layerrs: number;
    duration: number;
    waveformData: number[];
    isLiked: boolean;
}