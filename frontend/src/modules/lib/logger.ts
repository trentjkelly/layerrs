import { getEnvironment } from "../../stores/environment"

// Logger class for logging messages if in development environment
class Logger {
    public constructor() {}

    public debug(message: string) {
        if (getEnvironment() === 'DEVELOPMENT') {
            console.log(message)
        }
    }

    public error(error: any) {
        if (getEnvironment() === 'DEVELOPMENT') {
            console.error(error)
        }
    }
}

// Singleton instance of the logger
const logger = new Logger()
export { logger }