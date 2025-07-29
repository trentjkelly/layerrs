import { get, writable } from "svelte/store"
import { logger } from "../lib/logger"

// Defaults to production values
export const environment = writable('PRODUCTION')
export const urlBase = writable('https://layerrs.com')

// Changes the environment based on the .env file
export async function handleEnvironment(): Promise<void> {
	const envValue = import.meta.env.VITE_ENVIRONMENT
	if (envValue === 'DEVELOPMENT') {
		urlBase.set('http://localhost:8080')
		environment.set('DEVELOPMENT')
	} else if (envValue === 'PRODUCTION') {
		urlBase.set('https://layerrs.com')
		environment.set('PRODUCTION')
	} else {
		logger.error('Could not set environment')
	}
}

export function getEnvironment(): string {
	return get(environment)
}