import type { Handle } from '@sveltejs/kit';

const CHROME_DEVTOOLS_PATH = '/.well-known/appspecific/com.chrome.devtools.json';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname === CHROME_DEVTOOLS_PATH) {
		return new Response('{}', {
			headers: { 'Content-Type': 'application/json' }
		});
	}
	return resolve(event);
};
