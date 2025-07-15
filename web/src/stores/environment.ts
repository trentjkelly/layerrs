import { readable } from "svelte/store"

// Comment / uncomment line based on dev & prod environments

// Production
export const urlBase = readable('https://layerrs.com')

// Development
// export const urlBase = readable('http://localhost:8080')
