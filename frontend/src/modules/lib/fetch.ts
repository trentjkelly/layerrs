import { get } from 'svelte/store';
import { jwt, refreshToken } from '../../stores/auth';
import { handleBrowserLogin, handleBrowserLogout } from './session';
import { getUrlBase } from '../../stores/environment';
import { logger } from './logger';

async function attemptTokenRefresh(): Promise<boolean> {
    const currentRefreshToken = get(refreshToken);
    if (!currentRefreshToken) return false;

    try {
        const res = await fetch(`${getUrlBase()}/api/authentication/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refreshToken: currentRefreshToken })
        });

        if (!res.ok) {
            await handleBrowserLogout();
            return false;
        }

        const data = await res.json();
        // Refresh endpoint only returns a new JWT; re-use the existing refresh token
        await handleBrowserLogin(data.token, currentRefreshToken);
        return true;
    } catch (error) {
        logger.error(`Token refresh failed: ${error}`);
        return false;
    }
}

export async function fetchWithAuth(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    const currentJwt = get(jwt);

    const headers = new Headers(init?.headers);
    if (currentJwt) headers.set('Authorization', `Bearer ${currentJwt}`);

    const response = await fetch(input, { ...init, headers });

    if (response.status === 401) {
        const refreshed = await attemptTokenRefresh();
        if (refreshed) {
            headers.set('Authorization', `Bearer ${get(jwt)}`);
            return fetch(input, { ...init, headers });
        }
    }

    return response;
}
