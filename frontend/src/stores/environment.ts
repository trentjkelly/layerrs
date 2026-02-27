import { get, writable } from "svelte/store"
import { logger } from "../modules/lib/logger"

// Defaults to production values
export const environment = writable('PRODUCTION')
export const urlBase = writable('https://layerrs.com')

// Changes the environment based on the .env file (or Vite dev mode when .env is missing)
export async function handleEnvironment(): Promise<void> {
	const envValue = import.meta.env.VITE_ENVIRONMENT
	const isDev = import.meta.env.DEV
	if (envValue === 'DEVELOPMENT' || (envValue == null && isDev)) {
		urlBase.set('http://localhost:8080')
		environment.set('DEVELOPMENT')
	} else if (envValue === 'PRODUCTION' || (envValue == null && !isDev)) {
		urlBase.set('https://layerrs.com')
		environment.set('PRODUCTION')
	} else {
		logger.error('Could not set environment')
	}
}

export function getEnvironment(): string {
	return get(environment)
}

export function getUrlBase(): string {
	return get(urlBase)
}