import {useContext, useEffect, useState} from "react";
import {$api} from "@/http/api";
import {Context} from "@/index";
import { useStore } from "./store";
import { UserService } from "@/http/User";


const useAuth = () => {
    useEffect(() => {
        const interval = setInterval(() => {
            UserService.refreshToken();
        }, 2 * 1000 * 60);

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
                await UserService.refreshToken()
                const [ok, data] = await UserService.me()

                if (!ok) {
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