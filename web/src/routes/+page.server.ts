// Loads the jwt from cookies
export function load({ cookies }) {
    const refreshToken = cookies.get('refreshJWT')
    const token = cookies.get('jwt')
    return { token, refreshToken }
}
