import { StartConfig, FormBotConfig } from "@/types/Config"
import { $api } from "./api"


export class ConfigService {
    static async startConfig(config: StartConfig) {
        const {status, data} = await $api.post("/queue/", {Data: config.data, Prompt: config.prompt, UserId: 0, ConfigId: config.configId })
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