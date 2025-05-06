
import { BotConfig } from "@/types/Config";
import { IUser } from "@/types/user";


class userStore {
    IsAuth: boolean = false;
    User: IUser | undefined;
    Configs: BotConfig[] | undefined;
    
}


export default userStore;