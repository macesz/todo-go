import { createContext, useContext, useState } from "react";


const AuthContext = createContext();


const AuthProvider = ({ children }) => {

    const [user, setUser] = useState(() => {
        const token = localStorage.getItem("token");
        const userData = localStorage.getItem("userData");

        if (token && userData) {
            return { token, ...JSON.parse(userData) };
        }
        return null;
    });



    const login = (userData) => {

        localStorage.setItem("token", userData.token);

        const userDetails = {
            id: userData.id,
            name: userData.name,
            email: userData.email,
        }
        localStorage.setItem("user", JSON.stringify(userDetails));

        setUser(userData);
    };

    const logout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        setUser(null);
    };


    return (
        <AuthContext.Provider value={{ user, login, logout }}>
            {children}
        </AuthContext.Provider>
    )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => useContext(AuthContext);
export default AuthProvider;
