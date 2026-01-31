import { getUrlBase } from "../../stores/environment";

export async function loginServerRequest(email: string) : Promise<Response | null> {
        try {
            const res = await fetch(`${getUrlBase()}/api/authentication/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    email: email,
                })
            })

            return res
        } catch (error) {
            return null
        }
}