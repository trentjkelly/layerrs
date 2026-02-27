import { logger } from "../lib/logger";
import type { TrackData } from "../../models/types";
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
            name: responseData.name,
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

 