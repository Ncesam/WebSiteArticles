import {useEffect} from "react";
import { useStore } from "@/hooks/store";
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
                    userStore.SetIsLoading(false);
                    return;
                }
                userStore.SetIsAuth(true);
                userStore.SetIsLoading(false);
            } catch (e) {
                console.log("error: ", e);
                userStore.SetIsAuth(false);
                userStore.SetIsLoading(false);
            }
        };
        checkAuth();
    }, []);

    return;
};

export {useAutoLogin};