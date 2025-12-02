import { createContext, useContext, useState } from "react";


const AuthContext = createContext();


const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null); // User state management



    const login = (token, userData) => {
        localStorage.setItem("token", token);
        setUser({ token, ...userData });
    };

    const logout = () => {
        localStorage.removeItem("token");
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
