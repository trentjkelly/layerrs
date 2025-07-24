import { get, writable } from "svelte/store"

// Defaults to production values
export const environment = writable('PRODUCTION')
export const urlBase = writable('https://layerrs.com')

// Changes the environment based on the .env file
export async function handleEnvironment(): Promise<void> {
	const environment = import.meta.env.VITE_ENVIRONMENT
	if (environment === 'DEVELOPMENT') {
		urlBase.set('http://localhost:8080')
		environment.set('DEVELOPMENT')
	} else {
		urlBase.set('https://layerrs.com')
		environment.set('PRODUCTION')
	}
}

export function getEnvironment(): string {
	return get(environment)
}