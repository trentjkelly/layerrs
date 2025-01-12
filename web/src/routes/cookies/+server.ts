// Sets the cookie for jwt
export const POST = async ({ cookies, request }) => {
    try {
        const body = await request.json()

        const { token } = body
        
        cookies.set('jwt', token, { 
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