import { logger } from "../lib/logger";
import { getUrlBase } from "../../stores/environment";

export async function loginServerRequest(email: string, password: string) : Promise<any> {
    if ((email !== '') && (password !== '')) {
        try {
            const res = await fetch(`${getUrlBase()}/api/authentication/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    email: email,
                    password: password
                })
            })

            const jsonData = await res.json()
            return jsonData
        } catch (error) {
            logger.error(`Failed to login: ${error}`)
            return null
        }
    } else {
        logger.error('Email or password is empty')
        return null
    }
}