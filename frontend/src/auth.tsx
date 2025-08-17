import axios from "axios";
import { useState, useEffect, useCallback, ReactNode, createContext, useContext } from "react";
import { useNavigate } from "react-router-dom";

type UserInfo = {
    name?: string,
    email?: string,
}

type AuthContextType = {
    loading?: boolean;
    userInfo?: UserInfo;
    logout?: () => void;
    loadUserInfo?: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({children}: {children: ReactNode}) {
    const [loading, setLoading] = useState<boolean>(false);
    const [userInfo, setUserInfo] = useState<UserInfo | undefined>(undefined);

    const loadUserInfo = useCallback(() => {
        setLoading(true);
        axios.get("/api/auth/me")
            .then(response => {
                setLoading(false);
                if (Object.keys(response.data).length !== 0) {
                    const userInfo = response.data as UserInfo;
                    setUserInfo(userInfo);
                } else {
                    setUserInfo(undefined);
                }
            })
            .catch((error) => {
                setLoading(false);
                setUserInfo(undefined);
            });
    }, []);

    useEffect(() => {
        loadUserInfo();
    }, []);

    const logout = useCallback(() => {
        setLoading(true);
        axios.post("api/auth/logout")
            .then(response => {
                setLoading(false);
                setUserInfo(undefined);
            })
            .catch((error) => {
                console.log(error);
                setLoading(false);
            })
    }, []);

    const value: AuthContextType = {
        loading, userInfo, logout, loadUserInfo
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const authContext = useContext(AuthContext);
    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return authContext;
}
