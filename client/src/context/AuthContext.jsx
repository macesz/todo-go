import { createContext, useContext, useState } from "react";
import { isUserAuthenticated } from "../util/Util";


const AuthContext = createContext();


const AuthProvider = ({ children }) => {

    const [user, setUser] = useState(() => {
        if (isUserAuthenticated()) {
            const token = localStorage.getItem("token");
            const userData = localStorage.getItem("user");

            if (userData) {
                return { token, ...JSON.parse(userData) };
            }
        }

        localStorage.removeItem("token");
        localStorage.removeItem("user");
        return null;
    });



    const login = (userData) => {

        localStorage.setItem("token", userData.token);

        localStorage.setItem("user", JSON.stringify(userData.user));

        setUser({ token: userData.token, ...userData.user });
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
