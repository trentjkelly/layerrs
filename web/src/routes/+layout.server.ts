// Loads the jwt from cookies
export function load({ cookies }) {
    const newJWT = cookies.get('jwt')
    const newRefreshToken = cookies.get('refresh')

    return { newJWT, newRefreshToken }
}
