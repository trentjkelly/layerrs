import { writable, get } from 'svelte/store';
import { urlBase } from './environment';
import { logger } from '../modules/lib/logger';
import { fetchWithAuth } from '../modules/lib/fetch';

export const usernameStore = writable('');
export const emailStore = writable('');
export const bioStore = writable('');
export const portraitUrlStore = writable('');
export const validStripeSellerStore = writable(false);

export async function loadProfile(): Promise<void> {
    if (get(usernameStore) !== '') return;

    try {
        const res = await fetchWithAuth(`${get(urlBase)}/api/profile`, { method: 'GET' });
        if (res.ok) {
            const data = await res.json();
            usernameStore.set(data.username ?? '');
            emailStore.set(data.email ?? '');
            bioStore.set(data.bio ?? '');
            portraitUrlStore.set(data.portraitUrl ?? '');
            validStripeSellerStore.set(data.ValidStripeSeller ?? false)
        }
    } catch (err) {
        logger.error(`Failed to load profile: ${err}`);
    }
}
