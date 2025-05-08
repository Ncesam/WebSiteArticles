import {$api} from "@/http/api";


export class UserService {
    static async login(nickname: string, password: string) {
        const {status, data} = await $api.post("/auth/login", {
            nickname: nickname,
            password: password,
        })
        if (status === 200) {
            return [true, data]
        } else {
            return [false, data]
        }
    }
    static async register(email: string, nickname: string, password: string) {
        const {status, data} = await $api.post("/auth/register", {
            email: email,
            nickname: nickname,
            password: password,
        })
        if (status === 200) {
            return [true, data]
        } else {
            return [false, data]
        }
    }

    static async me() {
        const {status, data} = await $api.get("/auth/me")
        if (status === 200) {
            return [true, data]
        } else {
            return [false, data]
        }
    }
    static async logout() {
        const {status, data} = await $api.post("/auth/logout");
        if (status === 200) {
            return [true, data]
        } else {
            return [false, data]
        }
    }
}