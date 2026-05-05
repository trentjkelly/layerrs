declare global {
	interface Window {
		Stripe: any;
	}
}

export async function loadStripe(publishableKey: string): Promise<void> {
	return new Promise((resolve) => {
		const script = document.createElement('script');
		script.src = 'https://js.stripe.com/v3/';
		script.onload = () => {
			(window as any).Stripe = (window as any).Stripe || {};
			(window as any).Stripe(publishableKey);
			resolve();
		};
		document.head.appendChild(script);
	});
}

export function getStripe(): any {
	return (window as any).Stripe;
}
