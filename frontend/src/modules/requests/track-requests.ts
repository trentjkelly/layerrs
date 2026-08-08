import { logger } from "../lib/logger";
import { fetchWithAuth } from "../lib/fetch";
import type { TrackData, TrackInfo, TrackUsesBulkResponse } from "../../models/types";
import { audio } from "../../stores/player";

// Requests the metadata for the track
export async function getTrackData(urlBase: string, trackId: string): Promise<TrackData | null> {
    try {        
        const baseUrl = `${urlBase}/api/track/${trackId}/data`;
        const response = await fetch(baseUrl, { method: "GET"});
        if (!response.ok) {
            throw new Error("Failed to get track data");
        }
        const responseData = await response.json();

        const trackData: TrackData = {
            description: responseData.description,
            artistId: responseData.artistId,
            likes: responseData.likes,
            layerrs: responseData.layerrs,
            waveformData: responseData.waveformData,
            duration: responseData.trackDuration
        }
        
        return trackData;

    } catch (error) {
        logger.error(`Error catching track data: ${error}`);
        return null;
    }
}

// Requests recommendation-shaped data for a single track
export async function getTrackTrackInfo(urlBase: string, trackId: string): Promise<TrackInfo | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/track/${trackId}/recommendation`, { method: "GET" });
        if (!response.ok) {
            throw new Error("Failed to get track recommendation");
        }
        return await response.json() as TrackInfo;
    } catch (error) {
        logger.error(`Error getting track recommendation: ${error}`);
        return null;
    }
}

// Requests full TrackInfo for a batch of track IDs
export async function getTrackInfoBatch(urlBase: string, trackIds: number[]): Promise<TrackInfo[] | null> {
    if (trackIds.length === 0) {
        return [];
    }

    try {
        const response = await fetchWithAuth(`${urlBase}/api/track/batch`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ track_ids: trackIds })
        });
        if (!response.ok) {
            throw new Error("Failed to get track info batch");
        }
        return await response.json() as TrackInfo[];
    } catch (error) {
        logger.error(`Error getting track info batch: ${error}`);
        return null;
    }
}

// Requests the parent track IDs for a batch of track IDs
export async function getTrackUsesBulk(urlBase: string, trackIds: number[]): Promise<TrackUsesBulkResponse | null> {
    if (trackIds.length === 0) {
        return {};
    }

    try {
        const response = await fetch(`${urlBase}/api/tracks/uses?trackIds=${trackIds.join(",")}`, { method: "GET" });
        if (!response.ok) {
            throw new Error("Failed to get track uses bulk");
        }
        return await response.json() as TrackUsesBulkResponse;
    } catch (error) {
        logger.error(`Error getting track uses bulk: ${error}`);
        return null;
    }
}

// Requests the audio for the track
export async function getAudio(urlBase: string, trackId: string) {
    try {
        const baseUrl = `${urlBase}/api/track/${trackId}/audio`;
        const response = await fetch(baseUrl, { method: "GET"});
        if (!response.ok) {
            throw new Error("Failed to get audio");
        }
        const responseData = await response.json();
        const url = responseData.url;
        return url;

    } catch (error) {
        logger.error(`Error getting audio: ${error}`);
        return null;
    }
}

// Records a play for the track (requires auth; logged-out users are ignored)
export async function recordTrackPlay(urlBase: string, trackId: number): Promise<void> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/track/${trackId}/play`, {
            method: "POST"
        });
        if (!response.ok) {
            throw new Error("Failed to record track play");
        }
    } catch (error) {
        logger.error(`Error recording track play: ${error}`);
    }
}

 