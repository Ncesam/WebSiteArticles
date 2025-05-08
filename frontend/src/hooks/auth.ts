import {useContext, useEffect, useState} from "react";
import {$api} from "@/http/api";
import {Context} from "@/index";
import { useStore } from "./store";


const useAuth = () => {
    useEffect(() => {
        const refreshToken = async () => {
            try {
                console.log("Refresh token...");
                await $api.post("/auth/refresh", {}, {withCredentials: true});
            } catch (error) {
                console.error("Ошибка обновления токена:", error);
                window.location.href = "/login";
            }
        };

        const interval = setInterval(() => {
            refreshToken();
        }, 20 * 1000 * 60);

        return () => clearInterval(interval);
    }, []);

};

export {useAuth};

const useAutoLogin = () => {
    const {userStore} = useStore();

    useEffect(() => {
        const checkAuth = async () => {
            try {
                userStore.SetIsLoading(true)
                const {status, data} = await $api.post("/auth/me", {withCredentials: true});

                if (status === 403) {
                    window.location.href = "/login";
                    return;
                }
                userStore.SetIsAuth(true);
                userStore.SetIsLoading(false);
                return;
            } catch (e) {
                console.log("error: ", e);
                userStore.SetIsLoading(false);
            }

        };

        checkAuth();
    }, []);

    return;
};

export {useAutoLogin};