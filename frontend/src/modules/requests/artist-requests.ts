import { logger } from "../lib/logger";
import type { ArtistData } from "../../models/types";

// Gets the name of the artist
export async function getArtistName(urlBase: string, artistId: number): Promise<ArtistData | null> {
    try {
        const baseUrl = `${urlBase}/api/artist/${artistId}`;
        const response = await fetch(baseUrl, {
            method: "GET"
        })
        if(!response.ok) {
            throw new Error("Failed to get artist data");
        }

        const responseData = await response.json();
        const artistData: ArtistData = {
            name: responseData.name
        }

        return artistData;
    } catch (error) {
        logger.error(`Could not retrieve artist data: ${error}`);
        return null;
    }
}