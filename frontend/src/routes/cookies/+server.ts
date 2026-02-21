// Sets the cookie for JWT & refresh token
export const POST = async ({ cookies, request }) => {
    try {
        const url = new URL(request.url)
        const body = await request.json()
        const { jwtToken, refreshToken } = body

        const host = url.hostname
        const isLocalhost = host === 'localhost' || host === '127.0.0.1'
        const secure = !isLocalhost

        cookies.set('refresh', refreshToken, {
            httpOnly: true,
            secure,
            path: '/',
            sameSite: 'lax'
        })

        cookies.set('jwt', jwtToken, {
            httpOnly: true,
            secure,
            path: '/',
            sameSite: 'lax'
        })

        return new Response('JWT & Refresh Token cookie set successfully', { status: 200 })
    } catch (error) {
        console.error('failed to set cookies:', error)
        return new Response('Failed to process request', { status: 400 })
    }
}

// Deletes the cookies for JWT & refresh token
export const DELETE = async ({ cookies, request }) => {
    try {
        const host = new URL(request.url).hostname
        const isLocalhost = host === 'localhost' || host === '127.0.0.1'
        const secure = !isLocalhost

        cookies.delete('jwt', { httpOnly: true, secure, path: '/', sameSite: 'lax' })
        cookies.delete('refresh', { httpOnly: true, secure, path: '/', sameSite: 'lax' })

        return new Response('JWT & Refresh Token cookie deleted successfully', { status: 200 })
    } catch (error) {
        return new Response('Failed to delete JWT & Refresh Token cookie', { status : 400 })
    }
}