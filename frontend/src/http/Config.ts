import { StartConfig, FormBotConfig } from "@/types/Config"
import { $api } from "./api"


export class ConfigService {
    static async startConfig(config: StartConfig) {
        const formData = new FormData();
        formData.append("data", config.data); // файл
        formData.append("Prompt", config.prompt);
        formData.append("UserId", "0");
        formData.append("ConfigId", String(config.configId));
        const {status, data} = await $api.post("/queue/", formData, {
            headers: {
                "Content-Type": "miltipart/form-data"
            },
            maxBodyLength: Infinity,
            maxContentLength: Infinity,
        })
        if (status === 200) {
            return [true, data]
        } else {
            return [false, data]
        }
    }
    static async addConfig(config: FormBotConfig) {
        const {status, data} = await $api.post("/config/", {...config})
        if (status === 200) {
         return [true, data]
        } else {
         return [false, data]
        }
    }
    static async getConfigs() {
        const {status, data} = await $api.get("/config/");
        if (status === 200) {
            return [true, data]
        } else {
            return [false, data]
        }
    }
}