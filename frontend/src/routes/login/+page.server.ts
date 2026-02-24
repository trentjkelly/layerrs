import { redirect } from '@sveltejs/kit';

// Redirect already-logged-in users using layout cookie data (no client-store race).
export async function load({ parent }) {
	const data = await parent();
	if (data.newJWT || data.newRefreshToken) {
		throw redirect(302, '/');
	}
	return {};
}
