import { jwtDecode } from 'jwt-decode';

export const isUserAuthenticated = () => {
	const token = localStorage.getItem('jwt');
	if (!token) return false;

	try {
		const decoded = jwtDecode(token);
		return decoded.exp * 1000 > Date.now();
	} catch {
		return false;
	}
};

export function getMonthName(monthNum) {
	return new Intl.DateTimeFormat("en-US", { month: "long" }).format(new Date(2000, monthNum - 1));
}
