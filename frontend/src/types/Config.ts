
export interface BotConfig {
    name: string,
    id: number,
    userId: number,
    prompt: string,
    delay: number,
}

export interface FormBotConfig {
    name: string,
    prompt: string,
    delay: number
    email: string,
    password: string,
}

export interface StartConfig {
    configId: number,
    userId: number,
    prompt: string,
    data: string,
}

