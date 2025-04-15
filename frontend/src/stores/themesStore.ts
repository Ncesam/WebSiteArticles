import { ITheme, ThemeType } from "@/types/theme";
import { makeAutoObservable } from "mobx";



class themesStore {
    themes: ITheme[] | undefined;
    loading: boolean = false;
    error: string | null = null;

    constructor() {
        makeAutoObservable(this);
    }

    getFilteredThemes(filter: ThemeType) {
        return this?.themes?.filter((theme) => filter in theme.type)
    }

    async fetchThemes() {
        
    }


}

export default themesStore;