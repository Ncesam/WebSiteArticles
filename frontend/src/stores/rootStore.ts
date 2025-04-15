import themesStore from "./themesStore";
import userStore from "./userStore";



class RootStore{
    userStore: userStore;
    themesStore: themesStore
    constructor() {
        this.userStore = new userStore(),
        this.themesStore = new themesStore()
    }
}

export default RootStore;