// Sets the cookie for JWT & refresh token
export const POST = async ({ cookies, request }) => {
    try {
        const body = await request.json()

        const { jwtToken, refreshToken } = body
        
        cookies.set('refresh', refreshToken, {
            httpOnly: true, 
            secure: true, 
            path: "/", 
            sameSite: 'strict'
        })

        cookies.set('jwt', jwtToken, { 
            httpOnly: true, 
            secure: true, 
            path: "/", 
            sameSite: 'strict' 
        })

        return new Response('JWT Cookie set successfully', { status: 200 })

    } catch (error) {
        return new Response('Failed to process request', { status: 400 })
    }
}

// Deletes the cookies for JWT & refresh token
export const DELETE = async ({ cookies }) => {
    try {
        cookies.delete('jwt', { 
            httpOnly: true, 
            secure: true, 
            path: "/", 
            sameSite: 'strict' 
        })

        cookies.delete('refresh', { 
            httpOnly: true, 
            secure: true, 
            path: "/", 
            sameSite: 'strict' 
        })

        return new Response('JWT & Refresh Token cookie set successfully', { status: 200 })

    } catch (error) {
        return new Response('Failed to delete JWT & Refresh Token cookie', { status : 400 })
    }
}