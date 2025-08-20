import { jwt, isLoggedIn, refreshToken } from "../../stores/auth";
import { logger } from "./logger";

// Takes care of deleting cookies on the browser & in the stores
export async function handleBrowserLogout() : Promise<boolean> {
    try {
        const res  = await fetch('/cookies', { 
            method: 'DELETE'
        })
        if (!res.ok) {
            logger.error(`Failed to delete the JWT and refresh tokens: ${res.statusText}`);
            return false;
        }

        jwt.set('');
        refreshToken.set('');
        isLoggedIn.set(false);
        return true;

    } catch (error) {
        logger.error(`Failed to delete the JWT and refresh tokens: ${error}`);
        return false;
    }
}

// Takes care of setting cookies on the browser & in the stores
export async function handleBrowserLogin(newJwtToken: string, newRefreshToken: string) : Promise<boolean> {
    try {
        const res  = await fetch('/cookies', { 
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ 
                refreshToken: newRefreshToken,
                jwtToken: newJwtToken
            })
        })
        if (!res.ok) {
            logger.error(`Failed to set the JWT and refresh tokens: ${res.statusText}`);
            return false;
        }

        jwt.set(newJwtToken);
        refreshToken.set(newRefreshToken);
        isLoggedIn.set(true);
        return true;
    } catch (error) {
        logger.error(`Failed to set the JWT: ${error}`);
        return false;
    }
}