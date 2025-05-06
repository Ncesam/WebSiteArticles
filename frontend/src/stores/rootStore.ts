
import userStore from "./userStore";



class RootStore {
    userStore: userStore;
    constructor() {
        this.userStore = new userStore()
    }
}

export default RootStore;