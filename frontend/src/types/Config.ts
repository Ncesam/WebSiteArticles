
export interface BotConfig {
    name: string,
    id: number,
    userId: number,
    prompt: string,
    delay: number,
}

export interface FormBotConfig {
    name: string,
    userId: number,
    prompt: string,
    delay: number
}

export interface StartConfig {
    configId: number,
    userId: number,
    prompt: string,
    data: string,
}

