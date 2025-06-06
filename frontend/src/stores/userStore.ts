import {makeAutoObservable} from "mobx";
import { BotConfig } from "@/types/Config";
import { IUser } from "@/types/user";


class userStore {
    IsAuth: boolean = false;
    User: IUser | undefined;
    Configs: BotConfig[] | undefined;
    isLoading: boolean = false;

    constructor () {
        makeAutoObservable(this);
    }
    SetIsLoading(value: boolean) {
        this.isLoading = value;
    }
    
    SetIsAuth(value: boolean) {
        this.IsAuth = value
    }
    SetUser(user: IUser) {
        this.User = user
    }
    SetConfigs(configs: BotConfig[]) {
        this.Configs = configs;
    }
    clear() {
        this.logout();
    }
    logout() {
        this.Configs = undefined;
        this.IsAuth = false;
        this.User = undefined;
        this.isLoading = false;
    }
}


export default userStore;