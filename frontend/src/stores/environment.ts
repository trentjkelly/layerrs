import { get, writable } from "svelte/store"
import { logger } from "../modules/lib/logger"

function resolveEnvironment(): { env: string; url: string } {
	const envValue = import.meta.env.VITE_ENVIRONMENT
	const isDev = import.meta.env.DEV
	if (envValue === 'DEVELOPMENT' || (envValue == null && isDev)) {
		return { env: 'DEVELOPMENT', url: 'http://localhost:8080' }
	} else if (envValue === 'PRODUCTION' || (envValue == null && !isDev)) {
		return { env: 'PRODUCTION', url: 'https://layerrs.com' }
	} else {
		logger.error('Could not set environment')
		return { env: 'PRODUCTION', url: 'https://layerrs.com' }
	}
}

const { env, url } = resolveEnvironment()
export const environment = writable(env)
export const urlBase = writable(url)


export function getEnvironment(): string {
	return get(environment)
}

export function getUrlBase(): string {
	return get(urlBase)
}