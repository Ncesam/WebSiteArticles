import {$api} from "@/http/api";


export class UserService {
    static async login(nickname: string, password: string) {
        const {status, data} = await $api.post("/api/auth/login", {
            nickname: nickname,
            password: password,
        })
        if (status === 200) {
            return data
        } else {
            return null
        }
    }
    static async register(email: string, nickname: string, password: string) {
        const {status, data} = await $api.post("/api/auth/register", {
            email: email,
            nickname: nickname,
            password: password,
        })
        if (status === 200) {
            return data
        } else {
            return null
        }
    }

    static async me() {
        const {status, data} = await $api.get("/api/auth/me")
        if (status === 200) {
            return data
        } else {
            return null
        }
    }
}