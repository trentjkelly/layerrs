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
    username?: string;
}

export type TrackInfo = {
    id: number;
    description: string;
    artistId: number;
    artistName: string;
    artistPortraitUrl: string;
    likes: number;
    plays: number;
    layerrs: number;
    duration: number;
    waveformData: number[];
    isLiked: boolean;
    color: string;
}

export type Page = {
    id: number;
    editorId?: number;
    name: string;
    description: string;
    createdAt: string;
    updatedAt: string;
    followerCount?: number;
    isFollowing?: boolean;
    isEditor?: boolean;
    editorName?: string;
}

export type PageTrack = {
    id: number;
    pageId: number;
    trackId: number;
    addedAt: string;
    recommenderId?: number;
    track?: TrackInfo;
    notes?: PageTrackNote[];
}

export type PageSubmission = {
    id: number;
    pageId: number;
    trackId: number;
    submitterId: number;
    note: string;
    createdAt: string;
    track?: TrackInfo;
    submitter?: ArtistData;
}

export type PageTrackNote = {
    id: number;
    pageTrackId: number;
    note: string;
    createdAt: string;
}

export type PageWithFollowerCount = {
    id: number;
    editorId?: number;
    name: string;
    description: string;
    createdAt: string;
    updatedAt: string;
    followerCount: number;
    editorName?: string;
}

export type TrackPagesResponse = {
    topPage?: PageWithFollowerCount;
    otherCount: number;
    pages: PageWithFollowerCount[];
}

export type TrackPagesBulkResponse = {
    pages: PageWithFollowerCount[];
    pageCount: number;
}

export type TrackUsesBulkResponse = {
    [trackId: number]: number[];
}