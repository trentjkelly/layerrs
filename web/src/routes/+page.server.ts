// Loads the jwt from cookies
export function load({ cookies }) {
    const token = cookies.get('jwt')
    return { token }
}
