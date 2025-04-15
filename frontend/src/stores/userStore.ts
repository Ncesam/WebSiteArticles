import { IUser } from "@/types/user";


class userStore {
    IsAuth: boolean = false;
    User: IUser | undefined;

    async fetchUserData() {

    }

    
}


export default userStore;