import { logger } from "../lib/logger";
import { fetchWithAuth } from "../lib/fetch";
import type { Page, PageSubmission, PageTrack, PageWithFollowerCount, TrackPagesResponse, TrackPagesBulkResponse } from "../../models/types";

export async function createPage(urlBase: string, name: string, description: string): Promise<Page | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name, description })
        });
        if (!response.ok) {
            throw new Error(`Failed to create page: ${response.statusText}`);
        }
        return await response.json() as Page;
    } catch (error) {
        logger.error(`Error creating page: ${error}`);
        return null;
    }
}

export async function getPage(urlBase: string, pageId: string): Promise<Page | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get page: ${response.statusText}`);
        }
        return await response.json() as Page;
    } catch (error) {
        logger.error(`Error getting page: ${error}`);
        return null;
    }
}

export async function updatePage(urlBase: string, pageId: string, name: string, description: string): Promise<Page | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name, description })
        });
        if (!response.ok) {
            throw new Error(`Failed to update page: ${response.statusText}`);
        }
        return await response.json() as Page;
    } catch (error) {
        logger.error(`Error updating page: ${error}`);
        return null;
    }
}

export async function deletePage(urlBase: string, pageId: string): Promise<boolean> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}`, { method: "DELETE" });
        return response.ok;
    } catch (error) {
        logger.error(`Error deleting page: ${error}`);
        return false;
    }
}

export async function getArtistPages(urlBase: string, artistId: string): Promise<Page[]> {
    try {
        const response = await fetch(`${urlBase}/api/artists/${artistId}/pages`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get artist pages: ${response.statusText}`);
        }
        return await response.json() as Page[];
    } catch (error) {
        logger.error(`Error getting artist pages: ${error}`);
        return [];
    }
}

export async function followPage(urlBase: string, pageId: string): Promise<boolean> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/follow`, { method: "POST" });
        return response.ok;
    } catch (error) {
        logger.error(`Error following page: ${error}`);
        return false;
    }
}

export async function unfollowPage(urlBase: string, pageId: string): Promise<boolean> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/follow`, { method: "DELETE" });
        return response.ok;
    } catch (error) {
        logger.error(`Error unfollowing page: ${error}`);
        return false;
    }
}

export async function addTrackToPage(urlBase: string, pageId: string, trackId: number): Promise<PageTrack | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/tracks`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ trackId })
        });
        if (!response.ok) {
            throw new Error(`Failed to add track to page: ${response.statusText}`);
        }
        return await response.json() as PageTrack;
    } catch (error) {
        logger.error(`Error adding track to page: ${error}`);
        return null;
    }
}

export async function removeTrackFromPage(urlBase: string, pageId: string, trackId: number): Promise<boolean> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/tracks/${trackId}`, { method: "DELETE" });
        return response.ok;
    } catch (error) {
        logger.error(`Error removing track from page: ${error}`);
        return false;
    }
}

export async function getPageFeed(urlBase: string, pageId: string, limit: number = 20, offset: number = 0): Promise<PageTrack[]> {
    try {
        const response = await fetch(`${urlBase}/api/pages/${pageId}/feed?limit=${limit}&offset=${offset}`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get page feed: ${response.statusText}`);
        }
        return await response.json() as PageTrack[];
    } catch (error) {
        logger.error(`Error getting page feed: ${error}`);
        return [];
    }
}

export async function getSubmissions(urlBase: string, pageId: string, limit: number = 20, offset: number = 0): Promise<PageSubmission[]> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/submissions?limit=${limit}&offset=${offset}`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get submissions: ${response.statusText}`);
        }
        return await response.json() as PageSubmission[];
    } catch (error) {
        logger.error(`Error getting submissions: ${error}`);
        return [];
    }
}

export async function submitTrack(urlBase: string, pageId: string, trackId: number, note: string = ""): Promise<PageSubmission | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/submissions`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ trackId, note })
        });
        if (!response.ok) {
            throw new Error(`Failed to submit track: ${response.statusText}`);
        }
        return await response.json() as PageSubmission;
    } catch (error) {
        logger.error(`Error submitting track: ${error}`);
        return null;
    }
}

export async function approveSubmission(urlBase: string, pageId: string, submissionId: number, note: string = ""): Promise<PageTrack | null> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/pages/${pageId}/submissions/${submissionId}/approve`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ note })
        });
        if (!response.ok) {
            throw new Error(`Failed to approve submission: ${response.statusText}`);
        }
        return await response.json() as PageTrack;
    } catch (error) {
        logger.error(`Error approving submission: ${error}`);
        return null;
    }
}

export async function getTrackPages(urlBase: string, trackId: string): Promise<TrackPagesResponse | null> {
    try {
        const response = await fetch(`${urlBase}/api/tracks/${trackId}/pages`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get track pages: ${response.statusText}`);
        }
        return await response.json() as TrackPagesResponse;
    } catch (error) {
        logger.error(`Error getting track pages: ${error}`);
        return null;
    }
}

export async function getTrackPagesBulk(urlBase: string, trackIds: number[]): Promise<Record<string, TrackPagesBulkResponse> | null> {
    if (trackIds.length === 0) {
        return {};
    }

    try {
        const response = await fetch(`${urlBase}/api/tracks/pages?trackIds=${trackIds.join(",")}`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get track pages bulk: ${response.statusText}`);
        }
        return await response.json() as Record<string, TrackPagesBulkResponse>;
    } catch (error) {
        logger.error(`Error getting track pages bulk: ${error}`);
        return null;
    }
}

export async function getFollowingFeed(urlBase: string, limit: number = 20, offset: number = 0): Promise<PageTrack[]> {
    try {
        const response = await fetchWithAuth(`${urlBase}/api/feed/following?limit=${limit}&offset=${offset}`, { method: "GET" });
        if (!response.ok) {
            throw new Error(`Failed to get following feed: ${response.statusText}`);
        }
        return await response.json() as PageTrack[];
    } catch (error) {
        logger.error(`Error getting following feed: ${error}`);
        return [];
    }
}
